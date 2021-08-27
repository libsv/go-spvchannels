// +build integration

package integration

import (
	"context"

	"testing"

	spv "github.com/libsv/go-spvchannels"
	"github.com/stretchr/testify/assert"
)

func TestChannelIntegration(t *testing.T) {

	cfg := spv.ClientConfig{
		Insecure: true,
		BaseURL:  "localhost:5010",
		Version:  "v1",
		User:     duser,
		Passwd:   dpassword,
		Token:    "",
	}

	t.Run("CreateChannel", func(t *testing.T) {

		client := spv.NewClient(cfg)

		r := spv.CreateChannelRequest{
			AccountId:   accountid,
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
		assert.NotNil(t, reply)
		assert.NoError(t, err)
		assert.Equal(t, len(reply.AccessTokens), 1)
		assert.NotEmpty(t, reply.Id)
		assert.NotEmpty(t, reply.AccessTokens[0].Token)
	})
}
