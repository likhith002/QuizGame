package service

import (
	"math/rand"
	"og_ed/entity"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Player struct {
	Id         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Connection *websocket.Conn `json:"-"`
}

type GameState int

const (
	LobbyState  GameState = 0
	PlayState   GameState = 1
	RevealState GameState = 2
	EndState    GameState = 3
)

type Game struct {
	Id              uuid.UUID
	Quiz            entity.Quiz
	CurrentQuestion int
	Code            string
	State           GameState
	Time            int
	Players         []Player
	Host            *websocket.Conn
	netService      *NetService
}

func generateCode() string {
	return strconv.Itoa(100000 + rand.Intn(900000))

}

func NewGame(quiz entity.Quiz, host *websocket.Conn) Game {

	return Game{
		Id:      uuid.New(),
		Quiz:    quiz,
		Code:    generateCode(),
		Players: []Player{},
		Time:    60,
		State:   LobbyState,
		Host:    host,
	}
}

func (g *Game) BroadCastPacket(packet any, includeHost bool) error {

	for _, player := range g.Players {
		err := g.netService.SendPacket(player.Connection, packet)
		if err != nil {
			return err
		}
	}
	if includeHost {
		err := g.netService.SendPacket(g.Host, packet)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) Start() {
	g.ChangeGameState(PlayState)
	go func() {
		for {
			g.Tick()
			time.Sleep(time.Second * 2)
		}

	}()
}

func (g *Game) Tick() {

	g.Time--
	g.netService.SendPacket(g.Host, TickPacket{Tick: g.Time})

}

func (g *Game) ChangeGameState(state GameState) {
	g.State = state
	g.BroadCastPacket(ChangeGameState{
		State: state,
	}, true)

}

func (g *Game) OnPlayerAdd(name string, conn *websocket.Conn) {

	player := Player{
		Id:         uuid.New(),
		Name:       name,
		Connection: conn,
	}
	g.Players = append(g.Players, player)

	g.netService.SendPacket(conn, ChangeGameState{
		State: g.State,
	})

	g.netService.SendPacket(g.Host, PlayerJoinPacket{
		Player: player,
	})
}
