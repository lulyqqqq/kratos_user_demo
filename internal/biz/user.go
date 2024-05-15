package biz

import (
	"context"
	"fmt"

	v1 "userDemo/api/user/v1"
	"userDemo/common/middleware/auth"
	"userDemo/common/util"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.

type Users struct {
	Name     string
	Password string
	Sex      string
	Role     int32
}

// GreeterRepo is a Greater repo.
type UserRepo interface {
	Save(context.Context, *Users) error
	Update(context.Context, *Users) (*Users, error)
	FindByID(context.Context, int64) (*Users, error)
	ListByHello(context.Context, string) ([]*Users, error)
	ListAll(context.Context) ([]*Users, error)
	QueryUserByName(context.Context, string) (*Users, error)
	QueryAllUser(context.Context) ([]*Users, error)
	SaveRedis(context.Context, *Users)
	SaveTokenToRedis(context context.Context, name, token string)
	DeleteUserByName(context context.Context, name string) error
	CheckUserRole(context context.Context) (*Users, error)
}

// GreeterUsecase is a Greeter usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}
func newUser(user *Users) *Users {
	// 密码混淆
	hashPassword := util.HashAndSalt(user.Password)
	return &Users{
		Name:     user.Name,
		Password: hashPassword,
		Sex:      user.Sex,
		Role:     user.Role,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, u *Users) error {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", u)
	// 查询用户名称是否已存在
	user, _ := uc.QuertyUserByName(ctx, u.Name)
	if user != nil {
		return fmt.Errorf("用户名称已存在, 请重新输入")
	}
	newUser := newUser(u)
	err := uc.repo.Save(ctx, newUser)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, req *v1.UserLoginRequest) (string, error) {

	// 查询是否有这个人
	user, _ := uc.QuertyUserByName(ctx, req.Name)
	if user == nil {
		return "", fmt.Errorf("用户不存在")
	}
	// 验证密码
	if !util.ComparePasswords(user.Password, req.Password) {
		return "", fmt.Errorf("密码错误")
	}
	// 生成token
	token, err := auth.CreateToken(user.Name)
	if err != nil {
		return "", nil
	}
	// 将用户信息存入redis
	uc.repo.SaveRedis(ctx, user)

	// token信息存入redis 持久化登录
	// uc.repo.SaveTokenToRedis(ctx,user.Name,token)
	return token, nil
}

func (uc *UserUsecase) QuertyUserByName(ctx context.Context, name string) (*Users, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", name)
	return uc.repo.QueryUserByName(ctx, name)
}

func (uc *UserUsecase) QuertyAllUser(ctx context.Context) ([]*v1.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("QueryAllUser")
	userList, err := uc.repo.QueryAllUser(ctx)
	userInfoList := make([]*v1.UserInfo, len(userList))
	for i := 0; i < len(userList); i++ {
		user := &v1.UserInfo{
			Name:     userList[i].Name,
			Password: userList[i].Password,
			Sex:      userList[i].Sex,
			Role:     userList[i].Role,
		}
		userInfoList[i] = user
	}
	if err != nil {
		return nil, err
	}
	return userInfoList, err
}

func (uc *UserUsecase) DeleteUserByName(ctx context.Context, name string) error {
	uc.log.WithContext(ctx).Infof("DeleteUserByName", name)
	// 验证权限是否可以删除
	user, err := uc.repo.CheckUserRole(ctx)
	if err != nil || user == nil {
		return err
	}
	if user.Role != 1 {
		return fmt.Errorf("没有权限删除")
	}
	return uc.repo.DeleteUserByName(ctx, name)
}
