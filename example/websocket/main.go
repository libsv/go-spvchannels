package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"

	spv "github.com/libsv/go-spvchannels"
)

var channelid = "cc33DR4-U1ZJnvKvMChUcikxck5ANY1XhIbkIz5YPXJIZCS--mDIUc9Ot5HbLRSBIZuFFtNZzKIMz8StV46-cw"
var tok = "T_udxsz-1sE9RPeBmEgYNNTtFaQEW204ETf3DzcpwQydIz_Gvb7X3gB6rPypkQmIH5Fl9_cfiZdk6XD3uzTmoQ"

// PullUnreadMessages handle new coming notification. It does :
//  - pull the notified unread message's content
//  - mark them as read
//  - print the message content to the log
//  - close the websocket if the message content contain 'close' or 'Close'.
func PullUnreadMessages(c context.Context, t int, msg []byte, err error) error {
	ctx, cancelfn := context.WithTimeout(c, time.Second)
	defer cancelfn()
	// If notification error, then return the error
	if err != nil {
		return err
	}

	// Pull unread messages
	restClient := spv.NewClient(
		spv.WithBaseURL("localhost:5010"),
		spv.WithVersion("v1"),
		spv.WithUser("dev"),
		spv.WithPassword("dev"),
		spv.WithToken(tok),
		spv.WithInsecure(),
	)

	r := spv.MessagesRequest{
		ChannelID: channelid,
		UnRead:    true,
	}

	unreadMsg, err := restClient.Messages(ctx, r)
	if err != nil {
		return fmt.Errorf("unable to read new messages : %w", err)
	}

	for _, msg := range unreadMsg {

		// Mark the unread message as read
		msgSeq := msg.Sequence
		r2 := spv.MessageMarkRequest{
			ChannelID: channelid,
			Sequence:  msgSeq,
			Older:     true,
			Read:      true,
		}

		err := restClient.MessageMark(context.Background(), r2)
		if err != nil {
			return fmt.Errorf("unable mark message as read : %w", err)
		}

		// Print the unread message content
		binMsg, _ := b64.StdEncoding.DecodeString(msg.Payload)
		msgStr := string(binMsg)
		fmt.Println(msgStr)

		// If received a message ordering to close, then return the close error, so the websocket client know to close
		if strings.Contains(msgStr, "Close") || strings.Contains(msgStr, "close") {
			fmt.Println("Received closing message")
			return spv.ErrWSClose{}
		}

	}
	return nil
}

// This program run a websocket notification listener
// Anytime a new (unread) message is notified, it pull the new messages, mark them as read
func main() {

	ws, err := spv.NewWSClient(
		spv.WithBaseURL("localhost:5010"),
		spv.WithVersion("v1"),
		spv.WithChannelID(channelid),
		spv.WithToken(tok),
		spv.WithInsecure(),
		spv.WithWebsocketCallBack(PullUnreadMessages),
	)

	if err != nil {
		panic(err)
	}

	ws.Run()
	defer ws.Close()

	fmt.Println("Exit Success")
}
