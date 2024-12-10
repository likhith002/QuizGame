package service

import (
	"fmt"
	"log"
	"math/rand"
	"og_ed/internal/utility"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Player struct {
	Id            uuid.UUID       `json:"id"`
	Name          string          `json:"name"`
	Points        int             `json:"points,"`
	ProfilePic    string          `json:"profile"`
	isWordChoosen bool            `json:"-"`
	Connection    *websocket.Conn `json:"-"`
}

type GameState int

const (
	LobbyState GameState = iota
	PlayState
	UpdatePlayerState
	WaitState
	RevealState
	EndState
)

type Game struct {
	Id            uuid.UUID
	CurrentWord   string
	Code          string
	State         GameState
	Time          int
	Players       []*Player
	Coordinates   []CoordinatesPacket
	PlayerConn    map[string]uuid.UUID
	CurrentPlayer *Player
	GeneratedSets *map[string]struct{}
	levels        int32
	Host          *websocket.Conn
	netService    *NetService
	brodcastChan  chan struct{}
}

func generateCode() string {
	return strconv.Itoa(100000 + rand.Intn(900000))

}

func NewGame(host *websocket.Conn, netService *NetService) Game {

	return Game{
		Id:            uuid.New(),
		Code:          generateCode(),
		PlayerConn:    make(map[string]uuid.UUID),
		Players:       make([]*Player, 0),
		netService:    netService,
		Coordinates:   []CoordinatesPacket{},
		Time:          60,
		levels:        3,
		State:         LobbyState,
		Host:          host,
		GeneratedSets: &map[string]struct{}{},
	}
}

func (g *Game) BroadCastPacket(packet any, excludePlayers map[*websocket.Conn]struct{}) error {

	for _, player := range g.Players {
		if _, ok := excludePlayers[player.Connection]; !ok {

			err := g.netService.SendPacket(player.Connection, packet)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Game) loop() {

	for {

		select {

		case data := <-g.brodcastChan:
			{
				fmt.Println("DD", data)
			}

		}

	}

}

func (g *Game) startNewLevel() {

	//Update the current player

	g.CurrentPlayer = g.Players[0]

	time.Sleep(time.Second * 2)

	g.nextWord()

}

func (g *Game) startGame() error {

	if len(g.Players) > 0 {
		g.startNewLevel()
	} else {
		return fmt.Errorf("no players present in game")
	}
	return nil
}

func (g *Game) resetGame() {


	g.levels = 3

	for _, player := range g.Players {

		player.Points = 0
		player.isWordChoosen = false
	}

}

func (g *Game) resetLevel() {


	g.levels = g.levels-1

	for _, player := range g.Players {
		player.isWordChoosen = false
	}
	/////
	go cms

}

func (g *Game) getNextPlayer() *Player {

	for _, player := range g.Players {

		if !player.isWordChoosen {
			return player
		}

	}
	return nil
}

func (g *Game) nextWord() {

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	go g.Tick()

	g.Time = 60
	nextPlayer := g.getNextPlayer()

	nextPlayer.isWordChoosen = true

	g.CurrentPlayer = nextPlayer
	g.ChangeGameState(UpdatePlayerState,
		struct {
			Player Player `json:"player"`
		}{
			*g.CurrentPlayer,
		},
	)

	g.netService.SendPacket(
		nextPlayer.Connection,
		ChooseWord{
			Words: utility.GenerateUniqueRandomWords(g.GeneratedSets),
		},
	)

	g.BroadCastPacket(ChangeGameStatePacket{
		State: WaitState,
		Payload: struct {
			message string
		}{
			fmt.Sprintf("%s is choosing a word....", g.CurrentPlayer.Name),
		},
	}, nil)
}

func sortMapByValue(players []*Player) []*Player {

	// Step 2: Sort the slice by value (ascending order)
	sort.Slice(players, func(i, j int) bool {
		return players[i].Points < players[j].Points
	})

	// Step 3: Return the sorted slice
	return players
}

func (g *Game) processResult() {

	//sort the map by Points

	g.BroadCastPacket(LevelResult{
		Result: g.dereferencePlayers(sortMapByValue(g.Players)),
	}, nil)

}

func (g *Game) Start() {
	g.ChangeGameState(PlayState, struct{}{})

	if err := g.startGame(); err != nil {
		log.Fatal(err)
	}

	go g.loop()

	go func() {
		for {

			if g.levels == 0 {
				g.resetGame()

			}

			if g.Time == 0 {
				g.processResult()
				g.nextWord()
				g.Time = 60

			}

			time.Sleep(time.Second * 1)
		}
	}()
}

func (g *Game) Tick() {

	g.Time--

	g.BroadCastPacket(TickPacket{Tick: g.Time}, nil)

}

func (g *Game) ChangeGameState(state GameState, payload interface{}) {
	g.State = state
	g.BroadCastPacket(ChangeGameStatePacket{
		State:   state,
		Payload: payload,
	}, nil)

}

func (g *Game) dereferencePlayers(players []*Player) []Player {
	var playerValues []Player
	for _, p := range players {
		if p != nil {
			playerValues = append(playerValues, *p) // Dereference the pointer and append the value
		}
	}
	return playerValues

}

func (g *Game) OnPlayerAdd(name string, conn *websocket.Conn) {

	player := Player{
		Id:            uuid.New(),
		Points:        0,
		isWordChoosen: false,
		Name:          name,
		Connection:    conn,
	}
	g.Players = append(g.Players, &player)
	g.PlayerConn[conn.IP()] = player.Id

	g.netService.SendPacket(conn, ChangeGameStatePacket{
		State: g.State,
	})

	g.netService.SendPacket(conn, PlayerJoinPacket{
		Player:   player,
		GameCode: g.Code,
	})

	go func() {
		time.Sleep(time.Second * 2)
		g.netService.SendPacket(conn, GameSettings{
			Coordinates:   g.Coordinates,
			Players:       g.dereferencePlayers(g.Players),
			CurrentPlayer: g.CurrentPlayer,
		})
	}()

}
