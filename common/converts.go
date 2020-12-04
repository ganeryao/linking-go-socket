package common

import "strings"

func ConvertApi(api string) string {
	a := strings.Split(api, ".")
	return a[len(a)-2] + "." + a[len(a)-1]
}
