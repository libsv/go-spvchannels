// +build integration

package integration

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	spv "github.com/libsv/go-spvchannels"
	"github.com/stretchr/testify/assert"
)

// getChannelTokens create a new channel with 2 tokens
// Return chanelid, token1, token2 and error if any
func getChannelTokens() (string, string, string, error) {
	restClient := getRestClient()

	// Create a new channel
	replyCreateChannel, err := createChannel(restClient)
	if err != nil {
		return "", "", "", err
	}
	channelid := replyCreateChannel.ID
	token1 := replyCreateChannel.AccessTokens[0].Token

	// Create and additional token on the channel
	r := spv.TokenCreateRequest{
		AccountID:   accountid,
		ChannelID:   replyCreateChannel.ID,
		Description: "Create Second Token for websocket test",
		CanRead:     true,
		CanWrite:    true,
	}

	reply, err := restClient.TokenCreate(context.Background(), r)
	if err != nil {
		return "", "", "", err
	}
	token2 := reply.Token
	return channelid, token1, token2, nil
}

// TestWebsocket run 2 goroutines :
//     One open a websocket client, listen and process messages. The process exit when it fully received N messages
//     Other keep sending messages to spv channel server, exit only when the first one received enough messages
func TestWebsocket(t *testing.T) {
	channelid, token1, token2, err := getChannelTokens()
	assert.NoError(t, err)
	totalSend := uint64(10)
	totalRecv := uint64(0)
	tries := 0

	ws, err := spv.NewWSClient(
		spv.WithBaseURL(baseURL),
		spv.WithVersion(version),
		spv.WithToken(token1),
		spv.WithChannelID(channelid),
		spv.WithWebsocketCallBack(func(ctx context.Context, t int, msg []byte, err error) error {
			if err != nil {
				return nil
			}
			atomic.AddUint64(&totalRecv, 1)
			return nil
		}),
		spv.WithInsecure(),
	)
	assert.NoError(t, err)

	// Websocket client routine ---------------------------------------------------------
	go ws.Run()
	defer ws.Close()

	var wg sync.WaitGroup

	// Message writer client routine ------------------------------------------------------
	wg.Add(1)
	go func() {
		defer wg.Done()
		restClient := spv.NewClient(
			spv.WithBaseURL(baseURL),
			spv.WithVersion(version),
			spv.WithUser(duser),
			spv.WithPassword(dpassword),
			spv.WithToken(token2),
			spv.WithInsecure(),
		)

		for i := 0; uint64(i) < atomic.LoadUint64(&totalSend); i++ {
			r := spv.MessageWriteRequest{
				ChannelID: channelid,
				Message:   "Some random message",
			}

			_, err := restClient.MessageWrite(context.Background(), r)
			assert.NoError(t, err)
		}
	}()
	wg.Wait()
	for atomic.LoadUint64(&totalRecv) < atomic.LoadUint64(&totalSend) && tries <= 10 {
		time.Sleep(500 * time.Millisecond)
		tries++
	}

	assert.Equal(t, int(atomic.LoadUint64(&totalSend)), int(atomic.LoadUint64(&totalRecv)))
}
