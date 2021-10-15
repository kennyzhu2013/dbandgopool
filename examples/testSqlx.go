package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// sqlx官方文档: http://jmoiron.github.io/sqlx/.
var schema = `
CREATE TABLE userinfo (
    uid INT(10) NOT NULL AUTO_INCREMENT,
    create_time datetime DEFAULT NULL,
    username VARCHAR(64)  DEFAULT NULL,
    password VARCHAR(32)  DEFAULT NULL,
    department VARCHAR(64)  DEFAULT NULL,
    email varchar(64) DEFAULT NULL,
    PRIMARY KEY (uid)
)ENGINE=InnoDB DEFAULT CHARSET=utf8
`
//-- CREATE TABLE place (
//--     country  VARCHAR(30),
//--     city  VARCHAR(30),
//--     telcode INTEGER
//-- )

func connectMysql() (*sqlx.DB) {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddrees, port, dbName, charset)
	Db, err := sqlx.Open("mysql", "root:L0v1!@#$@tcp(10.153.90.12:3306)/mid_bussops_prod?charset=utf8")
	if err != nil {
		fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
	}
	return Db
}

func addRecord(Db *sqlx.DB) {
	for i:=0; i<2; i++ {
		result, err := Db.Exec("insert into userinfo  values(?,?,?,?,?,?)",0, "2019-07-06 11:45:20", "johny", "123456", "技术部", "123456@163.com")
		if err != nil {
			fmt.Printf("data insert faied, error:[%v]", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		fmt.Printf("insert success, last id:[%d]\n", id)
	}
}

func updateRecord(Db *sqlx.DB){
	//更新uid=1的username
	result, err := Db.Exec("update userinfo set username = 'anson' where uid = 1")
	if err != nil {
		fmt.Printf("update faied, error:[%v]", err.Error())
		return
	}
	num, _ := result.RowsAffected()
	fmt.Printf("update success, affected rows:[%d]\n", num)
}

func deleteRecord(Db *sqlx.DB){
	//删除uid=2的数据
	result, err := Db.Exec("delete from userinfo where uid = 2")
	if err != nil {
		fmt.Printf("delete faied, error:[%v]", err.Error())
		return
	}
	num, _ := result.RowsAffected()
	fmt.Printf("delete success, affected rows:[%d]\n", num)
}

// sqlx test.
func main() {
	var Db *sqlx.DB = connectMysql()
	defer Db.Close()
	Db.MustExec(schema)

	addRecord(Db)
	updateRecord(Db)
	deleteRecord(Db)
}