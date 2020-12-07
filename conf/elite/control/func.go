package control

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	dlog "gin-cladder/conf/elite/log"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var (
	TimeLocation *time.Location // 时间节点
	TimeFormat = "2006-01-02 15:04:05" // 时间原点
	DateFormat = "2006-01-02" // 日期原点
	LocalIP = net.ParseIP("127.0.0.1") // 默认ip设置
	devpath = "conf/dev_conf.toml" // dev 环境配置路径
	prodpath = "conf/prod_conf.toml" // prod 环境配置路径
	conf = flag.String("config","dev","string")
	// 默认的是 dev 配置
)

// 路由规则
type RulePath struct {
	Config string // 配置，使用说明模式 dev
	Path string // 路径
}

// 实例化模块操作
func InitModule() error {
	// 解析操作
	flag.Parse()
	// 对 dev 和 prod 进行分析
	r,err := matchingConfig(*conf)
	if err != nil {
		return err
	}

	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO]  config=%s\n", *conf)
	log.Printf("[INFO] %s\n", " start loading resources.")
	// 验证文件格式
	if err = ParseConfPath(r.Config,r.Path); err != nil {
		return err
	}
	// 实例化 InitViperConf
	if err = InitViperConf(); err != nil {
		return err
	}
	// 实例化 BaseConf 操作
	err = InitBaseConf(r.Path)
	if err != nil {
		return err
	}
	// 设置 location 时间
	if location, err := time.LoadLocation(ConfBase.TimeLocation); err != nil {
		return err
	} else {
		TimeLocation = location
	}

	log.Printf("[INFO] %s\n", " success loading resources.")
	log.Println("------------------------------------------------------------------------")
	return nil
}


// 对日志进行销毁
func Destroy() {
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] %s\n", " start destroy resources.")
	dlog.Close()
	log.Printf("[INFO] %s\n", " success destroy resources.")
}

// 验证设置的配置模式 ，就是 flag 对其操作
func matchingConfig(conf string) (*RulePath,error){
	switch conf {
	case "dev":
		return &RulePath{"dev",devpath},nil
	case "prod":
		return &RulePath{"prod",prodpath},nil
	}
	return nil,errors.New("There is no such pattern")
}


func NewTrace() *TraceContext {
	trace := &TraceContext{}
	trace.TraceId = GetTraceId()
	trace.SpanId = NewSpanId()
	return trace
}

func NewSpanId() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}

func GetTraceId() (traceId string) {
	return calcTraceId(LocalIP.String())
}

func calcTraceId(ip string) (traceId string) {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIP := net.ParseIP(ip)
	if netIP == nil {
		b.WriteString("00000000")
	} else {
		b.WriteString(hex.EncodeToString(netIP.To4()))
	}
	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0") // go sign

	return b.String()
}