package main

import (
	"context"
	"encoding/json"
	"fmt"

	spvchannels "github.com/libsv/go-spvchannels"
)

func main() {
	cfg := spvchannels.ClientConfig{
		Insecure: true,
		BaseURL:  "localhost:5010",
		Version:  "v1",
		User:     "dev",
		Passwd:   "dev",
		Token:    "",
	}

	client := spvchannels.NewClient(cfg)
	/*
		r := spvchannels.GetChannelRequest{
			AccountId: "1",
			ChannelId: "2vkapEui-Cfb3tY7l9FFviRjpsNGa0Iv4kFEHYoMWJdl4f9PSlvurjOCnTBzH1r_C8VUuvQsn-0NsO0Q2bKGUA",
		}
	*/
	r := spvchannels.DeleteTokenRequest{
		AccountId: "1",
		ChannelId: "IZjYsiWpRwg6qi1zCWUkUs2eDEN16YovfQeZMX1gD-FHEjC68kCZle_5NcwekCWae5ecvyBqdNqzql6ler6bxA",
		TokenId:   "1",
	}

	reply, err := client.DeleteToken(context.Background(), r)
	if err != nil {
		panic("foo")
	}

	bReply, _ := json.Marshal(reply)
	fmt.Println(string(bReply))
}
