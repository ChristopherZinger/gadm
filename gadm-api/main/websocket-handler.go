package main

import (
	"errors"
	"gadm-api/logger"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func getWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("websocket_request_received remote=%s origin=%s", r.RemoteAddr, r.Header.Get("Origin"))
	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// TODO: tighten for production. The webapp dev server runs on a
		// different origin than the API, so we accept any origin here.
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		logger.Error("failed_to_accept_websocket %v", err)
		return
	}
	defer wsConn.CloseNow()

	logger.Info("websocket_connected")

	ctx := r.Context()

	if err := wsjson.Write(ctx, wsConn, map[string]any{
		"type": "hello",
		"msg":  "welcome",
	}); err != nil {
		logger.Error("failed_to_write_welcome_message %v", err)
		return
	}

	for {
		var v any
		if err := wsjson.Read(ctx, wsConn, &v); err != nil {
			status := websocket.CloseStatus(err)
			if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
				logger.Info("websocket_closed_by_client status=%d", status)
				return
			}
			if errors.Is(err, ctx.Err()) {
				logger.Info("websocket_context_done")
				return
			}
			logger.Error("failed_to_read_websocket_message %v", err)
			return
		}

		log.Printf("received: %v", v)

		err := wsjson.Write(ctx, wsConn, map[string]any{
			"type": "echo",
			"data": v,
		})
		if err != nil {
			logger.Error("failed_to_write_echo_message %v", err)
			return
		}
	}
}
