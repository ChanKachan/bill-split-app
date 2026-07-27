package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func NewUpgrader(
	readBufferSize, writeBufferSize int,
) *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin:     requestFunc,
		Error:           handleUpgradeError,
	}
}

func requestFunc(r *http.Request) bool {
	//todo Здесь позже нужно проверять разрешённые Origin.
	return true
}

func handleUpgradeError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	reason error,
) {
	log.Printf(
		"websocket upgrade failed: method=%s path=%s status=%d error=%v",
		r.Method,
		r.URL.Path,
		status,
		reason,
	)

	response := struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{
		Error:   "websocket_upgrade_failed",
		Message: upgradeErrorMessage(status),
		Code:    status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("encode websocket upgrade error: %v", err)
	}
}

func upgradeErrorMessage(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "invalid websocket handshake"
	case http.StatusForbidden:
		return "websocket origin is not allowed"
	case http.StatusMethodNotAllowed:
		return "websocket connection requires GET request"
	case http.StatusInternalServerError:
		return "failed to establish websocket connection"
	default:
		return "websocket connection failed"
	}
}
