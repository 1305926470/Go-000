package biz

import (
	pb "Week04/api/hello/v1"
	"Week04/internal/dao"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.HelloServer), new(*Biz)))

// 这个会被 wire 注入到 grpc 和 http 的 New 方法中。
// 需实现 .pb.go 文件中的接口
// grpc 和 wire 的这种规避反射的思路值得学习下
// 虽然增加了总体代码量，但是都是自动化生成的，也无妨。
type Biz struct {
	ac  *paladin.Map
	dao dao.Dao
}

// 依赖的 Dao 访问数据层会在 wire_gen.go 生成的代码中注入。
// 没有依赖动态实现而是通过 wire 工具自动化生成依赖注入代码
func New(d dao.Dao) (s *Biz, cf func(), err error) {
	s = &Biz{
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

// SayHello grpc func.
func (s *Biz) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)

	// 通过 s.dao 实现对数据层的访问
	fmt.Printf("hello %s", req.Name)
	return
}

func (s *Biz) Close() {}
