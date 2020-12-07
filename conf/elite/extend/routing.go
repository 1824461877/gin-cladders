package extend

import (
	"fmt"
	"gin-cladder/conf/elite/control"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
	InitRouter func() *gin.Engine
)

type InitServerConf struct {
	Router func() *gin.Engine
}

func (isc *InitServerConf) HttpServerRun() {
	if isc.Router == nil || InitRouter != nil {
		return
	}
	InitRouter = isc.Router
	HttpServerRun()
}

func (isc *InitServerConf) HttpServerStop() {
	if isc == nil {
		return
	}
	HttpServerStop()
}

func HttpServerRun() {
	gin.SetMode(control.GetBaseConf().DebugMode)
	HttpSrvHandler = &http.Server{
		Addr:              control.GetStringConf(control.ConfEnv,"http.addr"),
		Handler:           InitRouter(),
		ReadTimeout:    time.Duration(control.GetIntConf(control.ConfEnv,"http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(control.GetIntConf(control.ConfEnv,"http.write_timeout")) * time.Second,
		MaxHeaderBytes:    control.GetIntConf(control.ConfEnv,"http.max_header_bytes"),
	}

	go func() {
		log.Printf("[INFO] This HttpsSrverRun",control.GetStringConf(control.ConfEnv,"http.addr"))
		log.Printf("[INFO] HttpServerRun:%s\n",control.GetStringConf(control.ConfEnv,"http.addr"))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf("[ERROR] HttpServerRun:%s err:%v\n", control.GetStringConf(control.ConfEnv,"http.addr"), err)
		}
	}()
}

func HttpServerStop() {
	fmt.Println("this is stop set")
}