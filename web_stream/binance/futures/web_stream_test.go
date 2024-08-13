package futures_web_stream_test

import (
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/futures"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	timeOut = 1 * time.Second
)

// Mock handler for WebSocket messages
func mockHandler(message *simplejson.Json) {
	logrus.Infof("Received message: %+v", message)
}

// Mock error handler for WebSocket errors
func mockErrHandler(err error) {
	logrus.Errorf("Error: %v", err)
}

func TestStartBinanceStreamer(t *testing.T) {
	// Test 2: Binance stream
	stream := web_stream.New("BTCUSDT", true)
	assert.NotNil(t, stream)

	doneC, stopC, err :=
		stream.Klines("1m").
			SetHandler(
				func(message *simplejson.Json) {
					mockHandler(message)
				}).
			SetErrHandler(
				func(err error) {
					mockErrHandler(err)
				}).
			Start()
	assert.NoError(t, err)
	assert.NotNil(t, doneC)
	assert.NotNil(t, stopC)

	// Stop the streamer after some time
	time.Sleep(timeOut)
	close(stopC)
}
