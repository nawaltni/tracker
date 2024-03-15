package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/nawaltni/tracker/domain"
)

// RecordPositionRequest represents the request to record a user's location update.
type RecordPositionRequest struct {
	UserID   string `json:"user_id"`
	Location struct {
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	} `json:"location"`
	Timestamp time.Time `json:"timestamp"`
	ClientID  string    `json:"client_id"`
	Metadata  *struct {
		DeviceID   string  `json:"device_id"`
		Brand      string  `json:"brand"`
		Model      string  `json:"model"`
		Os         string  `json:"os"`
		AppVersion string  `json:"app_version"`
		Carrier    string  `json:"carrier"`
		Battery    float64 `json:"battery"`
	}
}

// RecordPosition records a user's location update.
func (a *API) RecordPosition(c *gin.Context) {
	var req RecordPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get the id from the request metadata
	// if the id is not present, generate a new one
	slog.Info("RecordPosition request received", "user_id", req.UserID)
	uid, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate uuid"})
		return
	}

	userPostion := domain.UserPosition{
		UserID:    req.UserID,
		Reference: uid.String(),
		Location: domain.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		CreatedAt: req.Timestamp,
		PhoneMeta: domain.PhoneMeta{
			DeviceID:   req.Metadata.DeviceID,
			Brand:      req.Metadata.Brand,
			Model:      req.Metadata.Model,
			OS:         req.Metadata.Os,
			AppVersion: req.Metadata.AppVersion,
			Carrier:    req.Metadata.Carrier,
			Battery:    int(req.Metadata.Battery),
		},
	}

	err = a.services.UserPositionService.RecordPosition(c, userPostion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record position"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Position recorded successfully"})
}
