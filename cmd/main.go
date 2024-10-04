package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AsaHero/just-ask-bot/internal/app"
	"github.com/AsaHero/just-ask-bot/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// config
	config := config.New()

	// app
	app := app.NewApp(config)

	// run application
	go func() {
		log.Println("Listen: ", "address", config.Server.Host+config.Server.Port)
		if err := app.Run(); err != nil {
			log.Fatalln("app run", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// app stops
	log.Println("api esoterica bot stops")
	app.Stop()
}
