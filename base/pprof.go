package base

import (
	"github.com/mx5566/logm"
	"net/http"
	_ "net/http/pprof"
)

type Pprof struct {
}

func init() {

}

func (p *Pprof) Init() {
	p.Start()
}

// http://localhost:6060/debug/pprof/
func (p *Pprof) Start() {
	defer func() {
		if err := recover(); err != nil {
			TraceCode(err)
		}
	}()

	pprofAddr := "0.0.0.0:6060"
	go func(addr string) {
		if err := http.ListenAndServe(addr, nil); err != nil {
			logm.PanicfE("Pprof server ListenAndServe: %v", err)
			return
		}

		logm.DebugfE("启动pprof成功 %s", addr)
	}(pprofAddr)

}
