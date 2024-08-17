package streamer

import (
	web_api "github.com/fr0ster/turbo-restler/web_api"
	web_stream "github.com/fr0ster/turbo-restler/web_stream"
)

type (
	Stream struct {
		maintainWebApi *web_api.WebApi
		symbol         string
		wsHost         web_api.WsHost
		wsPath         web_api.WsPath
		low_stream     *web_stream.WebStream
	}
)
