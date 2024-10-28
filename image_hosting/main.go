package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"treehole_image_proxy/api"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler:          common.ErrorHandler,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		BodyLimit:             128 * 1024 * 1024,
	})
	// app.Use(recover.New(recover.Config{
	// 	EnableStackTrace:  true,
	// 	StackTraceHandler: common.StackTraceHandler,
	// }))
	// app.Use(common.MiddlewareGetUserID)
	//app.Use(common.MiddlewareCustomLogger)
	router := app.Group("/api")
	router.Post("/json", api.UploadImage)
	router.Get("/i/:year/:month/:day/:identifier", api.GetImage)

	go func() {
		err := app.Listen(":8000")
		if err != nil {
			log.Println(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		log.Println(err)
	}
}
