package grpc

import (
	pb "github.com/nawaltni/api/gen/go/nawalt/auth/v1"
	"github.com/nawaltni/tracker/domain"
)

// MapUserFromProto maps a User proto to a User domain object
func MapUserFromProto(in *pb.User) *domain.User {
	user := domain.User{
		ID:        in.Id,
		BackendID: int(in.BackendId),
		CreatedAt: in.CreatedAt.AsTime(),
		UpdatedAt: in.UpdatedAt.AsTime(),

		Name:     in.Name,
		Email:    in.Email,
		Phone:    in.Phone,
		Address:  in.Address,
		City:     in.City,
		State:    in.State,
		Password: in.Password,
	}

	switch in.Status {
	case pb.User_USER_STATUS_ENABLED:
		user.Status = "ENABLED"
	case pb.User_USER_STATUS_DISABLED:
		user.Status = "DISABLED"
	}

	switch in.Role {
	case pb.User_USER_ROLE_DEVELOPER:
		user.Role = "DEVELOPER"
	case pb.User_USER_ROLE_ADMINISTRATOR:
		user.Role = "ADMINISTRATOR"
	case pb.User_USER_ROLE_SUPERVISOR:
		user.Role = "SUPERVISOR"
	case pb.User_USER_ROLE_DISPLAY:
		user.Role = "DISPLAY"
	case pb.User_USER_ROLE_CLIENT:
		user.Role = "CLIENT"

	}

	return &user
}
