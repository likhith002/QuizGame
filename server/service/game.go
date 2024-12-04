package service

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Player struct {
	Id         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Points     int             `json:"ponts"`
	ProfilePic string          `json:"profile"`
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
	Id            uuid.UUID
	CurrentWord   string
	Code          string
	State         GameState
	Time          int
	Players       []Player
	Coordinates   []CoordinatesPacket
	PlayerConn    map[string]uuid.UUID
	CurrentPlayer Player
	Rounds        int
	Host          *websocket.Conn
	netService    *NetService
}

func generateCode() string {
	return strconv.Itoa(100000 + rand.Intn(900000))

}

func NewGame(host *websocket.Conn) Game {

	return Game{
		Id:          uuid.New(),
		Code:        generateCode(),
		PlayerConn:  make(map[string]uuid.UUID),
		Players:     []Player{},
		Coordinates: []CoordinatesPacket{},
		Time:        60,
		State:       LobbyState,
		Host:        host,
	}
}

func (g *Game) BroadCastPacket(packet any, includeHost bool) error {

	for _, player := range g.Players {
		if player.Connection != g.Host {

			err := g.netService.SendPacket(player.Connection, packet)
			if err != nil {
				return err
			}
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

	// Select a player to be the

	go func() {
		for {

			if g.Time == 0 {
				break
			}
			g.Tick()

			time.Sleep(time.Second * 1)
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
		Points:     0,
		Name:       name,
		Connection: conn,
	}
	g.Players = append(g.Players, player)
	g.PlayerConn[conn.IP()] = player.Id

	g.netService.SendPacket(conn, ChangeGameState{
		State: g.State,
	})

	g.netService.SendPacket(conn, PlayerJoinPacket{
		Player:   player,
		GameCode: g.Code,
	})

	go func() {
		time.Sleep(time.Second * 2)
		g.netService.SendPacket(conn, GameSettings{
			Coordinates: g.Coordinates,
			Players:     g.Players,
		})
	}()
}
