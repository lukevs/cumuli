// db.go contains functions for managing cumuli's Redis database

package main

import (
    "time"

    "github.com/garyburd/redigo/redis"
)

// NewPool creates a new Redis pool from the given server and password.
func NewPool(server, password string) *redis.Pool {

    if password == "" {
        return &redis.Pool{
            MaxIdle: 3,
            IdleTimeout: 240 * time.Second,
            Dial: func () (redis.Conn, error) {
                return redis.Dial("tcp", server)
            },
            TestOnBorrow: func(c redis.Conn, t time.Time) error {
                _, err := c.Do("PING")
                return err
            },
        }
    }

    return &redis.Pool{
        MaxIdle: 3,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
            if _, err := c.Do("AUTH", password); err != nil {
                c.Close()
                return nil, err
            }
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
}