package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/logger"

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
		GenerateExternalToken(ctx context.Context, userId uint) (token string, err error)
		GetUserIdFormToken(ctx context.Context, token string) (uId uint, err error)
		GetUserIdFormRefreshToken(ctx context.Context, token string) (uId uint, err error)
		GetUserIdFromExternalToken(ctx context.Context, token string) (uId uint, err error)
		GenerateTeamToken(ctx context.Context, teamId, userId uint) (token string, err error)
		GetTeamFormToken(ctx context.Context, token string) (res *TeamJwtBody, err error)
	}
	jwtService struct {
		logger logger.AppLogger
		conf   *config.Config
		helper helpers.Helper
	}

	// Type TeamId and UserID
	TeamJwtBody struct {
		TeamId uint `json:"team_id"`
		UserId uint `json:"user_id"`
	}
)

// GetUserIdFormRefreshToken implements JWTService.
func (j *jwtService) GetUserIdFormRefreshToken(ctx context.Context, token string) (uId uint, err error) {
	t, err := j.parserToken(ctx, token, j.conf.JWT.REFRESH_SECRET)
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

// GetTeamFormToken implements JWTService.
func (j *jwtService) GetTeamFormToken(ctx context.Context, token string) (res *TeamJwtBody, err error) {
	t, err := j.parserToken(ctx, token, j.conf.JWT.EMAIL_SECRET)
	if err != nil {
		j.logger.Error(err)
		return nil, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		j.logger.Error("parse claims not ok")
		return nil, errs.ErrGenerateJWTFail
	}
	bodyStr, err := claims.GetSubject()
	if err != nil {
		j.logger.Error(err)
		return nil, err
	}
	if err := json.Unmarshal([]byte(bodyStr), &res); err != nil {
		j.logger.Error(err)
		return nil, err
	}
	return res, nil
}

// GenerateTeamToken implements JWTService.
func (j *jwtService) GenerateTeamToken(ctx context.Context, teamId uint, userId uint) (token string, err error) {
	body := TeamJwtBody{
		TeamId: teamId,
		UserId: userId,
	}
	bodyStr, err := j.helper.ToJson(body)
	if err != nil {
		j.logger.Error(err)
		return "", err
	}
	claims := jwt.MapClaims{
		"sub": bodyStr,
		"exp": time.Now().Add(EmailExpireTime).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(j.conf.JWT.EMAIL_SECRET)
	return t.SignedString(key)
}

// GetUserIdFromExternalToken implements JWTService.
func (j *jwtService) GetUserIdFromExternalToken(ctx context.Context, token string) (uId uint, err error) {
	t, err := j.parserToken(ctx, token, j.conf.JWT.EMAIL_SECRET)
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

// GetUserIdFormToken implements JWTService.
func (j *jwtService) GetUserIdFormToken(ctx context.Context, token string) (uId uint, err error) {
	t, err := j.parserToken(ctx, token, j.conf.JWT.SECRET)
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

func (j *jwtService) parserToken(ctx context.Context, tokenString, secret string) (token *jwt.Token, err error) {
	_ = ctx
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return
}

// CreateVerifyEmailToken implements JWTService.
func (j *jwtService) GenerateExternalToken(ctx context.Context, userId uint) (token string, err error) {
	uId := strconv.Itoa(int(userId))
	claims := jwt.MapClaims{
		"sub": uId,
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

func newJWT(logger logger.AppLogger, conf *config.Config, helper helpers.Helper) JWTService {
	return &jwtService{
		logger: logger,
		conf:   conf,
		helper: helper,
	}
}
