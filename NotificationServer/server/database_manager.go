package server

import (
	"database/sql"
	"fmt"
	util "project_yd/util"

	_ "github.com/go-sql-driver/mysql"
)

type GameDatabase struct {
	Login *sql.DB
	Game  *sql.DB
	Log   *sql.DB
}

var DBManager *GameDatabase

func StartDBConnection() {
	println("Start DB Connect!!")

	//-- Login DB Open
	login := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", util.DbUser, util.DbPass, util.LoginDBIp, util.LoginDBPort, util.LoginDB)

	db, err := sql.Open(util.DbDriver, login)
	if err != nil {
		println("LoginDB::", err.Error())
	}
	db.SetMaxIdleConns(util.MaxIdleConnect)
	db.SetMaxOpenConns(util.MaxOpenConnect)
	loginDB := db

	//-- Game DB Open
	game := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", util.DbUser, util.DbPass, util.GameDBIp, util.GameDBPort, util.GameDB)

	db, err = sql.Open(util.DbDriver, game)
	if err != nil {
		println("GameDB::", err.Error())
	}
	db.SetMaxIdleConns(util.MaxIdleConnect)
	db.SetMaxOpenConns(util.MaxOpenConnect)
	gameDB := db

	//-- Log DB Open
	log := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", util.DbUser, util.DbPass, util.LogDBIp, util.LogDBPort, util.LogDB)

	db, err = sql.Open(util.DbDriver, log)
	if err != nil {
		println("LogDB::", err.Error())
	}

	db.SetMaxIdleConns(util.MaxIdleConnect)
	db.SetMaxOpenConns(util.MaxOpenConnect)
	logDB := db

	DBManager = &GameDatabase{
		Login: loginDB,
		Game:  gameDB,
		Log:   logDB,
	}
	println("End DB Connect!!")
}
