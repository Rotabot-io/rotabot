package slack

import (
	"github.com/rotabot-io/rotabot/gen/http/slack/server"
	"github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/lib/codecs"
	"github.com/rotabot-io/rotabot/lib/goaerrors"
	goahttp "goa.design/goa/v3/http"
)

func NewServer(mux goahttp.Muxer, service slack.Service) *server.Server {
	endpoints := slack.NewEndpoints(service)

	return server.New(
		endpoints,
		mux,
		codecs.RequestDecoderWithLogs,
		codecs.ResponseEncoderWithLogs,
		goaerrors.ErrorHandler(),
		nil,
	)
}
