package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// 定义全局变量pool存储redis连接池
var pool *redis.Pool

// 初始化连接池
func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle: maxIdle, // 最大空闲连接数
		MaxActive: maxActive, // 和数据库的最大连接数，0 表示没有限制
		IdleTimeout: idleTimeout, // 最大空闲时间
		Dial: func () (redis.Conn, error) { // 初始化连接
			return redis.Dial("tcp", address)
		},
	}
}