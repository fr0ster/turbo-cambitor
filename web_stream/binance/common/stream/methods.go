package streamer

import web_stream "github.com/fr0ster/turbo-restler/web_stream"

func (stream *Stream) Start() (err error) {
	return stream.low_stream.Start()
}

func (stream *Stream) Stop() {
	stream.low_stream.Stop()
}

func (stream *Stream) Close() {
	stream.low_stream.Close()
}

func (stream *Stream) SetHandler(handler web_stream.WsHandler) *Stream {
	stream.low_stream.SetDefaultHandler(handler)
	return stream
}

func (stream *Stream) SetErrHandler(errHandler web_stream.ErrHandler) *Stream {
	stream.low_stream.SetErrHandler(errHandler)
	return stream
}

func (stream *Stream) AddSubscription(handler web_stream.WsHandler, subscription string) {
	stream.low_stream.AddHandler(subscription, handler)
	stream.low_stream.Subscribe(subscription)
}

func (stream *Stream) ListOfSubscriptions(handler web_stream.WsHandler) error {
	return stream.low_stream.ListOfSubscriptions(handler)
}

func (stream *Stream) RemoveSubscriptions(subscription string) {
	stream.low_stream.RemoveHandler(subscription)
	stream.low_stream.Unsubscribe(subscription)
}
