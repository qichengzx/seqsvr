// package service
//
// https://mp.weixin.qq.com/s/F7WTNeC3OUr76sZARtqRjw

package service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)

var (
	err error
	db  *sql.DB
)

func InitDB(conf MySQL) {
	dsn := conf.User + ":" + conf.PassWord + "@" + conf.Host + "/" + conf.Database + "?charset=utf8&loc=Asia%2FShanghai"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

}

func New(uuid uuid.UUID) (int64, error) {
	res, err := db.Exec("REPLACE INTO `generator_table` (uuid) VALUES (?)", uuid)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
