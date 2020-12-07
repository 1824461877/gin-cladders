package extend

import (
	"strings"
	"sync"
	"time"
)

var (
	ConnPoolShareInitServer map[string]interface{}
	ConnPoolShare *ConnPool
	once sync.Once
)

type PoolConf interface {
	Transport(*ConnPool)
}

type ConnPool struct {
	wg sync.WaitGroup
	BufPool map[string]PoolFormat
}

type PoolFormat struct {
	MaxTime time.Duration
	Conn interface{}
}


func (pf PoolFormat) maxTime() {

}


func NewPool(p PoolConf) {
	once.Do(func() {
		ConnPoolShareInitServer = make(map[string]interface{})
		ConnPoolShare = &ConnPool{BufPool:make(map[string]PoolFormat)}
	})
	if ConnPoolShare == nil {
		ConnPoolShare = &ConnPool{BufPool:make(map[string]PoolFormat)}
	}
	if ConnPoolShareInitServer == nil {
		ConnPoolShareInitServer = make(map[string]interface{})
	}
	p.Transport(ConnPoolShare)
}

func (cp *ConnPool) Set(key string,format PoolFormat) error{
	cp.wg.Add(1)
	go func() {
		cp.BufPool[key] = format
		cp.wg.Done()
	}()
	cp.wg.Wait()
	return nil
}

func (cp *ConnPool) Get(key string) interface{} {
	key = strings.ToLower(key)
	return cp.BufPool[key].Conn
}


