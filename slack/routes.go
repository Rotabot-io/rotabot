package slack

import (
	"github.com/rotabot-io/rotabot/gen/http/slack/server"
	gen "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/lib/codecs"
	"github.com/rotabot-io/rotabot/lib/errors"
	goahttp "goa.design/goa/v3/http"
)

func NewServer(mux goahttp.Muxer, service gen.Service) *server.Server {
	endpoints := gen.NewEndpoints(service)

	return server.New(
		endpoints,
		mux,
		codecs.RequestDecoderWithLogs,
		codecs.ResponseEncoderWithLogs,
		errors.ErrorHandler(),
		nil,
	)
}
