package controller

import (
	"log"
	"og_ed/service"

	"github.com/gofiber/contrib/websocket"
)

type WebsocketController struct {
	netService *service.NetService
}

func Ws(netService *service.NetService) WebsocketController {
	return WebsocketController{
		netService: netService,
	}
}

func (c *WebsocketController) Ws(conn *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)

	for {
		if mt, msg, err = conn.ReadMessage(); err != nil {
			log.Panicln("read", err)
			// fmt.Print(err)
			break
		}

		c.netService.OnIncomingMessage(conn, mt, msg)

	}
}
