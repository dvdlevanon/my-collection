package push

import (
	"context"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server"
	"my-collection/server/pkg/utils"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var logger = utils.MustGetLogger("push")

type PushHandler interface {
	server.Handler
	Run(ctx context.Context) error
	Push(m model.PushMessage)
}

func NewPush() PushHandler {
	return &push{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		sockets:        make([]*websocket.Conn, 0),
		socketsChannel: make(chan *websocket.Conn),
		messages:       make(chan model.PushMessage),
	}
}

type push struct {
	upgrader       websocket.Upgrader
	socketsChannel chan *websocket.Conn
	sockets        []*websocket.Conn
	messages       chan model.PushMessage
}

func (p *push) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ws", p.websocket)
}

func (p *push) Push(m model.PushMessage) {
	p.messages <- m
}

func (p *push) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Infof(ctx, "Push service shutting down...")
			// Close all websocket connections
			for _, socket := range p.sockets {
				socket.Close()
			}
			return nil
		case socket := <-p.socketsChannel:
			p.sockets = append(p.sockets, socket)
		case message := <-p.messages:
			p.writeToAll(ctx, &message)
		case <-ticker.C:
			p.writeToAll(ctx, &model.PushMessage{MessageType: model.PUSH_PING})
		}
	}
}

func (p *push) writeToAll(ctx context.Context, m *model.PushMessage) {
	for _, socket := range p.sockets {
		err := socket.WriteJSON(m)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "write: broken pipe" {
				p.removeAndClose(socket)
			} else {
				logger.Warningf(ctx, "Error writing to websocket %s", err)
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
	ctx := server.ContextWithSubject(c)
	ws, err := p.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf(ctx, "Unable to upgrade to websocket %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	p.socketsChannel <- ws
}
