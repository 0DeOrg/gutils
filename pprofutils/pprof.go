package pprofutils

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
)

/**
 * @Author: lee
 * @Description:
 * @File: pprof
 * @Date: 2023-01-16 1:59 下午
 */

func InitPProf(cfg *PProfConfig) {
	if !cfg.Enable {
		return
	}

	perfMux := http.NewServeMux()
	perfMux.HandleFunc("/debug/pprof/", pprof.Index)
	perfMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	perfMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	perfMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	perfMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	srv := http.Server{
		Handler: perfMux,
		Addr:    addr,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("InitPProf fatal, err: ", err.Error(), addr)
		}
	}()

	log.Println("pprof enabled, addr: ", addr)
}
