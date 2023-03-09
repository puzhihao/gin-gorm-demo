package impl

import (
	"database/sql"
	"github.com/puzhihao/gin-gorm-demo/conf"
)

var Service *impl = &impl{}

type impl struct {
	//数据库模块
	db *sql.DB
}

func (i *impl) Init() error {
	db, err := conf.C().Mysql.GetDB()
	if err != nil {
		return err
	}
	i.db = db
	return nil
}
