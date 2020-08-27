package redis

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/fatih/color"
	"github.com/gomodule/redigo/redis"
	"github.com/kimtaek/gamora/pkg/slack"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Configure struct {
	Mode      string `env:"APP_MODE" envDefault:"debug"`
	Host      string `env:"REDIS_HOST" envDefault:"localhost"`
	Port      string `env:"REDIS_PORT" envDefault:"6379"`
	Password  string `env:"REDIS_PASSWORD" envDefault:"XNwsLoa..."`
	Database  int    `env:"REDIS_DATABASE" envDefault:"1"`
	MaxIdle   int    `env:"REDIS_MAX_IDLE" envDefault:"20"`
	MaxActive int    `env:"REDIS_MAX_ACTIVE" envDefault:"10"`
}

var Config Configure
var p *redis.Pool

func Setup() {
	_ = env.Parse(&Config)
	conTimeout := redis.DialConnectTimeout(240 * time.Second)
	readTimeout := redis.DialReadTimeout(240 * time.Second)
	writeTimeout := redis.DialWriteTimeout(240 * time.Second)
	password := redis.DialPassword(Config.Password)
	database := redis.DialDatabase(Config.Database)
	p = &redis.Pool{
		MaxIdle:     Config.MaxIdle,
		MaxActive:   Config.MaxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", Config.Host+":"+Config.Port, password,
				readTimeout, writeTimeout, conTimeout, database)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	if err := Ping(); err != nil {
		slack.SendMessage(slack.Message{
			Text: "Redis: " + err.Error(),
		})
		os.Exit(1)
	}

	cleanupHook()
	_, _ = color.New(color.FgWhite).Println(time.Now().Format(time.RFC3339), "[info]", "[redis connected!]")
}

func Ping() error {
	conn := p.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return err
	}

	return nil
}

func cleanupHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		p.Close()
		os.Exit(0)
	}()
}

func Get(key string) ([]byte, error) {
	conn := p.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, err
	}
	return data, err
}

func Set(key string, value []byte, seconds int) error {
	conn := p.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return err
	}
	if seconds != 0 {
		_, _ = conn.Do("EXPIRE", key, seconds)
	}
	return err
}

func Exist(key string) (bool, error) {
	conn := p.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func Delete(key string) error {
	conn := p.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func Deletes(key string) error {
	conn := p.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func Incr(key string) (int, error) {
	conn := p.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", key))
}

func IncrBy(key string, value int) (int, error) {
	conn := p.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCRBY", key, value))
}

func SetExpires(key string, seconds int) {
	conn := p.Get()
	defer conn.Close()

	if seconds != 0 {
		_, _ = conn.Do("EXPIRE", key, seconds)
	}
}

func GetKeyWithPrefix(prefix string, key string) string {
	if prefix == "" {
		prefix = "api"
	}

	if Config.Mode == "release" {
		return fmt.Sprintf("production_%s_cache:%s", prefix, key)
	}

	return fmt.Sprintf("dev_%s_cache:%s", prefix, key)
}
