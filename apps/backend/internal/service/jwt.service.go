package service

import (
	"context"
	"fmt"
	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/logger"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ExpireTime             = time.Hour * 24
	RefreshTokenExpireTime = time.Hour * 24 * 7
	EmailExpireTime        = time.Hour * 24 * 30
)

type (
	JWTService interface {
		GenerateToken(ctx context.Context, userId uint) (token string, err error)
		GenerateRefreshToken(ctx context.Context, userId uint) (token string, err error)
		CreateVerifyEmailToken(ctx context.Context, email string) (token string, err error)
		GetEmailFormToken(ctx context.Context, token string) (email string, err error)
		GetUserIdFormToken(ctx context.Context, token string) (uId uint, err error)
		ParserToken(ctx context.Context, t, secret string) (token *jwt.Token, err error)
	}
	jwtService struct {
		logger logger.AppLogger
		conf   *config.Config
	}
)

// GetUserIdFormToken implements JWTService.
func (j *jwtService) GetUserIdFormToken(ctx context.Context, token string) (uId uint, err error) {
	fmt.Println("token", token)
	t, err := j.ParserToken(ctx, token, j.conf.JWT.SECRET)
	if err != nil {
		j.logger.Error(err)
		return 0, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		j.logger.Error("parse claims not ok")
		return 0, errs.ErrGenerateJWTFail
	}
	uIdStr, err := claims.GetSubject()
	if err != nil {
		j.logger.Error(err)
		return 0, err
	}
	u, err := strconv.Atoi(uIdStr)
	if err != nil {
		j.logger.Error(err)
		return 0, err
	}
	return uint(u), nil
}

func (j *jwtService) ParserToken(ctx context.Context, tokenString, secret string) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return
}

func (j *jwtService) GetEmailFormToken(ctx context.Context, token string) (email string, err error) {
	t, err := j.ParserToken(ctx, token, j.conf.JWT.EMAIL_SECRET)
	if err != nil {
		j.logger.Error(err)
		return "", err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		j.logger.Error("parse claims not ok")
		return "", errs.ErrGenerateJWTFail
	}
	return claims.GetSubject()
}

// CreateVerifyEmailToken implements JWTService.
func (j *jwtService) CreateVerifyEmailToken(ctx context.Context, email string) (token string, err error) {
	claims := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(EmailExpireTime).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(j.conf.JWT.EMAIL_SECRET)
	return t.SignedString(key)
}

func (j *jwtService) GenerateToken(ctx context.Context, userId uint) (token string, err error) {
	uId := strconv.Itoa(int(userId))
	claims := jwt.MapClaims{
		"sub": uId,
		"exp": time.Now().Add(ExpireTime).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(j.conf.JWT.SECRET)
	return t.SignedString(key)
}

func (j *jwtService) GenerateRefreshToken(ctx context.Context, userId uint) (token string, err error) {
	uId := strconv.Itoa(int(userId))
	claims := jwt.MapClaims{
		"sub": uId,
		"exp": time.Now().Add(RefreshTokenExpireTime).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(j.conf.JWT.REFRESH_SECRET)
	return t.SignedString(key)
}

func newJWT(logger logger.AppLogger, conf *config.Config) JWTService {
	return &jwtService{
		logger: logger,
		conf:   conf,
	}
}
