package main

import (
	"log"
	"time"

	"github.com/lmatte7/gomesh"
	"github.com/lmatte7/gomesh/github.com/meshtastic/gomeshproto"
)

var radio gomesh.Radio
var channels []*gomeshproto.Channel
var nodes []*gomeshproto.FromRadio

// get a radio! may not immediately respond as radio can take a minute to
// reset after last access, e.g. if this is a service give it a restart delay
func startRadio(port string) error {
	err := radio.Init(port)
	if err != nil {
		log.Println("Error setting radio port: %v", err)
		return err
	}
	return nil
}

// entry point from main()
func meshStart(port string) {
	// start a loop that sleeps between start attempts
	for true {
		listen(port)
		time.Sleep(5 * time.Second)
	}
}

// Shut down the radio on a ctrl-c from main
func meshStop() {
	radio.Close()
}

// listening loop, return meshStart() if an error pops
func listen(port string) {
	err := startRadio(port)
	if err != nil {
		log.Println("Couldn't get radio")
		return
	}
	defer radio.Close()

	log.Println(".. listening on radio", port)
	count := 0
	for {
		// make sure we update our list of channels and nodes every once in a while
		if count == 0 {
			time.Sleep(100 * time.Millisecond)
			updateChannels()
			time.Sleep(100 * time.Millisecond)
			updateNodes()
		}
		count++
		if count == 5000 {
			count = 0
		}

		responses, _ := radio.ReadResponse(false)
		for _, response := range responses {
			if packet, ok := response.GetPayloadVariant().(*gomeshproto.FromRadio_Packet); ok {
				if packet.Packet.GetDecoded().GetPortnum() == gomeshproto.PortNum_TEXT_MESSAGE_APP {
					receiveBuffer.Push(*packet)
				}
			}
		}
	}
}
