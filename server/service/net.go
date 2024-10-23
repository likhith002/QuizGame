package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"og_ed/entity"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NetService struct {
	quizService *QuizService
	games       []*Game
	host        *websocket.Conn
	tick        int
}

func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
		games:       []*Game{},
	}
}

type ConnectionPacket struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type HostGamePacket struct {
	QuizId string `json:"quiz_id"`
}

type ShowQuestionPacket struct {
	Quetion entity.QuizQuestion `json:"question"`
}

type ChangeGameState struct {
	State GameState `json:"state"`
}

type PlayerJoinPacket struct {
	Player Player `json:"player"`
}
type StartGamePacket struct {
}
type TickPacket struct {
	Tick int `json:"tick"`
}

func (c *NetService) getGameByCode(code string) *Game {

	for _, game := range c.games {

		if game.Code == code {
			return game
		}

	}
	return nil
}

func (c *NetService) getGameByHost(host *websocket.Conn) *Game {

	for _, game := range c.games {

		if game.Host == host {
			return game
		}

	}
	return nil
}

func (c *NetService) packetIdToPacket(pid uint8) interface{} {

	switch pid {
	case 0:
		return &ConnectionPacket{}

	case 1:
		return &HostGamePacket{}

	case 5:
		return &StartGamePacket{}
	}

	return nil
}

func (c *NetService) packetToPacketId(packet interface{}) (uint8, error) {

	switch packet.(type) {

	case ShowQuestionPacket:
		return 2, nil

	case ChangeGameState:
		return 3, nil

	case PlayerJoinPacket:
		return 4, nil
	}

	return 0, errors.New("invalid type provided")

}

func (c *NetService) OnIncomingMessage(conn *websocket.Conn, mt int, msg []byte) {

	if len(msg) < 2 {
		return
	}

	pId := msg[0]
	data := msg[1:]

	packet := c.packetIdToPacket(pId)
	// fmt.Println(pId, string(data))
	err := json.Unmarshal(data, packet)

	if err != nil {

		log.Fatal(err)
	}

	switch data := packet.(type) {

	case *ConnectionPacket:
		{

			game := c.getGameByCode(data.Code)

			if game == nil {
				return
			}

			game.OnPlayerAdd(data.Name, conn)

			// fmt.Println(data.Name, "wants to join game", data.Code)
			break
		}

	case *HostGamePacket:
		{

			quizId, err := primitive.ObjectIDFromHex(data.QuizId)

			if err != nil {
				fmt.Println(err)
				return
			}

			quiz, err := c.quizService.quizCollection.GetById(quizId)

			if err != nil {
				fmt.Println(err)
				return
			}

			if quiz == nil {

				return
			}
			newGame := NewGame(*quiz, conn)
			fmt.Println("Code for New Game", newGame.Code)
			c.games = append(c.games, &newGame)

			c.SendPacket(conn, ChangeGameState{
				State: LobbyState,
			})

		}
	case *StartGamePacket:
		{

			game := c.getGameByHost(conn)

			if game == nil {
				return
			}
			game.Start()
			return
		}
	}

	// str := string(msg)
	// splittedString := strings.Split(str, ":")

	// cmd := splittedString[0]
	// arg := splittedString[1]

	// switch cmd {

	// case "host":
	// 	{
	// 		fmt.Println("host quiz", arg)
	// 		c.host = conn
	// 		c.tick = 100

	// 		go func() {

	// 			for {
	// 				c.tick--
	// 				c.host.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(c.tick)))
	// 				time.Sleep(time.Second)

	// 			}

	// 		}()
	// 		break
	// 	}

	// case "join":
	// 	{
	// 		fmt.Println("join with code", arg)
	// 		if c.host != nil {
	// 			c.host.WriteMessage(websocket.TextMessage, []byte("A new player joined the game"))
	// 		}

	// 	}

	// }

}

func (c *NetService) PacketToBytes(packet interface{}) ([]byte, error) {

	pId, err := c.packetToPacketId(packet)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}
	final := append([]byte{pId}, bytes...)
	return final, nil

}

func (c *NetService) SendPacket(conn *websocket.Conn, packet interface{}) error {

	bytes, err := c.PacketToBytes(packet)
	if err != nil {

		return err
	}

	return conn.WriteMessage(websocket.BinaryMessage, bytes)
}
