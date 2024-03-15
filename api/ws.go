package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"nhooyr.io/websocket"
)

func (a *API) ws(c *gin.Context) {
	// Upgrade the HTTP connection to a WebSocket

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("Failed to upgrade to websocket", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	defer conn.Close(websocket.StatusInternalError, "Closed unexepetedly")

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()
	ctx = conn.CloseRead(ctx)

	tick := time.NewTicker(time.Second * 5)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Connection closed by client")
			conn.Close(websocket.StatusNormalClosure, "")
			return
		case <-tick.C:
			err := conn.Write(ctx, websocket.MessageText, []byte("Hello"))
			if err != nil {
				slog.Error("Failed to write message", "error", err)
				return
			}
		}
	}
}
