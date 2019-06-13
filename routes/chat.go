package routes

import (
    "hypermedlab/backend-myblog/services/chat"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Chat struct {
	engine *gin.Engine
}

type Message struct {
	Type    string
	Content string
	To      string
	From    string
}

func NewChatRouter(engine *gin.Engine) *Chat {
	return &Chat{
		engine: engine,
	}
}

func (c *Chat) Init() {
	c.engine.GET("/chat", c.chatHandler)
}

func (c *Chat) chatHandler(ctx *gin.Context) {
	conn, err := websocket.Upgrade(ctx.Writer, ctx.Request, ctx.Writer.Header(), 200, 200)
	if err != nil {
		ctx.String(400, err.Error())
	}
	defer conn.Close()

	for {
		message := new(Message)
		err := conn.ReadJSON(message)
		if err != nil {
			conn.WriteJSON(message)
		}

		switch message.Type {
		case "close":
			return
		case "message":
			
		case "join":
		case "exit":
		case: "invite":

		default:
		}
	}
}
