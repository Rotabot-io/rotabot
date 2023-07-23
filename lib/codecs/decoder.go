package codecs

import (
	"fmt"
	stdhttp "net/http"

	"github.com/ajg/form"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
	"goa.design/goa/v3/http"
)

type decoder struct {
	decoder http.Decoder
	logger  *zap.Logger
	kind    string
}

func RequestDecoderWithLogs(r *stdhttp.Request) http.Decoder {
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		return &decoder{
			decoder: form.NewDecoder(r.Body),
			logger:  zapctx.Logger(r.Context()),
			kind:    "request",
		}
	}
	return &decoder{
		decoder: http.RequestDecoder(r),
		logger:  zapctx.Logger(r.Context()),
		kind:    "request",
	}
}

func (d decoder) Decode(v interface{}) error {
	err := d.decoder.Decode(v)
	if err != nil {
		d.logger.Error(fmt.Sprintf("failed to decode %s", d.kind), zap.Error(err))
	}

	return err
}
