package http

import (

	pb "Week04/api/hello/v1"
	"Week04/internal/model"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

var svc pb.HelloServer

// 依赖的 pb.HelloServer 接口实现在 wire_gen.go 文件中注入
func New(s pb.HelloServer) (engine *bm.Engine, err error) {
	var (
		cfg bm.ServerConfig
		ct paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	svc = s
	engine = bm.DefaultServer(&cfg)
	pb.RegisterHelloBMServer(engine, s)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	g := e.Group("/kratos-demo")
	{
		g.GET("/start", howToStart)
	}
}

// example for http request handler.
func howToStart(c *bm.Context) {
	a := &model.Article{
		ID:      1,
		Content: "洪湖水",
		Author:  "小泽",
	}
	c.JSON(a, nil)
}