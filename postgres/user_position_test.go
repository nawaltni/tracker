package postgres

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nawaltni/tracker/domain"
)

func TestUserPositionRepository_Insert(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		userPosition *domain.UserPosition
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		teardown func()
	}{
		{
			name: "when insert is successful",
			fields: fields{
				client: &Client{
					db: db,
				},
			},
			args: args{
				userPosition: &domain.UserPosition{
					UserID:    uuid.New().String(),
					Location:  domain.GeoPoint{Latitude: 40.7128, Longitude: -74.0060},
					CreatedAt: time.Now(),
					Metadata: domain.PhoneMetadata{
						DeviceID:    "device-uuid-1",
						Model:       "model-1",
						OSVersion:   "os-1",
						Carrier:     "carrier-1",
						CorporateID: "corp-id-1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "when no coordinates are provided",
			fields: fields{
				client: &Client{
					db: db,
				},
			},
			args: args{
				userPosition: &domain.UserPosition{
					UserID:    uuid.New().String(),
					CreatedAt: time.Now(),
					Metadata: domain.PhoneMetadata{
						DeviceID:    "device-uuid-1",
						Model:       "model-1",
						OSVersion:   "os-1",
						Carrier:     "carrier-1",
						CorporateID: "corp-id-1",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserPositionRepository{
				client: tt.fields.client,
			}
			if err := r.Insert(tt.args.userPosition); (err != nil) != tt.wantErr {
				t.Errorf("UserPositionRepository.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
