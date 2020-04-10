package logic

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"

	"micro-message-system/common/baseerror"
	"micro-message-system/common/middleware"
	"micro-message-system/userserver/models"
)

type (
	UserLogic struct {
		userModel *models.MembersModel
	}
	LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	LoginResponse struct {
		Token       string `json:"token"`
		AccessToken string `json:"accessToken"`
		ExpireAt    int64  `json:"expireAt"`
		TimeStamp   int64  `json:"timeStamp"`
	}

	RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	RegisterResponse struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}
)

var (
	NotFoundUserErr       = baseerror.NewBaseError("用户不存在")
	UserNameOrPasswordErr = baseerror.NewBaseError("用户不存在或者密码错误")
	AccessTokenErr        = baseerror.NewBaseError("生成签名错误")
	CreateMemberErr       = baseerror.NewBaseError("注册失败")
	ExistsUserErr         = baseerror.NewBaseError("用户已存在，无法注册")
)

func NewUserLogic(userModel *models.MembersModel) *UserLogic {

	return &UserLogic{userModel: userModel}
}

// 登录
func (l *UserLogic) Login(r *LoginRequest) (*LoginResponse, error) {
	user, err := l.userModel.FindByUserName(r.Username)
	if err != nil {
		return nil, NotFoundUserErr
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(r.Password))) {
		return nil, UserNameOrPasswordErr
	}

	expired := time.Now().Add(7 * 24 * time.Hour).Unix()
	accessToken, err := l.createAccessToken(expired)
	if err != nil {
		return nil, AccessTokenErr
	}
	return &LoginResponse{
		Token:       user.Token,
		AccessToken: accessToken,
		ExpireAt:    expired,
		TimeStamp:   time.Now().Unix(),
	}, nil
}

// 注册
func (l *UserLogic) Register(r *RegisterRequest) (*RegisterResponse, error) {
	member := &models.Members{
		Token:    uuid.NewV4().String(),
		Username: r.Username,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(r.Password))),
	}
	mem, err := l.userModel.FindByUserName(member.Username)
	if mem != nil || err != nil {
		return nil, ExistsUserErr
	}
	if _, err := l.userModel.InsertMember(member); err != nil {
		return nil, CreateMemberErr
	}
	return &RegisterResponse{Token: member.Token, Username: r.Username}, nil
}

// 生成jwt token
func (l *UserLogic) createAccessToken(expired int64) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: expired,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(middleware.UserSignedKey))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
