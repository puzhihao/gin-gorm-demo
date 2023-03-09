package conf

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

var db *sql.DB

// conf pkg的全局变量
// 全局配置对象
var global *Config

// 全局配置对象的访问方式
func C() *Config {
	if global == nil {
		panic("config not found")
	}
	return global
}

// 全局配置对象的设置方式
func SetGlobalConfig(conf *Config) {
	global = conf
}
func NewDefaultConfig() *Config {
	return &Config{
		App:   newDefaultApp(),
		Mysql: newDefaultMySQL(),
	}
}

// 配置通过对象来进行映射
// 我们定义的是配置对象的数据结构
type Config struct {
	App   *App   `toml:"app"'`
	Mysql *Mysql `toml:"mysql"`
}

// 程序配置
type App struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
	Port string `toml:"port"`
	Key  string `toml:"key"'`
}

func newDefaultApp() *App {
	return &App{
		Name: "restful-api",
		Port: "8080",
		Host: "127.0.0.1",
		Key:  "default",
	}
}

// mysql数据库配置
type Mysql struct {
	Host     string `toml:"host"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	Port     string `toml:"port"`
	lock     sync.Mutex
}

func newDefaultMySQL() *Mysql {
	return &Mysql{
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "root",
		Password: "123456",
		Database: "restful_api",
	}
}

// DB 利用mysql配置构造全局mysql单例模式
func (m *Mysql) DB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.Username, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}

func (m *Mysql) GetDB() (*sql.DB, error) {
	//加载全局单例
	m.lock.Lock()
	defer m.lock.Unlock()
	if db == nil {
		conn, err := m.DB()
		if err != nil {
			return nil, err
		}
		db = conn
	}
	return db, nil

}
