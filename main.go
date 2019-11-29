package main

//https://discordapp.com/oauth2/authorize?client_id=%3c<ID>%3e&scope=bot&permissions=2048

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"io/ioutil"
	"encoding/json"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var ming *littleming
var saveFile = "data.json"
var messagesFile = "messages.txt"

func main() {
	log.Println("Started...")

	//Init ming
	ming = &littleming{
		Users: make(map[string]*user),
	}

	//Get saved users
	saveJSON, err := ioutil.ReadFile(saveFile)
	if err == nil {
		log.Println("Reading save file...")
		err = json.Unmarshal(saveJSON, &ming.Users)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("No save file found.")
	}

	//Get messages
	messagesBytes, err := ioutil.ReadFile(messagesFile)
	if err != nil {
		log.Fatal(err)
	}
	ming.Messages = strings.Split(string(messagesBytes), "\n")

	//Connect to discord
	ming.Discord, err = discordgo.New("Bot " + superSecretBotKey)
	if err != nil {
		log.Fatal(err)
	}

	//Add handlers...
	ming.Discord.AddHandler(handleMessage)

	//Open a websocket connection to discord
	err = ming.Discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	//get old users back up and running
	for _, user := range ming.Users {
		if user.Enabled {
			user.SetupTimer()
		}
	}

	//all good!
	log.Printf("Bot is now running @ %s#%s!\n", ming.Discord.State.User.Username, ming.Discord.State.User.Discriminator)

	//setup kill signals config
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-signals //wait for a signal
	
	log.Println("Stopping... the bot now will save user config and gracefully exit.")
	jsonData, _ := json.Marshal(ming.Users)
	ioutil.WriteFile(saveFile, jsonData, 0644)
	log.Println("Saved users.")
	ming.Discord.Close()
}