package lib

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	logs.SetLogFuncCall(true)
	logs.SetLogger("console")
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"/opt/logs/go/fisher/def.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],"daily":true,"maxdays":365,"perm":"0755"}`)
}

func Goid() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic recover:panic info: %v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}

	return id
}

func Log(v ...interface{}) {
	logs.Info("[go-id: "+strconv.Itoa(Goid())+"]", v, "\n")
}
