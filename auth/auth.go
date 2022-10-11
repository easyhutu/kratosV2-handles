package auth

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

const (
	DefaultTokenExpire = 30 * 24 * time.Hour
)

type Auth struct {
	redis       *redis.Client
	expire      time.Duration
	validateMsg interface{}
	tokenKey    string
}

type authKey struct{}

type Options struct {
	Redis           *redis.Options
	TokenExpire     time.Duration
	ValidateMsgTemp interface{}
	TokenKey        string
}

type LoginFunc func(ctx http2.Context) (*UserInfo, error)

func New(opt *Options) *Auth {
	if opt.TokenExpire == 0 {
		opt.TokenExpire = DefaultTokenExpire
	}
	if opt.ValidateMsgTemp == nil {
		opt.ValidateMsgTemp = noLoginMsg
	}
	if opt.TokenKey == "" {
		opt.TokenKey = "token"
	}
	return &Auth{
		redis:  redis.NewClient(opt.Redis),
		expire: opt.TokenExpire,
	}
}

func (a *Auth) PutUser(info *UserInfo) error {

	bs, err := json.Marshal(info)
	if err != nil {
		log.Errorf("auth json marshal err: %v", err)
		return err
	}
	err = a.redis.Set(info.Token, bs, a.expire).Err()
	return err
}

func (a *Auth) GetUser(token string) (*UserInfo, error) {
	bs, err := a.redis.Get(token).Bytes()
	if err != nil {
		log.Errorf("token %s, bs: %x", token, bs)
		return nil, err
	}
	u := &UserInfo{}
	if err := json.Unmarshal(bs, &u); err != nil {
		log.Errorf("auth json marshal err: %v", err)
	}
	return u, nil
}

func (a *Auth) Guest() http2.FilterFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			u, _ := a.tokenUser(request)
			if u != nil {
				ctx := context.WithValue(request.Context(), authKey{}, u)
				request = request.WithContext(ctx)
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func (a *Auth) User() http2.FilterFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			u, err := a.tokenUser(request)
			if err != nil || u == nil {
				log.Errorf("auth failed err %+v", err)
				resp, _ := json.Marshal(a.validateMsg)
				writer.Write(resp)
				return
			}
			ctx := context.WithValue(request.Context(), authKey{}, u)
			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request)
		})
	}
}

func (a *Auth) tokenUser(request *http.Request) (u *UserInfo, err error) {
	token := ""
	if request.Form.Get(a.tokenKey) != "" {
		token = request.Form.Get(a.tokenKey)
	}
	if ctoken, err := request.Cookie(a.tokenKey); err == nil {
		if ctoken.Value != "" {
			token = ctoken.Value
		}
	}
	if htoken := request.Header.Get(a.tokenKey); htoken != "" {
		token = htoken
	}
	if token == "" {
		return
	}
	return a.GetUser(token)
}

func (a *Auth) RegisterLoginHandle(loginFunc LoginFunc) http2.HandlerFunc {
	return func(c http2.Context) error {
		userInfo, err := loginFunc(c)
		if err != nil {
			log.Errorf("login error %+v", err)
			return err
		}
		if userInfo != nil {
			a.PutUser(userInfo)
		}
		return nil
	}
}

func FromContext(ctx context.Context) (u interface{}, ok bool) {
	ui, ok := ctx.Value(authKey{}).(*UserInfo)
	return ui.Info, ok
}
