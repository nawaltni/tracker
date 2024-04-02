package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/nawaltni/tracker/domain"
)

// RecordPositionRequest represents the request to record a user's location update.
type RecordPositionRequest struct {
	UserID   string `json:"userId"`
	Location struct {
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	} `json:"location"`
	Timestamp time.Time `json:"timestamp"`
	ClientID  string    `json:"clientId"`
	Metadata  *struct {
		DeviceID   string  `json:"deviceId"`
		Brand      string  `json:"brand"`
		Model      string  `json:"model"`
		Os         string  `json:"os"`
		AppVersion string  `json:"appVersion"`
		Carrier    string  `json:"carrier"`
		Battery    float64 `json:"battery"`
	}
	RefID string `json:"refId"`
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
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record position"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": uid, "refId": req.RefID})
}

// RecordPositionBatch records a batch of user's location updates.
func (a *API) RecordPositionBatch(c *gin.Context) {
	var req []RecordPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out := map[string]string{}

	for _, r := range req {
		// get the id from the request metadata
		// if the id is not present, generate a new one
		uid, err := uuid.NewV7()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate uuid"})
			return
		}

		userPostion := domain.UserPosition{
			UserID:    r.UserID,
			Reference: uid.String(),
			Location: domain.GeoPoint{
				Latitude:  r.Location.Latitude,
				Longitude: r.Location.Longitude,
			},
			CreatedAt: r.Timestamp,
			PhoneMeta: domain.PhoneMeta{
				DeviceID:   r.Metadata.DeviceID,
				Brand:      r.Metadata.Brand,
				Model:      r.Metadata.Model,
				OS:         r.Metadata.Os,
				AppVersion: r.Metadata.AppVersion,
				Carrier:    r.Metadata.Carrier,
				Battery:    int(r.Metadata.Battery),
			},
		}

		err = a.services.UserPositionService.RecordPosition(c, userPostion)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record position"})
			return
		}

		out[r.RefID] = uid.String()

	}
	c.JSON(http.StatusOK, mapToUserPositionBatchResponse(out))
}

type UserPositionBatchResponse struct {
	RefID string `json:"refId"`
	ID    string `json:"id"`
}

func mapToUserPositionBatchResponse(m map[string]string) []UserPositionBatchResponse {
	var s []UserPositionBatchResponse
	for k, v := range m {
		s = append(s, UserPositionBatchResponse{RefID: k, ID: v})
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].RefID < s[j].RefID
	})
	return s
}
