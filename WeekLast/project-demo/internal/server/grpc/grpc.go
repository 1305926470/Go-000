package grpc

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	pb "project-demo/api/hello/v1"
)

// 依赖的 pb.HelloServer 接口实现在 wire_gen.go 文件中注入
func New(svc pb.HelloServer) (ws *warden.Server, err error) {
	var (
		cfg warden.ServerConfig
		ct paladin.TOML
	)
	if err = paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	ws = warden.NewServer(&cfg)
	pb.RegisterHelloServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}