/**
 * @Title  公用包
 * @Description 公用的属性和方法
 * @Author YaoWeiXin
 * @Update 2020/11/20 10:07
 */
package common

type ApiProcessMode string

const (
	ApiModeMain   ApiProcessMode = "Main"
	ApiModeThread ApiProcessMode = "Thread"
	ApiModeNone   ApiProcessMode = "None"
)
