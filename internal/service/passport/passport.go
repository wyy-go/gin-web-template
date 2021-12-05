package passport

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go/test"
	"github.com/wyy-go/go-web-template/internal/common/context"
	"github.com/wyy-go/go-web-template/internal/common/jwt"
	"github.com/wyy-go/go-web-template/internal/common/util"
	"github.com/wyy-go/go-web-template/pkg/idgen"
	"gorm.io/gorm"

	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	"golang.org/x/crypto/bcrypt"

	"github.com/wyy-go/go-web-template/internal/common/constant"
	"github.com/wyy-go/go-web-template/internal/common/errors"
	"github.com/wyy-go/go-web-template/internal/dao"
	"github.com/wyy-go/go-web-template/internal/model"
	log "github.com/wyy-go/go-web-template/pkg/logger"
	proto "github.com/wyy-go/go-web-template/proto/passport"
)

const (
	TokenExpiresIn        = 15 * 24 * 3600
	RefreshTokenExpiresIn = 30 * 24 * 3600
)

type IpLimit struct {
	IP            string
	LastErrorTime int64
	ErrorTimes    int64
}

type SmsHistory struct {
	CreateTime int64  `json:"create_time"` // 创建时间
	LastTime   int64  `json:"last_time"`   // 最近一次获取短信验证码时间
	SendCount  int    `json:"send_count"`  // 发送次数
	Code       string `json:"code"`        // 验证码
	ErrorCount int    `json:"error_count"` // 校验次数，同一个验证码超过3次作废
}

type Service struct {
	db             *gorm.DB
	authPrivateKey *rsa.PrivateKey
	authPublicKey  *rsa.PublicKey
}

var (
	service *Service
	once    sync.Once
)

func GetService() *Service {
	once.Do(func() {
		authPrivateKey := test.LoadRSAPrivateKeyFromDisk("./conf/auth_key")
		authPublicKey := test.LoadRSAPublicKeyFromDisk("./conf/auth_key.pub")
		service = &Service{db: dao.GetDB(), authPrivateKey: authPrivateKey, authPublicKey: authPublicKey}
	})
	return service
}

func (s *Service) AuthToken(token string) (acc *jwt.Account, err error) {
	if acc, err = jwt.Decode(s.authPublicKey, token); err != nil {
		return
	}
	var t string
	tokenKey := fmt.Sprintf(constant.RedisKeyToken, 1, acc.Platform, acc.Uid)
	client := dao.GetRedisClient()
	t, err = client.Get(tokenKey).Result()
	if err == redis.Nil {
		err = errors.ErrTokenRevoked
		return
	} else if err != nil {
		log.Error(err)
		return nil, err
	}

	if token != t {
		return nil, errors.ErrInvalidToken
	}

	return
}

func (s *Service) RefreshToken(ctx *context.Context, req *proto.RefreshTokenRequest) (token *proto.TokenInfo, err error) {
	var acc *jwt.Account
	if acc, err = jwt.Decode(s.authPublicKey, req.RefreshToken); err != nil {
		log.Error(err)
		return
	}

	refreshTokenKey := fmt.Sprintf(constant.RedisKeyRefreshToken, 1, acc.Platform, acc.Uid)

	var t string
	client := dao.GetRedisClient()
	t, err = client.Get(refreshTokenKey).Result()
	if err == redis.Nil {
		err = errors.ErrTokenRevoked
		log.Error(err)
		return
	} else if err != nil {
		log.Error(err)
		return
	}

	if t != req.RefreshToken {
		err = errors.ErrInvalidToken
		log.Error(err)
		return
	}

	token, err = s.updateToken(ctx, acc)
	return
}

func (s *Service) Login(ctx *context.Context, req *proto.LoginRequest) (rsp *proto.TokenInfo, err error) {
	v := model.User{}
	if req.Type == 1 {
		v.Mobile = req.Account
		v.MobileVerified = 1
	} else if req.Type == 2 {
		v.LuoboId = req.Account
	} else if req.Type == 3 {
		v.Email = req.Account
		v.EmailVerified = 1
	}
	if err = s.db.Where(&v).First(&v).Error; err != nil {
		if util.IsNotFound(err) {
			err = errors.ErrUserNotExists
		}
		return
	}

	if v.Passwd == "" {
		err = errors.New("未设置密码，不能用密码登录")
		return
	}

	var ipLimits []*IpLimit
	var curIpLimit *IpLimit
	if v.LoginIpLimit != "" {
		_ = json.Unmarshal([]byte(v.LoginIpLimit), &ipLimits)
	}

	// 清理过期记录
	i := 0
	for _, item := range ipLimits {
		if time.Now().Unix() < item.LastErrorTime+constant.RetryCD {
			// 未过冷却时间
			ipLimits[i] = item
			if item.IP == ctx.Context().ClientIP() {
				curIpLimit = item
			}
			i++
		}
	}
	ipLimits = ipLimits[:i]

	if curIpLimit != nil && curIpLimit.ErrorTimes >= 6 {
		msg := fmt.Sprintf("密码错误次数过多，请%s再试", util.HumanTime(curIpLimit.LastErrorTime+constant.RetryCD))
		err = errors.NewError(errors.ErrErrorTimesLimit.Code, msg)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(v.Passwd), []byte(req.Passwd)); err != nil {
		err = errors.ErrPassword

		// 更新错误次数
		if curIpLimit == nil {
			curIpLimit = &IpLimit{
				IP:            ctx.Context().ClientIP(),
				LastErrorTime: time.Now().Unix(),
				ErrorTimes:    0,
			}
			ipLimits = append(ipLimits, curIpLimit)
		}
		curIpLimit.LastErrorTime = time.Now().Unix()
		curIpLimit.ErrorTimes += 1
		if b, e := json.Marshal(&ipLimits); e == nil {
			v.LoginIpLimit = string(b)
			s.db.Model(&v).Updates(&v)
		}

		// UserLoginLog
		//l := model.UserLoginLog{
		//	Id:       idgen.Next(),
		//	UserId:   v.Id,
		//	DeviceId: req.DeviceId,
		//	Type:     1, // 登录
		//	Status:   2, // 失败
		//	LoginIp:  ctx.Context().ClientIP(),
		//}
		//s.db.Create(&l)

		return
	}
	acc := &jwt.Account{
		Uid:        v.Id,
		DeviceName: req.DeviceName,
		Platform:   req.DevicePlatform,
	}
	if rsp, err = s.updateToken(ctx, acc); err != nil {
		return nil, err
	}

	// 清空错误次数
	if curIpLimit != nil {
		curIpLimit.LastErrorTime = 0
		curIpLimit.ErrorTimes = 0
	}

	if b, e := json.Marshal(&ipLimits); e == nil {
		v.LoginIpLimit = string(b)
		v.LastLoginIp = ctx.Context().ClientIP()
		v.LastLoginTime = time.Now()
		s.db.Model(&v).Updates(&v)
	}

	// UserLoginLog
	//l := model.UserLoginLog{
	//	Id:       idgen.Next(),
	//	UserId:   v.Id,
	//	DeviceId: req.DeviceId,
	//	Type:     1,
	//	Status:   1,
	//	LoginIp:  ctx.Context().ClientIP(),
	//}
	//s.db.Create(&l)

	return
}

func (s *Service) Logout(ctx *context.Context) (err error) {
	uid := ctx.GetUid()
	client := dao.GetRedisClient()

	platform := ctx.Context().Request.Header.Get("DevicePlatform")

	tokenKey := fmt.Sprintf(constant.RedisKeyToken, 1, platform, uid)
	refreshTokenKey := fmt.Sprintf(constant.RedisKeyRefreshToken, 1, platform, uid)
	client.Del(tokenKey)
	client.Del(refreshTokenKey)

	return
}


func (s *Service) SetPwd(ctx *context.Context, req *proto.SetPwdRequest) (rsp *proto.TokenInfo, err error) {
	uid := ctx.GetUid()

	v := model.User{Id: uid}
	if err = s.db.First(&v).Error; err != nil {
		if util.IsNotFound(err) {
			err = errors.New("用户不存在")
		}
		return
	}

	req.Passwd = strings.TrimSpace(req.Passwd)
	req.Passwd = strings.ToLower(req.Passwd)

	passwdHash, err := bcrypt.GenerateFromPassword([]byte(req.Passwd), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return
	}
	v.Passwd = string(passwdHash)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&v).Updates(&v).Error; err != nil {
			return err
		}
		acc := &jwt.Account{
			Uid:        v.Id,
			DeviceName: ctx.GetDeviceName(),
			Platform:   ctx.GetPlatform(),
		}

		rsp, err = s.updateToken(ctx, acc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return
}

func (s *Service) ChangePwd(ctx *context.Context, req *proto.ChangePwdRequest) (rsp *proto.TokenInfo, err error) {
	uid := ctx.GetUid()

	v := model.User{Id: uid}
	if err = s.db.First(&v).Error; err != nil {
		if util.IsNotFound(err) {
			err = errors.ErrUserNotExists
		}
		return
	}

	req.OldPasswd = strings.TrimSpace(req.OldPasswd)
	req.OldPasswd = strings.ToLower(req.OldPasswd)

	req.NewPasswd = strings.TrimSpace(req.NewPasswd)
	req.NewPasswd = strings.ToLower(req.NewPasswd)

	if err = bcrypt.CompareHashAndPassword([]byte(v.Passwd), []byte(req.OldPasswd)); err != nil {
		err = errors.ErrPassword
		return
	}

	passwdHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPasswd), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return
	}
	v.Passwd = string(passwdHash)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Updates(&v).Error; err != nil {
			return err
		}

		acc := &jwt.Account{
			Uid:        v.Id,
			DeviceName: ctx.GetDeviceName(),
			Platform:   ctx.GetPlatform(),
		}

		rsp, err = s.updateToken(ctx, acc)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	return
}

func (s *Service) setToken(tokenKey, token, refreshTokenKey, refreshToken string) (err error) {
	client := dao.GetRedisClient()
	_, err = client.Pipelined(func(pipe redis.Pipeliner) error {
		if err := pipe.MSet(tokenKey, token, refreshTokenKey, refreshToken).Err(); err != nil {
			return err
		}
		if err := pipe.Expire(tokenKey, time.Duration(TokenExpiresIn)*time.Second).Err(); err != nil {
			return err
		}
		if err := pipe.Expire(refreshTokenKey, time.Duration(RefreshTokenExpiresIn)*time.Second).Err(); err != nil {
			return err
		}
		return nil
	})
	return
}

func (s *Service) register(ctx *context.Context, mobile string, info proto.DeviceInfo) (*proto.TokenInfo, error) {
	var (
		tokenInfo *proto.TokenInfo
	)

	// 执行本地事务
	err := s.db.Transaction(func(tx *gorm.DB) error {
		v := model.User{Mobile: mobile}
		if err := s.db.Where(&v).First(&v).Error; err != nil {
			if !util.IsNotFound(err) {
				log.Debug(err)
				return err
			}
		} else {
			return errors.New("用户已经存在，不能注册")
		}

		uid := idgen.Next()


		u := model.User{
			Id:             uid,
			Mobile:         mobile,
			MobileVerified: 1,
			LuoboId:        "",
			Nickname:       util.RandomNickname(),
			Avatar:         util.RandAvatar(),
			LastLoginTime:  time.Now(),
			LastLoginIp:    ctx.Context().ClientIP(),
			RegisterIp:     ctx.Context().ClientIP(),
		}
		if err := tx.Create(&u).Error; err != nil {
			return err
		}


		token, err := jwt.Encode(
			s.authPrivateKey,
			&jwt.Account{
				Uid:        uid,
				DeviceName: info.DeviceName,
				Platform:   info.DevicePlatform,
			},
			TokenExpiresIn,
		)
		if err != nil {
			return err
		}

		refreshToken, err := jwt.Encode(
			s.authPrivateKey,
			&jwt.Account{
				Uid:        uid,
				DeviceName: info.DeviceName,
				Platform:   info.DevicePlatform,
			},
			RefreshTokenExpiresIn,
		)
		if err != nil {
			return err
		}
		tokenKey := fmt.Sprintf(constant.RedisKeyToken, 1, info.DevicePlatform, uid)
		refreshTokenKey := fmt.Sprintf(constant.RedisKeyRefreshToken, 1, info.DevicePlatform, uid)
		if err := s.setToken(tokenKey, token, refreshTokenKey, refreshToken); err != nil {
			return err
		}

		tokenInfo = &proto.TokenInfo{
			Uid:          u.Id,
			Token:        token,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Unix() + TokenExpiresIn,
		}

		ud := model.UserDevice{
			Id:             idgen.Next(),
			Appid:          1,
			Uid:            uid,
			DevicePlatform: info.DevicePlatform,
			DeviceId:       info.DeviceId,
			DeviceName:     info.DeviceName,
			LastLoginTime:  time.Now(),
		}

		if err := tx.Create(&ud).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
}

func (s *Service) updateToken(ctx *context.Context, acc *jwt.Account) (tokenInfo *proto.TokenInfo, err error) {
	var (
		token        string
		refreshToken string
	)
	token, err = jwt.Encode(
		s.authPrivateKey,
		acc,
		TokenExpiresIn,
	)
	if err != nil {
		return
	}

	refreshToken, err = jwt.Encode(
		s.authPrivateKey,
		acc,
		RefreshTokenExpiresIn,
	)
	if err != nil {
		return
	}

	uid := acc.Uid
	tokenInfo = &proto.TokenInfo{
		Uid:          uid,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Unix() + TokenExpiresIn,
	}

	tokenKey := fmt.Sprintf(constant.RedisKeyToken, 1, acc.Platform, acc.Uid)
	refreshTokenKey := fmt.Sprintf(constant.RedisKeyRefreshToken, 1, acc.Platform, acc.Uid)
	err = s.setToken(tokenKey, token, refreshTokenKey, refreshToken)

	return
}
