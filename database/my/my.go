package my

import (
	"database/sql"
	"strconv"

	"github.com/dalpengida/portfolio-api-go-mysql/config"
	"github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var dsn string

func init() {
	port, err := strconv.Atoi(config.Config(config.MYSQL_PORT))
	if err != nil {
		panic(err)
	}

	dsn = driver.MySQLBuildQueryString(
		config.Config(config.MYSQL_USER),
		config.Config(config.MYSQL_PASS),
		config.Config(config.MYSQL_DB),
		config.Config(config.MYSQL_HOST),
		port,
		config.Config(config.MYSQL_SSL_MODE),
	)
}

// SetBoilerDatabas 는 sqlboiler 에 database 접속 정보 및 초기화
func SetBoilerDatabas() error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// sqlboiler 에 관련 정보를 셋팅 , 나중에 getcontextdb를 통해서 받아서 사용가능
	boil.SetDB(db)

	return nil
}
