package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaltni/tracker/domain"
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

	wsInstance := NewWebsocketInstance(conn, a.services.UserPositionService)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Minute)
	defer cancel()

	go heartbeat(ctx, conn, 15*time.Second)
	// Setup message channels
	messageChan := make(chan WsRequest)

	// Setup handlers
	handlers := map[string]HandlerFunc{
		"get_positions":    wsInstance.handleGetPositions,
		"stream_positions": wsInstance.handleStreamPositions,
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

type WebsocketInstance struct {
	Conn            *websocket.Conn
	PositionService domain.UserPositionService
}

func NewWebsocketInstance(conn *websocket.Conn, positionService domain.UserPositionService) *WebsocketInstance {
	return &WebsocketInstance{
		Conn:            conn,
		PositionService: positionService,
	}
}

// Example handler function
func (wi WebsocketInstance) handleGetPositions(ctx context.Context, conn *websocket.Conn,
	data interface{},
) error {
	slog.Info("Handling get_positions")

	positions, err := wi.PositionService.GetUsersCurrentPositionByDate(ctx, time.Now())

	// Here you would fetch the latest positions and send them.
	positionsData, err := json.Marshal(map[string]interface{}{
		"type":           "position_update",
		"correlation_id": "your_correlation_id", // You would need to manage correlation IDs appropriately.
		"data":           positions,
	})

	response, err := json.Marshal(positionsData)
	if err != nil {
		return err
	}

	err = conn.Write(ctx, websocket.MessageText, response)
	if err != nil {
		return err
	}

	return nil
}

func (wi WebsocketInstance) handleStreamPositions(ctx context.Context, conn *websocket.Conn, data interface{}) error {
	slog.Info("Handling stream_positions")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	refTime := time.Now()

	for {
		select {
		case <-ctx.Done(): // Make sure to handle context cancellation or deadlines.
			return ctx.Err()
		case <-ticker.C:

			positions, err := wi.PositionService.GetUsersCurrentPositionsSince(ctx, refTime)
			if err != nil {
				slog.Warn("Failed to get positions", "error", err)
				continue
			}

			// Here you would fetch the latest positions and send them.
			positionsData, err := json.Marshal(map[string]interface{}{
				"type":           "position_update",
				"correlation_id": "your_correlation_id", // You would need to manage correlation IDs appropriately.
				"data":           positions,
			})
			if err != nil {
				return err // Handling errors is crucial.
			}

			if err := conn.Write(ctx, websocket.MessageText, positionsData); err != nil {
				return err // Error handling for failed writes.
			}
			refTime = time.Now() // Update the reference time for the next iteration.
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
