package grpcutil

import (
	"errors"
	"strconv"

	"context"

	"github.com/ravlio/wallet/pkg/errutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ServerErrorInterceptor(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	response, err := handler(ctx, request)

	if err != nil {
		_, isOwn := err.(errutil.Error)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		if isOwn || status.Code(err) == codes.Unknown { // Here is own error or non-grpc error
			var code string

			if isOwn {
				code = strconv.Itoa(err.(errutil.Error).GetCode())
			} else {
				code = strconv.Itoa(500)
			}

			md = metadata.Join(metadata.Pairs("error-code", code), md)
			if msg := err.Error(); msg != "" {
				md = metadata.Join(metadata.Pairs("error-message-bin", msg), md)
			}

			if err := grpc.SetHeader(ctx, md); err != nil {
				return response, status.Errorf(codes.DataLoss, "metadata missing")

			}
		}
	}

	return response, err
}

func ClientErrorInterceptor(ctx context.Context, method string, request interface{},
	response interface{}, connection *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	var header metadata.MD

	opts = append(opts, grpc.Header(&header))
	err := invoker(ctx, method, request, response, connection, opts...)

	if header != nil {
		if code, ok := header["error-code"]; ok {
			var m string

			if msg, ok := header["error-message-bin"]; ok {
				m = msg[0]
			}

			ci, e := strconv.Atoi(code[0])
			if e != nil {
				return errors.New("invalid metadata")
			}

			return errutil.NewFromCode(ci, m)
		}
	}

	return err
}
