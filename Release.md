# Release Notes for Turbo-Cambitor

## v0.2.34

### Release Date: 2024-09-17

#### RestAPI Changes:
- **NewOrder**: Fixed endpoints for `binance/futures` and `binance/spot`.
- **QueryOrder**: Fixed endpoints for `binance/futures` and `binance/spot`.
- **CancelOrder**: Fixed endpoints for `binance/futures` and `binance/spot`.
- **ModifyOrder**: Added new function to support order modification for `binance/futures` and `binance/spot`.

---

## v0.2.33

### Release Date: 2024-09-12

#### RestAPI:
- **Updated**: Migrated to v3 endpoints for `Account Balance` and `PositionRisk`.

This release includes updates to the RestAPI component of Turbo-Cambitor. The `Account Balance` and `PositionRisk` functions have been migrated to use the v3 endpoints, ensuring improved performance and compatibility with the latest API standards.

---

## v0.2.32

### Release Date: 2024-09-12

#### RestAPI:
- **Added**: Function for `PositionRisk`.

#### WebApi:
- **Removed**: Functions that use deprecated endpoints.

This release includes updates to both the RestAPI and WebApi components of Turbo-Cambitor. The new `PositionRisk` function has been added to the RestAPI, providing enhanced risk management capabilities. Additionally, outdated functions that relied on deprecated endpoints have been removed from the WebApi to ensure compatibility and maintainability.

---

## v0.2.31

### Release Date: 2024-08-24

### Chore
- Upgraded to `turbo-restler` v0.2.27 and `turbo-signer` v0.1.7.

### Refactor
- For `RestAPI`:
  - Extracted common code from `RestApi` into a separate `common_rest_api.RestApiWrapper`.
  - Renamed `Request` to `RequestBuilder` and moved the call to `CallRestAPI` out of the `Do` method to `common_rest_api.RestApiWrapper`.
- For `WebApi`:
  - Renamed `Request` to `RequestBuilder`.

---

## v0.2.30

### Release Date: 2024-08-22

### Chore
- Removed unnecessary code.

---

## v0.2.29

### Release Date: 2024-08-21

### Chore
- Upgraded to `turbo-restler` v0.2.26 and `turbo-signer` v0.1.6.

### Fix
- Fixed implementation of `RestAPI` functions.
- Refactored code related to the change of parameter in `turbo_restler` `CallRestAPI` to `*http.Request`.

---

## v0.2.28

### Release Date: 2024-08-19

### Changes
- **Turbo-Cambitor**:
  - **Update**: Upgraded to `turbo-restler` v0.2.19. Version v0.2.27 was uncorrected.

---

## v0.2.27

### Release Date: 2024-08-19

### Changes
- **Turbo-Cambitor**:
  - **Update**: Upgraded to `turbo-restler` v0.2.19.

---

## v0.2.26

### Release Date: 2024-08-19

### Changes
- **Turbo-Cambitor**:
  - **Update**: Upgraded to `turbo-restler` v0.2.18.
  - **WebApi**:
    - Retained only the synchronous `Call` method.
  - **WebStream**:
    - Utilized the asynchronous `Call` method.

---

## v0.2.25

### Release Date: 2024-08-18

### Changes
- **Turbo-Cambitor**:
  - **Update**: Upgraded to `turbo-restler` v0.2.17.
  - **WebStream**:
    - Added `Lock` and `Unlock` functions.

---

## v0.2.24

### Release Date: 2024-08-18

### Changes
- **Turbo-Cambitor**:
  - **Update**: Added `Lock` and `Unlock` functions to the `futures.WebApi` and `spot.WebApi` interfaces.

---

## v0.2.23

### Release Date: 2024-08-18

### Changes
- **Turbo-Cambitor**:
  - **Update**: Upgraded to `turbo-restler` v0.2.16.
  - **Refactor**:
    - Moved the `Do` method from `Request` to `WebApi` and renamed `WebApi` to `WebApiWrapper`.
    - Renamed `Stream` to `StreamWrapper`.

---

## v0.2.22

### Release Date: 2024-08-18

### Changes
- **Turbo-Cambitor**:
  - **Refactor**: WebSocket connection handling in Binance stream for improved reliability and performance.
    - This commit refactors the WebSocket connection handling in the Binance stream module. The changes aim to enhance the reliability and performance of the WebSocket connections, ensuring more stable and efficient data streaming.
  - **Update**: Updated `fr0ster/turbo-restler` dependency to v0.2.13.

---

## v0.2.22

### Release Date: 2024-08-17

### Changes
- **Turbo-Cambitor**:
  - **Refactor**: WebSocket connection handling in Binance stream for improved reliability and performance.
    - This commit refactors the WebSocket connection handling in the Binance stream module. The changes aim to enhance the reliability and performance of the WebSocket connections, ensuring more stable and efficient data streaming.

---

## v0.2.21

### Release Date: 2024-08-17

### Changes
- **Turbo-Cambitor**:
  - **Refactor**: WebSocket connection handling in Binance stream to simplify code and improve reliability and performance.
    - This commit refactors the WebSocket connection handling in the Binance stream module. The changes aim to simplify the codebase and enhance the reliability and performance of the WebSocket connections, ensuring more stable and efficient data streaming.

---

## v0.2.20

### Release Date: 2024-08-17

### Changes
- **Turbo-Cambitor**:
  - **Refactor**: WebSocket connection handling in Binance stream to simplify code and improve reliability and performance.
    - This commit refactors the WebSocket connection handling in the Binance stream module. The changes aim to simplify the codebase and enhance the reliability and performance of the WebSocket connections, ensuring more stable and efficient data streaming.

---

## v0.2.19

### Release Date: 2024-08-17

### Changes
- **Turbo-Cambitor**:
  - **Refactor**: WebSocket connection handling in Binance stream to improve reliability and performance.
    - This commit refactors the WebSocket connection handling in the Binance stream module. The changes aim to enhance the reliability and performance of the WebSocket connections, ensuring more stable and efficient data streaming.
---

## v0.2.18

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - For `SetHandler(handler WsHandler) *Stream`, removed default checks. Now all checks must be implemented in user-defined handlers.

---

## v0.2.17

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Updated `fr0ster/turbo-restler` dependency to v0.2.8.
  - Removed timeout functionality from `WebStream.Start` to simplify connection management.


---

## v0.2.16

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Merge commit '65bd5729c070933d4751d8acfe8fed017f68f0aa'

---

## v0.2.15

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Refactored WebSocket connection handling in Binance stream for improved reliability and performance.

---

## v0.2.14

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Refactored WebSocket connection handling in Binance stream for improved reliability and performance.

---

## v0.2.13

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Refactored WebSocket connection handling in Binance stream to improve reliability and performance.

---

## v0.2.12

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Refactored web stream functions to improve code maintainability and performance.

### Example Usage
#### Refactored Web Stream Functions
```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/fr0ster/turbo-cambitor/stream"
    "github.com/bitly/go-simplejson"
)

var (
	interrupt chan os.Signal = make(chan os.Signal, 1)
	quit                     = make(chan struct{})
)

// Mock handler for WebSocket messages
func mockHandler(message *simplejson.Json) {
	logrus.Infof("Received message: %+v", message)
}

// Mock error handler for WebSocket errors
func mockErrHandler(err error) {
	logrus.Errorf("Error: %v", err)
}

func main() {
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt
		// Закриваємо канал quit, щоб сповістити всі горутина про необхідність завершення
		close(quit)
	}()
	stream := web_stream.
		New(true).
		Stream().
		SetSymbol("BTCUSDT").
		SetHandler(mockHandler).
		SetErrHandler(mockErrHandler).
		SetTimeOut(5 * time.Second)
	_ = stream.Start()
	assert.NoError(t, err)
	response, _ := stream.ListSubscriptions()
	logrus.Infof("Subscriptions: %+v", response)
	stream.AddSubscriptions(mockHandler, "btcusdt@aggTrade")
	response, _ = stream.ListSubscriptions()
	logrus.Infof("Subscriptions: %+v", response)
	stream.RemoveSubscriptions("btcusdt@aggTrade")
	time.Sleep(timeOut)
	response, _ = stream.ListSubscriptions()
	logrus.Infof("Subscriptions: %+v", response)
	stream.Stop()
}
```

---

## v0.2.11

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Set ping handler for WebSocket connection in Binance stream to maintain connection stability.

---

## v0.2.10

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Refactored web stream functions to improve code maintainability and performance.

### Example Usage
#### Refactored Web Stream Functions
```go
package main

import (
    "fmt"
    "log"
    "github.com/fr0ster/turbo-cambitor/stream"
    "github.com/bitly/go-simplejson"
)

// Mock handler for WebSocket messages
func mockHandler(message *simplejson.Json) {
	logrus.Infof("Received message: %+v", message)
}

// Mock error handler for WebSocket errors
func mockErrHandler(err error) {
	logrus.Errorf("Error: %v", err)
}

func main() {
	stream := web_stream.
		New(true).
		Stream().
		SetSymbol("BTCUSDT").
		SetHandler(mockHandler).
		SetErrHandler(mockErrHandler).
		SetTimeOut(5 * time.Second)
	_ = stream.Start()
	response, _ := stream.ListSubscriptions()
	logrus.Infof("Subscriptions: %+v", response)
	stream.AddSubscriptions("btcusdt@aggTrade")
	response, _ = stream.ListSubscriptions()
	logrus.Infof("Subscriptions: %+v", response)
	stream.RemoveSubscriptions("btcusdt@aggTrade")
	time.Sleep(5 * time.Second)
	response, _ = stream.ListSubscriptions()
	logrus.Infof("Subscriptions: %+v", response)
	stream.Stop()
}
```

---

## v0.2.9

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Resolved issues caused by the removal of the following function implementations:
    - `Subscribe([]string) (response *simplejson.Json, err error)`
    - `ListOfSubscriptions() (response *simplejson.Json, err error)`
    - `Unsubscribe([]string) (response *simplejson.Json, err error)`
  - These functions are still under redevelopment and will be reintroduced in a future release.

---

## v0.2.8

### Release Date: 2024-08-16

### Changes
- **Turbo-cambitor**:
  - Removed implementations of the following functions:
    - `Subscribe([]string) (response *simplejson.Json, err error)`
    - `ListOfSubscriptions() (response *simplejson.Json, err error)`
    - `Unsubscribe([]string) (response *simplejson.Json, err error)`
  - These functions are now under redevelopment and will be reintroduced in a future release.

---

## v0.2.7

### Release Date: 2024-08-15

### Changes
- **Turbo-cambitor**:
  - Removed implementations of the following functions from `futures_web_stream.WebStream` and `spot_web_stream.WebStream` and moved them to `Stream`:
    - `Subscribe([]string) (response *simplejson.Json, err error)`
    - `ListOfSubscriptions() (response *simplejson.Json, err error)`
    - `Unsubscribe([]string) (response *simplejson.Json, err error)`
  - Upgraded to `github.com/fr0ster/turbo-restler` version `v0.2.6`, which includes new features and improvements.

### Example Usage
#### Stream Functions
```go
package main

import (
  web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/futures"
	"github.com/sirupsen/logrus"
)

func main() {
  stream := web_stream.New(true).DynamicStream()
  logrus.Infof("Stream: %v", stream)
  response, err := stream.Subscribe([]string{"btcusdt@aggTrade"})
  if err != nil {
    logrus.Fatalf("Error: %v", err)
  }
  logrus.Infof("Response: %v", response)
  response, err = stream.ListOfSubscriptions()
  if err != nil {
    logrus.Fatalf("Error: %v", err)
  }
  logrus.Infof("Response: %v", response)
  response, err = stream.Unsubscribe([]string{"btcusdt@aggTrade"})
  if err != nil {
    logrus.Fatalf("Error: %v", err)
  }
  logrus.Infof("Response: %v", response)
}
```

---

## v0.2.6

### Release Date: 2024-08-15

### Changes
- **New Functions Implemented**:
  - `Subscribe([]string) (response *simplejson.Json, err error)`: Allows subscribing to a list of topics.
  - `ListOfSubscriptions() (response *simplejson.Json, err error)`: Retrieves the current list of subscriptions.
  - `Unsubscribe([]string) (response *simplejson.Json, err error)`: Allows unsubscribing from a list of topics.

### Example Usage
#### Subscribe
```go
package main

import (
    "fmt"
    "log"
    "github.com/fr0ster/turbo-cambitor/web_stream"
    "github.com/bitly/go-simplejson"
)

func main() {
    stream := web_stream.NewWebStream("https://api.example.com")

    topics := []string{"topic1", "topic2"}
    response, err := stream.Subscribe(topics)
    if err != nil {
        log.Fatal("Subscribe error:", err)
    }

    fmt.Println("Subscribe response:", response)
}
```

---

## v0.2.5

### Release Date: 2024-08-15

### Changes
- **Updated Dependencies**:
  - Upgraded to `github.com/fr0ster/turbo-signer` version `v0.1.4`, which includes new features and improvements.
  - Upgraded to `github.com/fr0ster/turbo-restler` version `v0.2.4`, which includes new features and improvements.

---

## v0.2.4

### Release Date: 2024-08-15

### Changes
- **Fixed rq.params Call**: Corrected the call to `rq.params, _ = rq.sign.SignParameters(rq.params)` to ensure proper signing of parameters.

---

## v0.2.3

### Release Date: 2024-08-15

### Changes
- **Updated Dependencies**:
  - Upgraded to `github.com/fr0ster/turbo-signer` version `v0.1.2`, which includes new features and improvements.
  - Upgraded to `github.com/fr0ster/turbo-restler` version `v0.2.2`, which includes new features and improvements.

---

## v0.2.2

### Release Date: 2024-08-14

### Changes
- **Dependency Update**: Updated `github.com/fr0ster/turbo-restler` to version `v0.2.1` for improved stability and new features.

---

## v0.2.1

### Release Date: 2024-08-14

### Introduction
Turbo-Cambitor is a wrapper library for RESTful and WebSocket APIs of cryptocurrency exchanges, initially designed for Binance. This library aims to simplify the interaction with exchange APIs, providing a more intuitive and efficient way to integrate trading functionalities into your applications.

---

Thank you for using Turbo-Cambitor! If you have any questions or feedback, please contact our support team.