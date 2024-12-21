package main

import (
	"log"

	"github.com/lmatte7/gomesh/github.com/meshtastic/gomeshproto"
)

// channels!

type Channel struct {
	Index   uint32
	Name    string
	Private bool
	Role    string
}

// is a given channel private
func getChannelPrivate(channelIndex uint32, to uint32) bool {
	// return private if it was a DM (not "to" all FFFF..)
	if to != 4294967295 {
		return true
	}
	for _, channel := range channels {
		if channel.Index == int32(channelIndex) {
			// treats primary as public, otherwise checks for PSK
			// if it has a private key, it ain't public
			if channel.GetRole() == gomeshproto.Channel_PRIMARY {
				return false
			} else {
				return len(channel.Settings.Psk) > 0
			}
		}
	}
	return true
}

// get a channel's name if it exists
func getChannelName(channelIndex uint32, to uint32) string {
	// return channel name as "DM" if it was a DM (not "to" all FFFF..)
	if to != 4294967295 {
		return "DM"
	}
	for _, channel := range channels {
		if channel.Index == int32(channelIndex) {
			// Labels primary as "Public", otherwise pulls a name
			if channel.GetRole() == gomeshproto.Channel_PRIMARY {
				return "Public"
			} else {
				return channel.Settings.Name
			}
		}
	}
	return ""
}

// refresh our global list of channels if we get a good pull
func updateChannels() {
	newChannels, err := radio.GetChannels()
	if err != nil {
		log.Println("Couldn't get channels", err)
	} else {
		channels = newChannels
	}
}
