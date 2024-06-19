package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping is a simple checker for User's connection tk to the server.
// It returns the caller's ip.
func (a *API) Ping(c *gin.Context) {
	// // print all the headers
	// for k, v := range c.Request.Header {
	// 	a.log.Info().Msgf("Header: %s: %s", k, v)
	// }

	c.JSON(http.StatusOK, gin.H{"ip": c.ClientIP()})
}
