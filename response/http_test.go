package response

import (
	"fmt"
	"github.com/easyhutu/kratosV2-handles/ecode"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/parnurzeal/gorequest"
	"testing"
	"time"
)

func TestHttpContext(t *testing.T) {
	svr := http.NewServer(http.Address(":8000"), http.Timeout(10*time.Second))

	router := svr.Route("/x")
	router.GET("/ping", Ping)
	app := kratos.New(
		kratos.Server(svr),
	)
	defer app.Stop()
	go app.Run()
	time.Sleep(1 * time.Second)
	cli := gorequest.New()
	_, body, _ := cli.Get("http://127.0.0.1:8000/x/ping").End()
	println(fmt.Sprintf("body %v", body))
}

func Ping(ctx http.Context) error {

	return JSON(ctx, &struct {
		Val int `json:"val"`
	}{Val: 20}, ecode.New(ecode.Code500, "服务器错误"))
}
