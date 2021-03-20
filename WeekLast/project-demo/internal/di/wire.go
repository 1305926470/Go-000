// +build wireinject

package di

import (
	"Week04/internal/biz"
	"Week04/internal/dao"
	"Week04/internal/server/grpc"
	"Week04/internal/server/http"
	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, biz.Provider, http.New,  grpc.New, NewApp))
}