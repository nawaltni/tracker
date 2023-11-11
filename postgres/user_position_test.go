package postgres

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nawaltni/tracker/domain"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

func TestUserPositionRepository_GetUserPosition(t *testing.T) {
	// Create some test data

	now := time.Now()

	// With placeID and checkedIn/checkedOut
	id1 := uuid.New().String()
	placeID1 := uuid.New().String()
	placeName1 := "place-name-1"
	err := createSampleUserPosition(db, id1, &placeID1, &placeName1, &now, &now)
	require.Equal(t, nil, err)

	// Without placeID and no checkedIn/checkedOut
	id2 := uuid.New().String()
	err = createSampleUserPosition(db, id2, nil, nil, nil, nil)
	require.Equal(t, nil, err)

	// With placeID and checkin but no checkout
	id3 := uuid.New().String()
	placeID3 := uuid.New().String()
	placeName3 := "place-name-3"
	err = createSampleUserPosition(db, id3, &placeID3, &placeName3, &now, nil)
	require.Equal(t, nil, err)

	t.Cleanup(func() {
		err := db.Exec("DELETE FROM user_positions").Error
		require.Equal(t, nil, err)
	})

	type fields struct {
		client *Client
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.UserPosition
		wantErr bool
	}{
		{
			name: "when the user position is not found we expect and error and nil value",
			fields: fields{
				client: &Client{
					db: db,
				},
			},
			args: args{
				userID: uuid.New().String(), // not existent
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when the user position is found we expect the user position and no error",
			fields: fields{
				client: &Client{
					db: db,
				},
			},
			args: args{
				userID: id1,
			},
			want: &domain.UserPosition{
				UserID:    id1,
				Location:  domain.GeoPoint{Latitude: 40.7128, Longitude: -74.0060},
				CreatedAt: now,
				PlaceID:   &placeID1,
				PlaceName: &placeName1,
				Metadata: domain.PhoneMetadata{
					DeviceID:    "device-uuid-1",
					Model:       "model-1",
					OSVersion:   "os-1",
					Carrier:     "carrier-1",
					CorporateID: "corp-id-1",
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
			got, err := r.GetUserPosition(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPositionRepository.GetUserPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				require.Equal(t, tt.want.UserID, got.UserID)
				require.Equal(t, tt.want.Location.Latitude, got.Location.Latitude)
				require.Equal(t, tt.want.Location.Longitude, got.Location.Longitude)
				require.Equal(t, tt.want.CreatedAt.Round(time.Minute), got.CreatedAt.Round(time.Minute))
				require.Equal(t, tt.want.PlaceID, got.PlaceID)
				require.Equal(t, tt.want.PlaceName, got.PlaceName)
				require.Equal(t, tt.want.Metadata, got.Metadata)
				require.Equal(t, tt.want.Location, got.Location)
			}
		})
	}
}

func createSampleUserPosition(
	db *gorm.DB, userID string, placeID *string, placeName *string, checkedIn, checkedOut *time.Time) error {

	position := &UserPosition{
		UserID:    userID,
		Latitude:  40.7128,
		Longitude: -74.0060,
		PlaceID:   placeID,
		PlaceName: placeName,
		Location:  GeoPoint{Point: orb.Point{40.7128, -74.0060}},

		CreatedAt: time.Now(),
		PhoneMetadata: PhoneMetadata{
			DeviceID:    "device-uuid-1",
			Model:       "model-1",
			OSVersion:   "os-1",
			Carrier:     "carrier-1",
			CorporateID: "corp-id-1",
		},
	}

	position.CheckedIn = checkedIn
	position.CheckedOut = checkedOut

	return db.Create(position).Error
}

func TestUserPositionRepository_GetUsersPositionByCoordinates(t *testing.T) {

	// Create some test data

	now := time.Now()

	// With placeID and checkedIn/checkedOut
	id1 := uuid.New().String()
	placeID1 := uuid.New().String()
	placeName1 := "place-name-1"
	err := createSampleUserPosition(db, id1, &placeID1, &placeName1, &now, &now)
	require.Equal(t, nil, err)

	// Without placeID and no checkedIn/checkedOut
	id2 := uuid.New().String()
	err = createSampleUserPosition(db, id2, nil, nil, nil, nil)
	require.Equal(t, nil, err)

	// With placeID and checkin but no checkout
	id3 := uuid.New().String()
	placeID3 := uuid.New().String()
	placeName3 := "place-name-3"
	err = createSampleUserPosition(db, id3, &placeID3, &placeName3, &now, nil)
	require.Equal(t, nil, err)

	t.Cleanup(func() {
		err := db.Exec("DELETE FROM user_positions").Error
		require.Equal(t, nil, err)
	})

	type fields struct {
		client *Client
	}
	type args struct {
		lat      float64
		lon      float64
		distance int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "when no users are found we expect an empty list and no error",
			fields: fields{
				client: &Client{
					db: db,
				},
			},
			args: args{
				lat:      40.7128,
				lon:      -80.0060,
				distance: 1000,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "when users are found we expect a list of users and no error",
			fields: fields{
				client: &Client{
					db: db,
				},
			},
			args: args{
				lat:      40.7128,
				lon:      -74.0061,
				distance: 1000,
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserPositionRepository{
				client: tt.fields.client,
			}
			got, err := r.GetUsersPositionByCoordinates(tt.args.lat, tt.args.lon, tt.args.distance)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPositionRepository.GetUsersPositionByCoordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, len(got))
		})
	}
}
