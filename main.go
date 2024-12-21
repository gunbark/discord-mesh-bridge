package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load our Discord auth token
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	usbPort := os.Getenv("USB_PORT")
	if len(usbPort) == 0 {
		log.Fatal("No USB port defined in .env")
	}
	discordAuthToken := os.Getenv("DISCORD_AUTH_TOKEN")
	if len(discordAuthToken) == 0 {
		log.Fatal("No Discord auth token defined in .env")
	}
	discordPrivateChannel := os.Getenv("DISCORD_PRIVATE_CHANNEL")
	discordPublicChannel := os.Getenv("DISCORD_PUBLIC_CHANNEL")
	if len(discordPrivateChannel) == 0 && len(discordPublicChannel) == 0 {
		log.Fatal("No Discord channels defined to sync to in .env")
	}

	// Start a new session
	s, err := discordgo.New("Bot " + discordAuthToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	// Open the session
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// begin our message processor
	go discordQueue(s, discordPrivateChannel, discordPublicChannel)

	// Mesh start
	meshStart(usbPort)

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord and Mesh sessions
	s.Close()
	meshStop()
}
