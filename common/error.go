package common

const (
	Err_OK     = 0
	Err_Params           = 1002
	Err_Method_Not_EXIST = 1003
	Err_Unknow = 9999
)

var errDesc = map[int]string{
	Err_OK:     "ok",
	Err_Params:           "Params invalid",
	Err_Method_Not_EXIST: "Method doesnot exist",
	Err_Unknow: "Unknow error",
}

func GetErrDesc(errorCode int) string {
	desc, _ := errDesc[errorCode]
	return desc
}
