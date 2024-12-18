package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"og_ed/entity"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MAX_PLAYERS = 6

type NetService struct {
	quizService *QuizService
	games       []*Game
	mu          sync.Mutex
}

func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
		games:       []*Game{},
	}
}

type ConnectionPacket struct {
	Name string `json:"name"`
}

type HostGamePacket struct {
	QuizId string `json:"quizId"`
}

type ShowQuestionPacket struct {
	Question entity.QuizQuestion `json:"question"`
}

type ChangeGameStatePacket struct {
	State   GameState   `json:"state"`
	Payload interface{} `json:"payload"`
}

type PlayerJoinPacket struct {
	Player   Player `json:"player"`
	GameCode string `json:"gameCode"`
}
type StartGamePacket struct {
}

type LevelResult struct {
	Result []Player `json:"result"`
	Type   string   `json:"type"`
}

type TickPacket struct {
	Tick int `json:"tick"`
}
type GameSettings struct {
	Coordinates   []CoordinatesPacket `json:"coordinates"`
	Players       []Player            `json:"players"`
	CurrentPlayer *Player             `json:"currentPlayer,omitempty"`
}

type CoordinatesPacket struct {
	X1        float32 `json:"x1"`
	Y1        float32 `json:"y1"`
	X2        float32 `json:"x2"`
	Y2        float32 `json:"y2"`
	Color     string  `json:"color"`
	LineWidth string  `json:"lineWidth"`
}

type ChooseWord struct {
	Words []string `json:"words"`
}

type SelectedWordPacket struct {
	Woerd string `json:"word"`
}

func getRandomIndex(arr []*Game) (int, error) {
	// Check if the array is empty
	if len(arr) == 0 {
		return -1, fmt.Errorf("no happening games")
	}

	var randomIndex int
	for {
		rand.Seed(time.Now().UnixNano())

		randomIndex = rand.Intn(len(arr))

		if len(arr[randomIndex].Players) != MAX_PLAYERS {
			return randomIndex, nil

		}

	}

}

func (c *NetService) addToGame(conn *websocket.Conn, name string) *Game {

	index, _ := getRandomIndex(c.games)

	if index == -1 {
		newGame := NewGame(conn, c)
		c.games = append(c.games, &newGame)
		return &newGame

	} else {
		return c.games[index]
	}

}

func (c *NetService) getGameByHost(host *websocket.Conn) *Game {

	for _, game := range c.games {

		if game.Host == host {
			return game
		}

	}
	return nil
}

func (c *NetService) getGameByConn(conn *websocket.Conn) *Game {

	for _, game := range c.games {

		if _, flag := game.PlayerConn[conn.IP()]; flag {
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

	case 7:
		return &CoordinatesPacket{}
	case 10:
		return &SelectedWordPacket{}
	}

	return nil
}

func (c *NetService) packetToPacketId(packet interface{}) (uint8, error) {

	switch packet.(type) {

	case ShowQuestionPacket:
		return 2, nil

	case ChangeGameStatePacket:
		return 3, nil

	case PlayerJoinPacket:
		return 4, nil
	case TickPacket:
		return 6, nil

	case CoordinatesPacket:
		return 7, nil
	case GameSettings:
		return 8, nil
	case ChooseWord:
		return 9, nil
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

	err := json.Unmarshal(data, packet)

	if err != nil {

		log.Fatal(err)

	}

	switch data := packet.(type) {
	case *ConnectionPacket:
		{

			game := c.addToGame(conn, data.Name)

			if game == nil {
				return
			}

			currentPlayersCount := len(game.Players)

			game.OnPlayerAdd(data.Name, conn)

			if currentPlayersCount == 1 && game.State == LobbyState {

				game.Start()
			}

			break
		}

	case *HostGamePacket:
		{

			_, err := primitive.ObjectIDFromHex(data.QuizId)

			if err != nil {
				fmt.Println(err)

				return
			}

			// quiz, err := c.quizService.quizCollection.GetById(quizId)

			if err != nil {
				fmt.Println(err)
				return
			}

			newGame := NewGame(conn, c)
			c.games = append(c.games, &newGame)

			c.SendPacket(conn, ChangeGameStatePacket{
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
	case *CoordinatesPacket:
		{

			game := c.getGameByConn(conn)

			if game == nil {
				return
			}

			fmt.Println("Got coordinates packet")

			game.Coordinates = append(game.Coordinates, *data)

			err = game.BroadCastPacket(*data, nil)

			if err != nil {
				log.Fatal("Failed in broadcasting...", err)
			}
		}
	case *SelectedWordPacket:
		{

			game := c.getGameByConn(conn)

			if game == nil {
				return
			}
			game.BroadCastPacket(ChangeGameStatePacket{
				State: PlayState,
			}, nil)
			time.Sleep(time.Second * 2)
			game.BroadCastPacket(GameSettings{
				Coordinates:   game.Coordinates,
				Players:       game.dereferencePlayers(game.Players),
				CurrentPlayer: game.CurrentPlayer,
			}, nil)

			go game.Tick()

		}
	default:
		{

			fmt.Println(pId, data)
		}
	}

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

	c.mu.Lock()         // Lock before writing
	defer c.mu.Unlock() // Unlock

	bytes, err := c.PacketToBytes(packet)

	if err != nil {

		return err
	}

	return conn.WriteMessage(websocket.BinaryMessage, bytes)
}
