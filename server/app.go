package internal

import (
	"context"
	"log"
	"og_ed/internal/collection"
	"og_ed/internal/controller"
	"og_ed/service"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var quizCollection *mongo.Collection

type App struct {
	server      *fiber.App
	db          *mongo.Database
	quizService *service.QuizService
	netService  *service.NetService
}

func (a *App) Init() {

	a.setUpDb()
	a.setUpServices()
	a.setUpServer()

	log.Fatal(a.server.Listen(":5001"))

}
func (a *App) setUpServer() {
	app := fiber.New()

	app.Use(cors.New())

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })

	// app.Get("/api/quizzes", getQuizes)

	// app.Get("/ws", websocket.New(func(c *websocket.Conn) {
	// 	// c.Locals is added to the *websocket.Conn
	// 	log.Println(c.Locals("allowed"))  // true
	// 	log.Println(c.Params("id"))       // 123
	// 	log.Println(c.Query("v"))         // 1.0
	// 	log.Println(c.Cookies("session")) // ""

	// 	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	// 	var (
	// 		mt  int
	// 		msg []byte
	// 		err error
	// 	)
	// 	for {
	// 		if mt, msg, err = c.ReadMessage(); err != nil {
	// 			log.Println("read:", err)
	// 			break
	// 		}
	// 		log.Printf("recv: %s", msg)

	// 		if err = c.WriteMessage(mt, msg); err != nil {
	// 			log.Println("write:", err)
	// 			break
	// 		}
	// 	}

	// }))

	quizController := controller.Quiz(a.quizService)
	app.Get("/api/quizzes", quizController.GetQuizzes)

	wsController := controller.Ws(a.netService)
	app.Get("/ws", websocket.New(wsController.Ws))

	a.server = app
}

func (a *App) setUpDb() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin%40123@localhost:27017"))

	if err != nil {
		panic(err)
	}

	a.db = client.Database("quiz")

}

func (a *App) setUpServices() {

	a.quizService = service.Quiz(collection.Quiz(*a.db.Collection("quizzes")))
	a.netService = service.Net(a.quizService)

}
