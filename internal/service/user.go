package service

import (
	"context"
	"fmt"

	v1 "userDemo/api/user/v1"
	"userDemo/internal/biz"

	"github.com/jinzhu/copier"
)

// GreeterService is a greeter service.
type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUsecase
}

// NewGreeterService new a greeter service.
func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
// func (s *UserService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
// 	g, err := s.uc.CreateUser(ctx, &biz.User{Hello: in.Name})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
// }

// 注册
func (s *UserService) Register(ctx context.Context, in *v1.UserRegister) (*v1.UserRegisterReply, error) {
	err := s.uc.Register(ctx, &biz.Users{
		Name:     in.Name,
		Password: in.Password,
		Sex:      in.Sex,
		Role:     in.Role,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UserRegisterReply{
		Code:    200,
		Message: "注册成功",
	}, nil
}

// 登录
func (s *UserService) Login(ctx context.Context, req *v1.UserLoginRequest) (*v1.UserLoginReply, error) {
	token, err := s.uc.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.UserLoginReply{
		Code:    200,
		Message: "登录成功",
		Token:   token,
	}, nil
}

// 查询全部用户
func (s *UserService) QuertyAllUser(ctx context.Context, in *v1.QueryAllUserRequest) (*v1.QueryAllUserReply, error) {
	userList, err := s.uc.QuertyAllUser(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.QueryAllUserReply{
		User: userList,
	}, nil
}

// 根据名称查询单个用户
func (s *UserService) QuertyUserByName(ctx context.Context, in *v1.QueryUserRequest) (*v1.QueryUserReply, error) {
	user, err := s.uc.QuertyUserByName(ctx, in.Name)
	fmt.Println(in)
	if err != nil {
		return nil, err
	}
	var userInfo v1.UserInfo
	copier.Copy(&userInfo, &user)
	return &v1.QueryUserReply{
		User: &userInfo,
	}, nil
}

// 删除用户
func (s *UserService) DeleteUserByName(ctx context.Context, req *v1.DeleteUserByNameRequest) (*v1.DeleteUserByNameReply, error) {
	err := s.uc.DeleteUserByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteUserByNameReply{
		Code:    200,
		Message: "删除成功",
	}, nil
}

// 修改用户信息
