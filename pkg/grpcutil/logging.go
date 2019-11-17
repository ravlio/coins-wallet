package grpcutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ravlio/wallet/pkg/errutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Log struct {
	req         json.Marshaler
	svc, method string
	st          time.Time
}

func LogRequest(svc, method string, req json.Marshaler) *Log {
	return &Log{svc: svc, method: method, req: req, st: time.Now()}
}

func (l *Log) LogResponse(ctx context.Context, resp json.Marshaler, err error) {
	var event *zerolog.Event
	var localErr error
	var reqb, respb []byte

	stime := time.Now()

	if l.req != nil {
		if v, ok := l.req.(json.Marshaler); ok {
			reqb, localErr = v.MarshalJSON()
			if localErr != nil {
				err = fmt.Errorf("marshal error: %w", localErr)
			}
		}
	}

	if resp != nil {
		if v, ok := l.req.(json.Marshaler); ok {
			respb, localErr = v.MarshalJSON()
			if localErr != nil {
				err = fmt.Errorf("marshal error: %w", localErr)
			}
		}
	}

	if err != nil {
		event = log.Error()

		if ownErr, ok := err.(errutil.Error); ok {
			event = event.Int("errorCode", ownErr.GetCode())
		} else {
			event = event.Int("errorCode", http.StatusInternalServerError)
		}
		event = event.Str("error", err.Error())
	} else {
		event = log.Info()
	}

	event = event.Str("service", l.svc)
	event = event.Str("method", l.method)

	if reqb != nil {
		event = event.Str("request", string(reqb))
	}

	if respb != nil {
		event = event.Str("response", string(respb))
	}
	event = event.Dur("took", time.Since(stime))

	event.Msg("handle request")
}
