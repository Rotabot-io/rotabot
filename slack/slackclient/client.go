//go:generate sh -c "interfacer -for github.com/slack-go/slack.Client -as slackclient.SlackClient | grep -v _search > client_interface.go"
//go:generate mockgen -package=mock_slackclient -destination=mock_slackclient/client.go -source=client_interface.go . SlackClient
package slackclient

import (
	"context"
	"os"

	"github.com/rotabot-io/rotabot/slack/slackclient/mock_slackclient"

	"github.com/slack-go/slack"
	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
)

type contextKey string

var clientContextKey contextKey = "slackclient.client"

// WithClient returns a new context.Context where calls to ClientFor
// will return the given client, rather than creating a new one.
//
// This can be used for testing, when we want any calls to build a
// client to return a mock, rather than a real implementation.
func WithClient(ctx context.Context, client SlackClient) context.Context {
	return context.WithValue(ctx, clientContextKey, client)
}

// ClientFor will return a SlackClient, either from the client stored in
// the context or by generating a new client for the given organisation.
func ClientFor(ctx context.Context, organisationID string) (SlackClient, error) {
	client, ok := ctx.Value(clientContextKey).(SlackClient)
	if ok {
		return client, nil
	}

	c := getCredentials(ctx, organisationID)
	return slack.New(c.SlackAccessToken), nil
}

type Credentials struct {
	SlackAccessToken string
}

func getCredentials(_ context.Context, _ string) *Credentials {
	return &Credentials{
		SlackAccessToken: os.Getenv("SLACK_CLIENT_SECRET"),
	}
}

// MockSlackClient is used in tests to generate a mock client and stash
// it into a context, ensuring all code will use the mock client instead
// of reaching out into the real world.
//
// Example is:
//
//	Describe("subject", func() {
//	  slackclient.MockSlackClient(&ctx, &sc, nil)
//	})
func MockSlackClient(ctxPtr *context.Context, scPtr **mock_slackclient.MockSlackClient, ctrlPtr **gomock.Controller) {
	var ctrl *gomock.Controller
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		if ctrlPtr != nil {
			*ctrlPtr = ctrl
		}

		*scPtr = mock_slackclient.NewMockSlackClient(ctrl)
		*ctxPtr = WithClient(*ctxPtr, *scPtr)
	})
}
