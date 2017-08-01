package common

const (
	_ = iota * -1
	ERR_USER_EXIST
)

var ERRMAP = map[int]string{
	ERR_USER_EXIST: "user name is exist",
}
