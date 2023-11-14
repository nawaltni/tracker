package services

import (
	"testing"
	"time"

	"github.com/nawaltni/tracker/domain"
	"github.com/stretchr/testify/require"
)

func TestUserPositionService_CalculateUserPosition(t *testing.T) {
	now := time.Now()

	type fields struct {
		repo             domain.UserPositionRepository
		placesClientGRPC domain.PlacesClientGRPC
	}
	type args struct {
		userPosition  *domain.UserPosition
		knownPosition domain.UserPosition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.UserPosition
	}{
		{
			name: "user is entering a place",
			fields: fields{
				repo:             nil,
				placesClientGRPC: nil,
			},
			args: args{
				knownPosition: domain.UserPosition{
					PlaceID: nil,
				},
				userPosition: &domain.UserPosition{
					CreatedAt: now,
					PlaceID:   stringPointer("c390300b-300e-49ce-9a8f-157834883d4d"),
				},
			},
			want: &domain.UserPosition{
				CreatedAt: now,
				PlaceID:   stringPointer("c390300b-300e-49ce-9a8f-157834883d4d"),
				CheckedIn: timePointer(now),
			},
		},
		{
			name: "user is leaving a place",
			fields: fields{
				repo:             nil,
				placesClientGRPC: nil,
			},
			args: args{
				knownPosition: domain.UserPosition{
					PlaceID: stringPointer("c390300b-300e-49ce-9a8f-157834883d4d"),
				},
				userPosition: &domain.UserPosition{
					CreatedAt: now,
					PlaceID:   nil,
				},
			},
			want: &domain.UserPosition{
				CreatedAt:  now,
				PlaceID:    nil,
				CheckedOut: timePointer(now),
			},
		},
		{
			name: "user is not entering or leaving a place",
			fields: fields{
				repo:             nil,
				placesClientGRPC: nil,
			},
			args: args{
				knownPosition: domain.UserPosition{
					PlaceID: nil,
				},
				userPosition: &domain.UserPosition{
					CreatedAt: now,
					PlaceID:   nil,
				},
			},
			want: &domain.UserPosition{
				CreatedAt: now,
				PlaceID:   nil,
			},
		},
		{
			name: "user goes from street to street",
			fields: fields{
				repo:             nil,
				placesClientGRPC: nil,
			},
			args: args{
				knownPosition: domain.UserPosition{
					PlaceID: stringPointer("0000-0000-0000-0000"),
				},
				userPosition: &domain.UserPosition{
					CreatedAt: now,
					PlaceID:   stringPointer("0000-0000-0000-0000"),
				},
			},
			want: &domain.UserPosition{
				CreatedAt: now,
				PlaceID:   stringPointer("0000-0000-0000-0000"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserPositionService{
				repo:             tt.fields.repo,
				placesClientGRPC: tt.fields.placesClientGRPC,
			}
			s.CalculateUserPosition(tt.args.userPosition, tt.args.knownPosition)

			require.Equal(t, tt.want, tt.args.userPosition)
		})
	}
}
