package server

import (
	"my-collection/server/pkg/model"
	processor "my-collection/server/pkg/processor"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type push struct {
	server         *Server
	upgrader       websocket.Upgrader
	socketsChannel chan *websocket.Conn
	sockets        []*websocket.Conn
	messages       chan model.PushMessage
}

func newPush(processor processor.Processor, server *Server) *push {
	result := &push{
		server: server,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		sockets:        make([]*websocket.Conn, 0),
		socketsChannel: make(chan *websocket.Conn),
		messages:       make(chan model.PushMessage),
	}

	processor.SetProcessorNotifier(result)

	return result
}

func (p *push) OnFinishedTasksCleared() {
	p.pushQueueMetadata()
}

func (p *push) PauseToggled(paused bool) {
	p.pushQueueMetadata()
}

func (p *push) OnTaskAdded(task *model.Task) {
	p.pushQueueMetadata()
}

func (p *push) OnTaskComplete(task *model.Task) {
	p.pushQueueMetadata()
}

func (p *push) pushQueueMetadata() {
	queueMetadata, err := p.server.buildQueueMetadata()
	if err != nil {
		logger.Errorf("Unable to build queue metadata %s", err)
		return
	}

	p.messages <- model.PushMessage{MessageType: model.PUSH_QUEUE_METADATA, Payload: queueMetadata}
}

func (p *push) run() {
	for {
		select {
		case socket := <-p.socketsChannel:
			p.sockets = append(p.sockets, socket)
		case message := <-p.messages:
			p.writeToAll(&message)
		case <-time.After(30 * time.Second):
			p.writeToAll(&model.PushMessage{MessageType: model.PUSH_PING})
		}
	}
}

func (p *push) writeToAll(m *model.PushMessage) {
	for _, socket := range p.sockets {
		err := socket.WriteJSON(model.PushMessage{MessageType: model.PUSH_PING})
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "write: broken pipe" {
				p.removeAndClose(socket)
			} else {
				logger.Warningf("Error writing to websocket %s", err)
			}
		}
	}
}

func (p *push) removeAndClose(socket *websocket.Conn) {
	for i, cur := range p.sockets {
		if cur == socket {
			p.sockets = append(p.sockets[:i], p.sockets[i+1:]...)
			break
		}
	}

	socket.Close()
}

func (p *push) websocket(c *gin.Context) {
	ws, err := p.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("Unable to upgrade to websocket %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	p.socketsChannel <- ws
}
