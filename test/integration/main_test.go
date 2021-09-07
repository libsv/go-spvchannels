// +build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	spv "github.com/libsv/go-spvchannels"
)

var duser = "dev"
var dpassword = "dev"
var accountid = ""

func getClientConfig() spv.ClientConfig {
	return spv.ClientConfig{
		Insecure: true,
		BaseURL:  "localhost:5010",
		Version:  "v1",
		User:     duser,
		Passwd:   dpassword,
		Token:    "",
	}
}

func createChannel(cfg spv.ClientConfig) (*spv.ChannelCreateReply, error) {
	client := spv.NewClient(cfg)

	r := spv.ChannelCreateRequest{
		AccountID:   accountid,
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

	reply, err := client.ChannelCreate(context.Background(), r)
	return reply, err
}

func getChannel(cfg spv.ClientConfig, accountid string, channelid string) (*spv.ChannelReply, error) {
	client := spv.NewClient(cfg)

	r := spv.ChannelRequest{
		AccountID: accountid,
		ChannelID: channelid,
	}

	reply, err := client.Channel(context.Background(), r)
	return reply, err
}

func getChannels(cfg spv.ClientConfig, accountid string) (*spv.ChannelsReply, error) {
	client := spv.NewClient(cfg)

	r := spv.ChannelsRequest{
		AccountID: accountid,
	}

	reply, err := client.Channels(context.Background(), r)
	return reply, err
}

func setup() error {

	cmdcreateUser := exec.Command("docker", "exec", "spvchannels", "./SPVChannels.API.Rest", "-createaccount", "spvchannels_dev", duser, dpassword)
	out, err := cmdcreateUser.CombinedOutput()
	if err != nil {
		return err
	}

	parts := strings.Split(strings.TrimSpace(string(out)), ":")
	if len(parts) != 2 {
		return errors.New("Issue with creating account command")
	}
	accountid = parts[1]
	return nil
}

func teardown() error {
	// TODO teardown : clear spv database inside spvchannel_db
	return nil
}

func TestMain(m *testing.M) {

	serr := setup()

	if serr == nil {

		code := m.Run()
		terr := teardown()

		if terr == nil {
			os.Exit(code)
		} else {
			fmt.Printf("%v", terr)
		}

	} else {
		fmt.Printf("%v", serr)
	}

	os.Exit(1)
}
