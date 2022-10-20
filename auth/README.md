
### 使用

```go
import (
"github.com/easyhutu/kratosV2-handles/auth"
)

func main() {
    svr := http.NewServer(http.Address(":8000"), http.Timeout(10*time.Second))
    au := auth.New(&Options{
    Redis: &redis.Options{Addr: "0.0.0.0:6379"},
    })
    router := svr.Route("/x")
    router.GET("/userNeedLogin", GetUserInfo, au.User())
    router.POST("/login", au.RegisterLoginHandle(Login))
    router.POST("/logout", au.RegisterLogoutHandle(Logout), au.User())
    app := kratos.New(
    kratos.Server(svr),
    )
    defer app.Stop()
    go app.Run()
}

func GetUserInfo(ctx http.Context) error {
    u, _ := auth.FromContext(ctx)
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
```
