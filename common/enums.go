/**
 * @Title  公用包
 * @Description 公用的属性和方法
 * @Author YaoWeiXin
 * @Update 2020/11/20 10:07
 */
package common

var MainApiType = map[string]string{
	"chat.msg": "ApiChatMsg",
}

var ThreadApiType = map[string]string{
	"chat.msg": "ApiChatMsg",
}

type ApiProcessMode string

const (
	ApiModeMain   ApiProcessMode = "Main"
	ApiModeThread ApiProcessMode = "Thread"
	ApiModeNone   ApiProcessMode = "None"
)

func ContainsApi(api string) (bool, ApiProcessMode) {
	_, ok := MainApiType[api]
	if ok {
		return true, ApiModeMain
	}
	_, ok = ThreadApiType[api]
	if ok {
		return true, ApiModeThread
	}
	return false, ApiModeNone
}

func RetrieveApi(api string) (bool, string) {
	fun, ok := MainApiType[api]
	return ok, fun
}
