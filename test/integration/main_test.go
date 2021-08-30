// +build integration

package integration

import (
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
	fmt.Println("TODO teardown : clear spv database inside spvchannel_db")

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
