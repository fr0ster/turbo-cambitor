package streamer

import (
	web_api "github.com/fr0ster/turbo-restler/web_api"
	web_stream "github.com/fr0ster/turbo-restler/web_stream"
	"github.com/sirupsen/logrus"
)

func (stream *Stream) SetSymbol(symbol string) *Stream {
	stream.symbol = symbol
	return stream
}

func New(
	host web_api.WsHost,
	path web_api.WsPath,
	scheme ...web_api.WsScheme) *Stream {
	stream, err := web_stream.New(host, path, scheme...)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	wa, err := web_api.New(host, path)
	if err != nil {
		panic(err)
	}
	return &Stream{
		maintainWebApi: wa,
		wsHost:         host,
		wsPath:         path,
		low_stream:     stream,
	}
}
