package codecs

import (
	"context"
	"fmt"
	stdhttp "net/http"

	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
	"goa.design/goa/v3/http"
)

type encoder struct {
	encoder http.Encoder
	logger  *zap.Logger
	kind    string
}

func ResponseEncoderWithLogs(ctx context.Context, r stdhttp.ResponseWriter) http.Encoder {
	return &encoder{
		encoder: http.ResponseEncoder(ctx, r),
		logger:  zapctx.Logger(ctx),
		kind:    "response",
	}
}

func (e encoder) Encode(v interface{}) error {
	err := e.encoder.Encode(v)
	if err != nil {
		e.logger.Error(fmt.Sprintf("failed to decode %s", e.kind), zap.Error(err))
	}

	return err
}
