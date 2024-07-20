package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type callbackServer struct {
	done        chan bool
	callBackURL *url.URL
	server      http.Server
}

func startCallbackServer(port int) *callbackServer {
	callback := &callbackServer{
		done:        make(chan bool, 1),
		callBackURL: nil,
	}
	httpServer := &http.Server{Addr: fmt.Sprintf("localhost:%d", port)}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/oidc/callback", callback.handleCallback)
	httpServer.Handler = serverMux
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("failed to start callback webserver: %s", err.Error())
		}
	}()
	return callback
}

func (c *callbackServer) stop(ctx context.Context) {
	err := c.server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("failed to stop callback webserver: %v", err)
	}
}

func (c *callbackServer) waitForCallback() *url.URL {
	<-c.done
	return c.callBackURL
}

func (c *callbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, `
	<!DOCTYPE html>
	<html>
	<head>
	<script>
		setTimeout(function() {
			window.close()
		}, 2000);
	</script>
	</head>
	<body><p>Authenticated! This page closes in 2 seconds</p></body>
	</html>
	`)
	if err != nil {
		fmt.Printf("failed to write response: %v", err)
	}
	c.callBackURL = r.URL
	c.done <- true
}
