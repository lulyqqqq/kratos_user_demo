package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"userDemo/common/util"
	"userDemo/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) QueryUserByName(ctx context.Context, name string) (*biz.Users, error) {
	var user biz.Users
	err := r.data.gormDB.Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	r.log.WithContext(ctx).Info("gormDB: GetUser, user: ", name)
	return &biz.Users{
		Name:     user.Name,
		Password: user.Password,
		Sex:      user.Sex,
		Role:     user.Role,
	}, nil
}

func (r *userRepo) QueryAllUser(ctx context.Context) ([]*biz.Users, error) {
	var userList []*biz.Users
	err := r.data.gormDB.Find(&userList).Error
	if err != nil {
		return nil, err
	}
	return userList, nil
}

func (r *userRepo) DeleteUserByName(ctx context.Context, name string) error {
	var user biz.Users
	err := r.data.gormDB.Where("name = ?", name).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) Save(ctx context.Context, u *biz.Users) error {
	err := r.data.gormDB.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) CheckUserRole(ctx context.Context) (*biz.Users, error) {
	userName := ctx.Value("name")
	var user *biz.Users
	err := r.data.redis.Get(ctx, fmt.Sprintf("%s", userName)).Scan(&user)
	if err != nil {
		return nil, fmt.Errorf("redis中无该用户信息")
	}
	return user, nil
}

func (r *userRepo) SaveRedis(ctx context.Context, user *biz.Users) {
	data, _ := json.Marshal(user)
	r.data.redis.Set(ctx, user.Name, data, time.Hour*2)
}

func (r *userRepo) SaveTokenToRedis(ctx context.Context, name string, token string) {
	r.data.redis.Set(ctx, name, token, time.Hour*2)
}

func (r *userRepo) UpdateUser(ctx context.Context, user *biz.Users) error {
	password := util.HashAndSalt(user.Password)
	err := r.data.gormDB.Model(&biz.Users{}).Where("name = ?", user.Name).Update("password", password).Update("sex", user.Sex).Update("role", user.Role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) FindByID(context.Context, int64) (*biz.Users, error) {
	return nil, nil
}

func (r *userRepo) ListByHello(context.Context, string) ([]*biz.Users, error) {
	return nil, nil
}

func (r *userRepo) ListAll(context.Context) ([]*biz.Users, error) {
	return nil, nil
}
