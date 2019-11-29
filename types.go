package main

import (
	"time"
	"math/rand"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

//the server itself. each application should only have one of these running.
type littleming struct {
	Users   map[string]*user //map discord ID to user
	Discord *discordgo.Session
	Messages []string
}

//one person using the bot
type user struct {
	DMChannelID string `json:"channelid"`
	Enabled   bool          `json:"enabled"`
	Names     []string      `json:"names"`
	Frequency time.Duration `json:"frequency"`
	Timer *time.Timer `json:"-"` //this counts down the time till the next message.
}

//Use this to start a timer.
func (u *user) SetupTimer() {
	// fmt.Printf("DEBUG: Timer setup from %#v\n", u)
	u.Enabled = true
	freq := int64(u.Frequency)
	duration := time.Duration(freq + rand.Int63n(freq / 2) - freq / 4)
	u.Timer = time.AfterFunc(duration, func() {
		ming.Discord.ChannelMessageSend(u.DMChannelID, fmt.Sprintf(
			ming.Messages[rand.Intn(len(ming.Messages))], 
			u.Names[rand.Intn(len(u.Names))]))
		u.SetupTimer() //call itself again.
	})
}