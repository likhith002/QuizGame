package service

import (
	"fmt"
	"log"
	"math/rand"
	"og_ed/internal/logger"
	"og_ed/internal/utility"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	UpdateLevelState
	RevealState
	EndState
)
const LEVEL_TIME int = 10
const TOTAL_LEVELS int = 1

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
	levels        int
	Host          *websocket.Conn
	netService    *NetService
	logger        *logrus.Logger
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
		Time:          LEVEL_TIME,
		levels:        TOTAL_LEVELS,
		State:         LobbyState,
		Host:          host,
		GeneratedSets: &map[string]struct{}{},
		logger:        logger.GetLogger(),
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

// func (g *Game) loop() {

// 	for {

// 		select {

// 		case data := <-g.brodcastChan:
// 			{
// 				fmt.Println("DD", data)
// 			}

// 		}

// 	}

// }

func (g *Game) startGame() error {

	if len(g.Players) > 0 {
		g.nextWord()
	} else {
		return fmt.Errorf("no players present in game")
	}
	return nil
}

func (g *Game) resetGame() {

	for _, player := range g.Players {

		player.Points = 0
		player.isWordChoosen = false
	}
	g.GeneratedSets = nil

	g.levels = TOTAL_LEVELS
	time.Sleep(time.Second * 3)
	g.Time = LEVEL_TIME

}

func (g *Game) resetLevel() int {

	g.levels = g.levels - 1

	g.freezeLevel()

	if g.levels == 0 {
		return 0
	}

	for _, player := range g.Players {
		player.isWordChoosen = false
	}

	g.BroadCastPacket(ChangeGameStatePacket{
		State: UpdateLevelState,
		Payload: struct {
			Level int `json:"level"`
		}{
			Level: TOTAL_LEVELS - g.levels + 1,
		},
	}, nil)

	return 1

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

	g.Time = LEVEL_TIME
	g.Coordinates = []CoordinatesPacket{}
	nextPlayer := g.getNextPlayer()

	if nextPlayer == nil {

		fmt.Println("Calling reset level")
		status := g.resetLevel()

		if status == 0 {
			g.resetGame()
			return
		} else {

			nextPlayer = g.getNextPlayer()

		}

	}

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
	g.logger.Info("Choosing a word")
	g.BroadCastPacket(ChangeGameStatePacket{
		State: WaitState,
		Payload: struct {
			message string
		}{
			fmt.Sprintf("%s is choosing a word....", g.CurrentPlayer.Name),
		},
	}, map[*websocket.Conn]struct{}{
		nextPlayer.Connection: {},
	})
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
	var levelType string = "intermediate"
	if g.levels == 0 {
		levelType = "final"
	}

	g.BroadCastPacket(LevelResult{
		Result: g.dereferencePlayers(sortMapByValue(g.Players)),
		Type:   levelType,
	}, nil)

}

func (g *Game) freezeLevel() {
	g.processResult()

}

func (g *Game) Start() {
	// g.ChangeGameState(PlayState, struct{}{})

	if err := g.startGame(); err != nil {
		log.Fatal(err)
	}

	// go g.loop()

	go func() {
		for {

			if g.Time == 0 {

				g.logger.Info("Freezing the level")
				g.nextWord()

			}

			time.Sleep(time.Second * 5)
		}
	}()
}

func (g *Game) Tick() {

	for {

		if g.Time == 0 {
			break
		}
		g.Time--
		time.Sleep(time.Second)
		g.BroadCastPacket(TickPacket{Tick: g.Time}, nil)
	}

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
		time.Sleep(time.Second * 3)

		g.BroadCastPacket(GameSettings{
			Coordinates:   g.Coordinates,
			Players:       g.dereferencePlayers(g.Players),
			CurrentPlayer: g.CurrentPlayer,
		}, nil)

	}()

}
