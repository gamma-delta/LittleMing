package main

import (
	"fmt"
	"time"
	"strings"
)

var commandRouter = map[string]func(arg string, u *user, send func(string)){
	"help": cmdHelp,
	"addname": cmdAddName,
	"removename": cmdRemoveName,
	"listnames": cmdListNames,
	"start": cmdStart,
	"stop": cmdStop,
	"frequency": cmdFrequency,
	"whatfrequency": cmdWhatFrequency,
}

//yes yes I know no underscores in variable names
var _help = strings.ReplaceAll(`**Commands**
> - "~help": Get this list.
> - "~addname [name]": Add "name" to your list of chosen names.
> - "~removename [name]": Remove "name" from your list of chosen names.
> - "~listnames": Get a list of all your chosen names.
> - "~start": Start sending you messages.
> - "~stop": Stop sending you messages.
> - "~frequency [duration]": Set the average frequency of the message to "duration". Messages will be sent within a quarter of that time period (so if you set this to 2 hours, the time between messages will be between 1.5 hours and 2.5 hours.)
> - "~whatfrequency": Get the frequency of the message.`, "\"", "`") //please no annoying ```"`"+"`"` stuff

func cmdHelp(_ string, _ *user, send func(string)) {
	send(_help)
}

func cmdAddName(name string, u *user, send func(string)) {
	for _, alreadyName := range u.Names {
		if name == alreadyName {
			send(fmt.Sprintf("Hey, I've already been calling you %s, silly!", name))
			return
		}
	}
	u.Names = append(u.Names, name)
	send(fmt.Sprintf("Alright, I'll start calling you %s from now on!", name))
}

func cmdRemoveName(name string, u *user, send func(string)) {
	var idx int
	var alreadyName string
	for idx, alreadyName = range u.Names {
		if name == alreadyName {
			send("Alright, I'll stop calling you that from now on.")
			u.Names = append(u.Names[:idx], u.Names[idx+1:]...)
			return
		}
	}
	send("...but I'm not calling you that...")
}

func cmdListNames(_ string, u *user, send func(string)) {
	if len(u.Names) == 0 {
		send("You haven't told me what you want to be called yet!")
	} else {
		send("I'm going to call you one of these names:")
		for _, name := range u.Names {
			send("- " + name)
		}
	}
}

func cmdStart(_ string, u *user, send func(string)) {
	if len(u.Names) == 0 {
		send("You haven't told me what you want to be called yet!")
		return
	}
	if u.Enabled {
		//restart timer.
		u.SetupTimer()
		send("You already asked me to start sending you messages! So, I restarted the timer for you.")
	} else {
		//start timer.
		u.Enabled = true
		u.SetupTimer()
		send("OK! I'll start sending you messages!")
	}
}

func cmdStop(_ string, u *user, send func(string)) {
	if u.Enabled {
		//stop timer
		u.Enabled = false
		u.Timer.Stop()
		send("OK, I'll stop sending you messages.")
	} else {
		//wait it's already stopped ;(
		send("...but I wasn't sending you any messages...")
	}
}

func cmdFrequency(freqStr string, u *user, send func(string)) {
	freq, err := time.ParseDuration(freqStr)
	if err != nil {
		send("Hm, it looks like that wasn't a valid time... \n> Maybe reading `https://golang.org/pkg/time/#ParseDuration` could help you?")
		return
	}

	if freq.Minutes() < 1 {
		send("Aahhh... I can't send messages *that* fast! Please don't make it happen more than once a minute!")
		return
	}

	u.Frequency = freq
	if u.Enabled {
		//restart timer
		u.Timer.Stop()
		u.SetupTimer()
	}
	send("Frequency has been updated!")
}

func cmdWhatFrequency(_ string, u *user, send func(string)) {
	if u.Frequency.Hours() >= 1.5 {
		send(fmt.Sprintf("I'll send you a message about every %.1f hours.", u.Frequency.Hours()))
	} else {
		send(fmt.Sprintf("I'll send you a message about every %.1f minutes.", u.Frequency.Minutes()))
	}
}