package server

import (
	"context"
	"errors"

	pb "github.com/golang-grpc-proxy/generated"
	"github.com/golang-grpc-proxy/src/entities"
	"github.com/golang-grpc-proxy/src/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type apiServer struct {
	pb.UnimplementedApiServer
	api services.Api
}

func RegisterGrpcServer(grpcServer *grpc.Server, api services.Api) {
	pb.RegisterApiServer(grpcServer, &apiServer{api: api})
}

func (s *apiServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := entities.User{Nickname: in.Nickname, Email: in.Email, Phone: in.Phone, Password: in.Password, FirstName: in.FirstName, Surname: in.Surname, LastName: in.LastName, Work: in.Work, Study: in.Study, Telegram: in.Telegram}

	id, err := s.api.Register(user)
	if err == nil {
		return &pb.RegisterResponse{Result: true, ErrorMessage: "", Id: int64(id)}, nil
	}

	if errors.Is(err, services.ErrUserExist) {
		return &pb.RegisterResponse{Result: false, ErrorMessage: err.Error(), Id: 0}, status.Error(codes.AlreadyExists, err.Error())
	}

	if errors.Is(err, services.ErrInvalidBody) {
		return &pb.RegisterResponse{Result: false, ErrorMessage: err.Error(), Id: 0}, status.Error(codes.InvalidArgument, err.Error())
	}

	internalErr := services.ErrInternal.Error()
	return &pb.RegisterResponse{Result: false, ErrorMessage: internalErr, Id: 0}, status.Error(codes.Internal, internalErr)
}

func (s *apiServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	err := s.api.AuthentificateUser(in.Email, in.Password)
	if err == nil {
		return &pb.GetUserResponse{Result: true, ErrorMessage: ""}, nil
	}

	if errors.Is(err, services.ErrUserNotExist) {
		return &pb.GetUserResponse{Result: false, ErrorMessage: err.Error()}, status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, services.ErrInvalidCredentials) {
		return &pb.GetUserResponse{Result: false, ErrorMessage: err.Error()}, status.Error(codes.InvalidArgument, err.Error())
	}

	internalErr := services.ErrInternal.Error()
	return &pb.GetUserResponse{Result: false, ErrorMessage: internalErr}, status.Error(codes.InvalidArgument, internalErr)
}

func (s *apiServer) GetToken(ctx context.Context, in *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	token, err := s.api.GetToken(in.Email, in.Password)

	if err == nil {
		return &pb.GetTokenResponse{Result: true, Token: token, ErrorMessage: ""}, nil
	}

	if errors.Is(err, services.ErrUserNotExist) {
		return &pb.GetTokenResponse{Result: false, Token: "", ErrorMessage: err.Error()}, status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, services.ErrInvalidCredentials) {
		return &pb.GetTokenResponse{Result: false, Token: "", ErrorMessage: err.Error()}, status.Error(codes.InvalidArgument, err.Error())
	}

	internalErr := services.ErrInternal.Error()
	return &pb.GetTokenResponse{Result: false, Token: "", ErrorMessage: internalErr}, status.Error(codes.Internal, internalErr)
}

func (s *apiServer) GetContent(ctx context.Context, in *pb.GetContentRequest) (*pb.GetContentResponse, error) {
	return &pb.GetContentResponse{Result: true}, nil
}