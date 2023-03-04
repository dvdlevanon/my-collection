package server

import (
	"my-collection/server/pkg/model"
	processor "my-collection/server/pkg/processor"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type push struct {
	server   *Server
	messages chan model.PushMessage
	upgrader websocket.Upgrader
}

func newPush(processor processor.Processor, server *Server) *push {
	result := &push{
		server:   server,
		messages: make(chan model.PushMessage),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
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

func (p *push) websocket(c *gin.Context) {
	ws, err := p.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("Unable to upgrade to websocket %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer ws.Close()

	for {
		select {
		case message := <-p.messages:
			err := ws.WriteJSON(message)
			if err != nil {
				logger.Warningf("Error writing to websocket %s", err)
				return
			}
		case <-time.After(60 * time.Second):
			err := ws.WriteJSON(model.PushMessage{MessageType: model.PUSH_PING})
			if err != nil {
				logger.Warningf("Error writing to websocket %s", err)
				return
			}
		}
	}
}
