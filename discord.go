package main

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// discord! honestly not much happens here, just a webhook
// async processing of messages, mostly message formatting
func discordQueue(s *discordgo.Session, privateChannel string, publicChannel string) {
	// check every second for messages to process, as it calls hardware and
	// hardware can't process as fast as code runs
	for true {
		if receiveBuffer.Size() > 0 {
			// baseline info
			p := receiveBuffer.Pop()
			name := getNodeName(p.Packet.From)
			channel := getChannelName(p.Packet.Channel, p.Packet.To)
			private := getChannelPrivate(p.Packet.Channel, p.Packet.To)
			sourceNode := strconv.FormatUint(uint64(p.Packet.From), 10)
			re := regexp.MustCompile(`\r?\n`)
			message := re.ReplaceAllString(string(p.Packet.GetDecoded().Payload), "")
			snr := p.Packet.RxSnr

			// if we're logging, save out
			if logging {
				info := "{from:" + strconv.FormatFloat(float64(p.Packet.From), 'f', 0, 32) + ","
				info += "name:" + name + ","
				info += "channel:" + channel + ","
				info += "private:" + strconv.FormatBool(private) + ","
				info += "snr:" + strconv.FormatFloat(float64(snr), 'f', 2, 32) + ","
				info += "msg:\"" + message + "\"}"
				saveLog("mesh", channel, info)
			}

			// discord send
			discordMsg := "###\n"
			discordMsg += "\nmesh from: **" + name + "**\n"
			discordMsg += "```" + message + "```"
			discordMsg += "_channel: " + channel + ", "
			if private {
				discordMsg += "private, "
			} else {
				discordMsg += "public, "
			}
			discordMsg += "SNR: " + strconv.FormatFloat(float64(snr), 'f', 2, 32) + "_\n"
			discordMsg += "_source node: " + sourceNode + "_\n"

			if private && len(privateChannel) > 0 {
				// private message and private channel defined
				_, err := s.ChannelMessageSend(privateChannel, discordMsg)
				if err != nil {
					log.Println("mesh failed to send")
				}
			} else if !private && len(publicChannel) > 0 {
				// public and there's a public channel defined
				_, err := s.ChannelMessageSend(publicChannel, discordMsg)
				if err != nil {
					log.Println("mesh failed to send")
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}
