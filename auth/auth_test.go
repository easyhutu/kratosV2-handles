package auth

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	server := http.NewServer(http.Address("0.0.0.0:8000"))
	router := server.Route("/x")
	router.GET("/userInfo", nil)
	app := kratos.New(
		kratos.Server(server),
	)
	defer app.Stop()
	go app.Run()
	time.Sleep(20 * time.Second)
}
