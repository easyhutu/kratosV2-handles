package auth

import (
	"fmt"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-redis/redis"
	"github.com/parnurzeal/gorequest"
	"testing"
	"time"
)

type UInfo struct {
	Uid   int64  `json:"uid"`
	Token string `json:"token"`
}

func TestAuth(t *testing.T) {
	svr := http.NewServer(http.Address(":8000"), http.Timeout(10*time.Second))
	au := New(&Options{
		Redis: &redis.Options{Addr: "0.0.0.0:6379"},
	})
	router := svr.Route("/x")
	router.GET("/userNeedLogin", GetUserInfo, au.User())
	router.GET("/user", GetUserInfo, au.Guest())
	router.POST("/login", au.RegisterLoginHandle(Login))
	router.POST("/logout", au.RegisterLogoutHandle(Logout), au.User())

	app := kratos.New(
		kratos.Server(svr),
	)
	defer app.Stop()
	go app.Run()
	time.Sleep(1 * time.Second)
	cli := gorequest.New()
	_, body, _ := cli.Post("http://127.0.0.1:8000/x/login").End()
	println(fmt.Sprintf("body %v", body))
	_, body, _ = cli.Get("http://127.0.0.1:8000/x/userNeedLogin?token=sdfdfdfdfd").Set("token", "aaaaaa").End()
	println(fmt.Sprintf("body %v", body))
	_, body, _ = cli.Get("http://127.0.0.1:8000/x/user").Set("token", "aaaaaa").End()
	println(fmt.Sprintf("body %v", body))
	_, body, _ = cli.Post("http://127.0.0.1:8000/x/logout").Set("token", "aaaaaa").End()
	println(fmt.Sprintf("body %v", body))
	_, body, _ = cli.Get("http://127.0.0.1:8000/x/userNeedLogin").Set("token", "aaaaaa").End()
	println(fmt.Sprintf("body %v", body))

}

func GetUserInfo(ctx http.Context) error {
	u, _ := FromContext(ctx)
	return ctx.JSON(200, u)

}

func Login(ctx http.Context) (*UserInfo, error) {
	info := &UInfo{Uid: 22, Token: "aaaaaa"}
	return &UserInfo{Token: "aaaaaa", Info: info}, ctx.JSON(200, &UserInfo{Token: "aaaaaa", Info: info})
}

func Logout(ctx http.Context) error {
	ctx.JSON(200, &struct {
		Code int `json:"code"`
	}{Code: 0})
	return nil
}
