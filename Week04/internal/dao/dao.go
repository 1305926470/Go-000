package dao

import (
	"context"
	"github.com/go-kratos/kratos/pkg/cache/memcache"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	"github.com/google/wire"
	"time"
)

var Provider = wire.NewSet(New, NewDB, NewRedis, NewMC)

type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
}

// New new a dao and return.
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d Dao, cf func(), err error) {
	return newDao(r, mc, db)
}

type dao struct {
	db          *sql.DB
	redis       *redis.Redis
	mc          *memcache.Memcache
	cache *fanout.Fanout
	demoExpire int32
}

func (d dao) Close() {
	panic("implement me")
}

func (d dao) Ping(ctx context.Context) (err error) {
	panic("implement me")
}

func newDao(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d *dao, cf func(), err error) {
	var cfg struct{
		DemoExpire time.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &dao{
		db: db,
		redis: r,
		mc: mc,
		cache: fanout.New("cache"),
		demoExpire: int32(time.Duration(cfg.DemoExpire) / time.Second),
	}
	cf = d.Close
	return
}

func NewDB() (db *sql.DB, cf func(), err error) {
	panic("implement me")
}

func NewRedis() (r *redis.Redis, cf func(), err error) {
	panic("implement me")
}

func NewMC() (mc *memcache.Memcache, cf func(), err error) {
	panic("implement me")
}