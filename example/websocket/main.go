package main

import (
	"context"
	"encoding/json"
	"fmt"

	spv "github.com/libsv/go-spvchannels"
)

var channelid string = "b1j-Vd94XrU9NJlnrtQPfkzOgFLxun5oWZLUhXfvnZk2cekKe4QY7YKh_hbXivAroApDtVn3pmmOo848R6BhAw"
var token string = "hqaDcOY-3svYqZv5RXID7AphKp9Bm8obQ_74K7mFLjcjq_Bw-Vwng6Q0q7PvJqhKawikfmd0Kr2OpYFKpFrKcg"

// PullUnreadMessages pull the notified unread messages, mark them as read
func PullUnreadMessages(t int, msg []byte, err error) error {
	// If notification error, then return the error
	if err != nil {
		return err
	}

	// Pull unread messages
	cfg := spv.ClientConfig{
		Insecure: true,
		BaseURL:  "localhost:5010",
		Version:  "v1",
		User:     "dev",
		Passwd:   "dev",
		Token:    token,
	}

	restClient := spv.NewClient(cfg)

	r := spv.MessagesRequest{
		ChannelId: channelid,
		UnRead:    true,
	}

	unreadMsg, err := restClient.Messages(context.Background(), r)
	if err != nil {
		return fmt.Errorf("unable to read new messages : %v", err)
	}

	for _, msg := range *unreadMsg {
		msg_seq := msg.Sequence
		r2 := spv.MessageMarkRequest{
			ChannelId: channelid,
			Sequence:  msg_seq,
			Older:     true,
			Read:      true,
		}

		_, err := restClient.MessageMark(context.Background(), r2)
		if err != nil {
			return fmt.Errorf("unable mark message as read : %v", err)
		}
	}

	bReply, _ := json.MarshalIndent(unreadMsg, "", "    ")
	fmt.Println("\nNew unread messages ===================")
	fmt.Println(string(bReply))

	return nil
}

// This program run a websocket notification listener
// Anytime a new (unread) message is notified, it pull the new messages, mark them as read
func main() {

	cfg := spv.WSConfig{
		Insecure:  true,
		BaseURL:   "localhost:5010",
		Version:   "v1",
		ChannelId: channelid,
		Token:     token,
	}

	ws := spv.NewWSClient(
		cfg,
		PullUnreadMessages,
		10,
	)

	err := ws.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println("Exit Success : total processed ", ws.NbNotified(), " messages")
}
