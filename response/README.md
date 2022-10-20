### 使用

```go
import (
"github.com/easyhutu/kratosV2-handles/response"
)
func main() {
    svr := http.NewServer(http.Address(":8000"), http.Timeout(10*time.Second))
    
    router := svr.Route("/x")
    router.GET("/ping", Ping)
    app := kratos.New(
    kratos.Server(svr),
    )
    defer app.Stop()
    go app.Run()
}

func Ping(ctx http.Context) error {

return response.JSON(ctx, &struct {
Val int `json:"val"`
}{Val: 20}, ecode.New(ecode.Code500, "服务器错误"))
}

```

```shell
curl http:127.0.0.1:8000/x/ping
response: {"code":500,"data":{"val":20},"message":"服务器错误"}
```

