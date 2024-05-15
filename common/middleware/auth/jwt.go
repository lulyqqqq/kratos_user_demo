package auth

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	jwt2 "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	// 可根据需要自行添加字段
	Name                 string `json:"name"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("user_demo")

// TokenExpireDuration 超时时间
const TokenExpireDuration = time.Hour * 2

// CreateToken 生成JWT
func CreateToken(name string) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		name, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "mfn", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (string, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return "", err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims.Name, nil
	}
	return "", err
}

func Middleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				token := tr.RequestHeader().Get("Authorization")
				if token == "" {
					return nil, jwt2.ErrMissingJwtToken
				}
				name, err := ParseToken(token)
				if name == "" || err != nil {
					return nil, jwt2.ErrTokenInvalid
				}
				// 将userId放到context中
				ctx = context.WithValue(ctx, "name", name)
			}
			return handler(ctx, req)
		}
	}
}
