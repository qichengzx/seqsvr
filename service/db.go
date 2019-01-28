// package service
package service

import "database/sql"

type MySQL struct {
	User     string `yaml:"user"`
	PassWord string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

func conn(conf MySQL) *sql.DB {
	dsn := conf.User + ":" + conf.PassWord + "@" + conf.Host + "/" + conf.Database + "?charset=utf8&loc=Asia%2FShanghai"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func (svr *Service) newDBID() (int64, error) {
	res, err := svr.db.Exec("REPLACE INTO `generator_table` (uuid) VALUES (?)", svr.uuid)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
