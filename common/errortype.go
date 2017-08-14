package common

const (
	SUCCESS = 0

	USERNAMEEXIST = (iota * -1) - 10000 // -10000
	PARAMSERR

	DATABASEERR = (iota * -1) - 20000 // -20000

	MIDDLEWAREERR = (iota * -1) - 30000 // -30000

	NEEDLOGIN
)

var ERRMAP = map[int]string{
	SUCCESS: "success",

	USERNAMEEXIST: "user name was exist",
	PARAMSERR:     "params error",
	DATABASEERR:   "database create user error",
	MIDDLEWAREERR: "middler ware error",
	NEEDLOGIN:     "not login",
}
