package common

const (
	SUCCESS = 0

	USER_WASEXIST = (iota * -1) - 10000 // -10000
	USER_PARSEFORM
	USER_PARAMAGET
	USER_PARAMVALID

	PAGE_NOPAGEGUID = (iota * -1) - 20000 // -20000
	PAGE_GUIDNOTFOUND

	DATA_CREATEUSER = (iota * -1) - 30000 // -30000
)

var ERRMAP = map[int]string{
	SUCCESS:           "sucess",
	USER_WASEXIST:     "user name is exist",
	USER_PARSEFORM:    "parse form error",
	USER_PARAMAGET:    "param get error",
	USER_PARAMVALID:   "param error",
	PAGE_NOPAGEGUID:   "dont have page guid",
	PAGE_GUIDNOTFOUND: "page guid not found",
	DATA_CREATEUSER:   "database create user error",
}
