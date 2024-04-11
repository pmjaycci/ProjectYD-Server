package util

const (
	Success     = 200
	BadRequest  = 400
	NotFound    = 404
	Conflict    = 409
	ServerError = 500
)
const (
	HEARTBEAT       = 1
	DUPLICATE_LOGIN = 2
	ERROR           = 3
)

const (
	LOGIN = "login_db"
	GAME  = "game_db"
	LOG   = "log_db"
)

// #region REDIS Setting
const (
	RedisIp   = "13.125.254.231"
	RedisPort = "6379"
)

// #endregion

// #region DB Setting
const (
	DbDriver = "mysql"
	DbUser   = "root"
	DbPass   = "jaycci1@"
)
const (
	LoginDB     = "login_db"
	LoginDBIp   = "13.125.254.231"
	LoginDBPort = "3306"
	GameDB      = "game_db"
	GameDBIp    = "13.125.254.231"
	GameDBPort  = "3306"
	LogDB       = "log_db"
	LogDBIp     = "13.125.254.231"
	LogDBPort   = "3306"

	MaxIdleConnect = 10
	MaxOpenConnect = 10
)

//#endregion
