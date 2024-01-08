package gameuts

import (
	"baseutils/baseuts"
	"net/http"
	_ "net/http/pprof"
)

var pprofHost = ""
var pprofPort = ""

func PPROFCheck(host string, port string) {
	pprofHost = host
	pprofPort = port
	go func() {
		defer baseuts.ChkRecover()
		runPPROFCheck()
	}()
}

// http://58.87.95.155:6060/debug/pprof

//go tool pprof -inuse_space http://58.87.95.155:6060/debug/pprof/heap
// -alloc_objects：已分配的对象总量（不管是否已释放）
// -alloc_space：已分配的内存总量（不管是否已释放）
// -inuse_objects： 已分配但尚未释放的对象数量
// -inuse_sapce：已分配但尚未释放的内存数量

//go tool pprof -inuse_space -cum -svg http://58.87.95.155:6060/debug/pprof/heap > heap_inuse.svg

// top
// top -cum 函数调用关系 中的数据进行累积

// go tool pprof http://58.87.95.155:6060/debug/pprof/profile
// -http=:6066 本地以指定端口打开profile
// -seconds=5 采样时间设定
// http://localhost:6066/ui

func runPPROFCheck() {
	pprofStr := pprofHost + ":" + pprofPort
	baseuts.Log("pprof", pprofStr)
	err := http.ListenAndServe(pprofStr, nil)
	baseuts.ChkErr(err)
}
