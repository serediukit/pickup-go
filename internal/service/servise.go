package service

import (
	"context"
	"pickup-srv/internal/repository"
	"pickup-srv/proto"
)

type PickupService struct {
	proto.UnimplementedPickupServiceServer

	userRepo *repository.UserRepository
}

func NewPickupService(userRepo *repository.UserRepository) *PickupService {
	return &PickupService{
		userRepo: userRepo,
	}
}

func (s *PickupService) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	users, err := s.userRepo.GetUsers(req.UserSearchParams, req.Limit)
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*proto.User, len(users))
	for i, user := range users {
		protoUsers[i] = &proto.User{
			Id:    user.Id,
			Name:  user.Name,
			Age:   user.Age,
			City:  user.City,
			RegDt: user.RegDt,
		}
	}

	return &proto.GetUsersResponse{
		Users: protoUsers,
		Total: int32(len(protoUsers)),
	}, nil
}
