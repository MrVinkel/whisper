package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type CallbackServer struct {
	done        chan bool
	callBackURL *url.URL
	server      http.Server
}

func StartCallbackServer(port int) *CallbackServer {
	callback := &CallbackServer{
		done:        make(chan bool, 1),
		callBackURL: nil,
	}
	httpServer := &http.Server{Addr: fmt.Sprintf("localhost:%d", port)}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/oidc/callback", callback.handleCallback)
	httpServer.Handler = serverMux
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start callback webserver: %s", err.Error())
		}
	}()
	return callback
}

func (c *CallbackServer) Stop(ctx context.Context) {
	err := c.server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Failed to stop callback webserver: %v", err)
	}
}

func (c *CallbackServer) WaitForCallback() *url.URL {
	<-c.done
	return c.callBackURL
}

func (c *CallbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	c.callBackURL = r.URL
	_, err := io.WriteString(w, `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Authenticated</title>
  <script>setTimeout(function() { window.close() }, 2000);</script>
</head>
<body>
  <h1>Authenticated &#x1F600; This page closes in 2 seconds</h1>
</body>
</html>`)
	if err != nil {
		fmt.Printf("failed to write response: %v", err)
	}
	c.done <- true
}
