package main

import (
	"log"

	"github.com/lmatte7/gomesh/github.com/meshtastic/gomeshproto"
)

// nodes!

type Node struct {
	Num       uint32
	UserId    string
	LongName  string
	ShortName string
	SNR       float32
	Hops      uint32
}

// get a node's long name if it exists
func getNodeName(nodeNum uint32) string {
	for _, response := range nodes {
		if n, ok := response.GetPayloadVariant().(*gomeshproto.FromRadio_NodeInfo); ok {
			// make sure the user struct exists if we're trying to return a name
			if n.NodeInfo.User != nil && nodeNum == n.NodeInfo.Num {
				return n.NodeInfo.User.LongName
			}
		}
	}
	return "unknown"
}

// refresh our global list of nodes if we get a good pull
func updateNodes() {
	newNodes, err := radio.GetRadioInfo()
	if err != nil {
		log.Println("Couldn't get nodes", err)
	} else {
		nodes = newNodes
	}
}
