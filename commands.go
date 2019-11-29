package main

import (
	"fmt"
	"time"
)

var commandRouter = map[string]func(arg string, u *user, send func(string)){
	"addname": cmdAddName,
	"removename": cmdRemoveName,
	"listnames": cmdListNames,
	"start": cmdStart,
	"stop": cmdStop,
	"frequency": cmdFrequency,
	"whatfrequency": cmdWhatFrequency,
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