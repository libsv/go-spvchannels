package main

import (
	"context"
	"encoding/json"
	"fmt"

	spv "github.com/libsv/go-spvchannels"
)

func main() {
	cfg := spv.ClientConfig{
		Insecure: true,
		BaseURL:  "localhost:5010",
		Version:  "v1",
		User:     "dev",
		Passwd:   "dev",
		Token:    "",
	}

	client := spv.NewClient(cfg)

	r := spv.CreateChannelRequest{
		AccountId:   "1",
		PublicRead:  true,
		PublicWrite: true,
		Sequenced:   true,
		Retention: struct {
			MinAgeDays int  "json:\"min_age_days\""
			MaxAgeDays int  "json:\"max_age_days\""
			AutoPrune  bool "json:\"auto_prune\""
		}{
			MinAgeDays: 0,
			MaxAgeDays: 99999,
			AutoPrune:  true,
		},
	}

	reply, err := client.CreateChannel(context.Background(), r)
	if err != nil {
		panic("Problem with the request")
	}

	bReply, _ := json.MarshalIndent(reply, "", "    ")
	fmt.Println(string(bReply))
}
