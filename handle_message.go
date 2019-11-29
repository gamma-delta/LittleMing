package main

import (
	"strings"
	"github.com/bwmarrin/discordgo"
	"time"
)

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	//ignore my own messages, of course.
	if m.Author.ID == s.State.User.ID {
		return //noupe
	}

	//check to make sure this channel is a DM
	currentChannel, _ := s.Channel(m.ChannelID)
	if currentChannel.Type != discordgo.ChannelTypeDM {
		return
	}

	content := m.Content

	//make sure the message starts with ~
	if content[0] != '~' {
		return
	}
	//Nice, we can get started!

	id := m.Author.ID
	u, isKnownUser := ming.Users[id]
	if !isKnownUser {
		//better initialise them!
		s.ChannelMessageSend(m.ChannelID, "> Looks like you're a new user! Once you're done setting up, use `~start` to start the messages.")
		ming.Users[id] = &user{
			DMChannelID: m.ChannelID,
			Names: []string{},
			Frequency: 2 * time.Hour,
		}
		u = ming.Users[id] //gotta re-get the user
	}

	split := strings.SplitN(content[1:], " ", 2) //"addname Petra Kat" -> ["addname" "Petra Kat"]
	cmdStr := split[0]
	var args string
	if len(split) == 2 { args = split[1] }
	if cmd, ok := commandRouter[cmdStr]; ok { //yayayayaya we got a command~~~~
		out := strings.Builder{}
		send := func(output string) {
			out.WriteString("> ")
			out.WriteString(output)
			out.WriteRune('\n')
		}

		cmd(args, u, send)
		outStr := out.String()
		s.ChannelMessageSend(m.ChannelID, outStr[:len(outStr)-1]) //trim the last \n

	}
}