package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"nhooyr.io/websocket"
)

// WsRequest represents the request made to the websocket server
type WsRequest struct {
	Type          string          `json:"type"`
	CorrelationID string          `json:"correlation_id"`
	Data          json.RawMessage `json:"data"`
}

// Define a handler type
type HandlerFunc func(ctx context.Context, conn *websocket.Conn, data interface{}) error

func (a *API) ws(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("Failed to upgrade to websocket", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "The connection closed unexpectedly")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Minute)
	defer cancel()

	go heartbeat(ctx, conn, 15*time.Second)
	// Setup message channels
	messageChan := make(chan WsRequest)

	// Setup handlers
	handlers := map[string]HandlerFunc{
		"get_positions":    handleGetPositions,
		"stream_positions": handleStreamPositions,
	}

	// Dispatcher goroutine
	go func() {
		for msg := range messageChan {
			if handler, exists := handlers[msg.Type]; exists {
				err := handler(ctx, conn, msg.Data)
				if err != nil {
					slog.Error("Handler error", "error", err)
					continue
				}

			} else {
				slog.Error("No handler for message type", "type", msg.Type)
			}
		}
	}()

	// Main loop to read messages
	for {
		slog.Info("Reading message")
		_, msg, err := conn.Read(ctx)
		if err != nil {
			slog.Error("Failed to read message", "error", err)
			close(messageChan)
			return
		}

		var incomingMsg WsRequest
		if err := json.Unmarshal(msg, &incomingMsg); err != nil {
			slog.Error("Failed to unmarshal message", "error", err)
			continue
		}

		messageChan <- incomingMsg
	}
}

// Example handler function
func handleGetPositions(ctx context.Context, conn *websocket.Conn,
	data interface{},
) error {
	slog.Info("Handling get_positions")
	out := map[string]interface{}{
		"message": "Initial positions",
	}

	response, err := json.Marshal(out)
	if err != nil {
		return err
	}

	err = conn.Write(ctx, websocket.MessageText, response)
	if err != nil {
		return err
	}

	return nil
}

func handleStreamPositions(ctx context.Context, conn *websocket.Conn, data interface{}) error {
	slog.Info("Handling stream_positions")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // Make sure to handle context cancellation or deadlines.
			return ctx.Err()
		case <-ticker.C:
			// Here you would fetch the latest positions and send them.
			positionsData, err := json.Marshal(map[string]interface{}{
				"type":           "position_update",
				"correlation_id": "your_correlation_id", // You would need to manage correlation IDs appropriately.
				"data": map[string]interface{}{
					"positions": []string{"Position1", "Position2"}, // Example positions.
				},
			})
			if err != nil {
				return err // Handling errors is crucial.
			}

			if err := conn.Write(ctx, websocket.MessageText, positionsData); err != nil {
				return err // Error handling for failed writes.
			}
		}
	}
}

func heartbeat(ctx context.Context, c *websocket.Conn, d time.Duration) {
	slog.Info("Starting heartbeat")
	t := time.NewTimer(d)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}

		slog.Info("Sending ping")
		err := c.Ping(ctx)
		if err != nil {
			slog.Error("Failed to send ping", "error", err)
			return
		}

		t.Reset(time.Second * 15)
	}
}
