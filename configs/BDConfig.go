package configs

import (
	"awesomeProject/model"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //导入包但不使用，init()
	"github.com/spf13/viper"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

func buildDBConnectUrl() (dsn string) {
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	userName := viper.GetString("database.username")
	password := viper.GetString("database.password")
	database := viper.GetString("database.database")
	return userName + ":" + password + "@tcp(" + host + ":" + port + ")/" + database
}

func InitBD() (err error) {

	dsn := buildDBConnectUrl()
	fmt.Println("Connecting to database ...")
	db, err := sql.Open("mysql", dsn) //open不会检验用户名和密码
	if err != nil {
		fmt.Printf("dsn:%s invalid,err:%v\n", dsn, err)
		return err
	}
	err = db.Ping() //尝试连接数据库
	if err != nil {
		fmt.Printf("open %s faild,err:%v\n", dsn, err)
		return err
	}
	fmt.Println("database connected successfully!")
	model.DAO = db

	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		fmt.Printf("dsn:%s invalid,err:%v\n", dsn, err)
		return err
	}
	model.XDAO = engine
	engine.SetTableMapper(names.GonicMapper{})
	engine.SetColumnMapper(names.GonicMapper{})

	return nil
}

func InitMapper() {

}
