package control

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

var (
	ConfEnvPath string
	ConfBase *BaseConf
	ConfEnv string
	ViperConfMap map[string]*viper.Viper
)


func GetBaseConf() *BaseConf {
	return ConfBase
}

//获取配置环境名
func GetConfEnv() string{
	return ConfEnv
}

// 验证配置文件格式操作！
func ParseConfPath(ce,cep string) error{
	path := strings.Split(cep,"/")[1]
	pr := ce + "_"
	prefix := strings.HasPrefix(path,pr)
	suffix := strings.HasSuffix(path,"toml")
	if prefix && suffix {
		ConfEnvPath = cep
		ConfEnv = ce
		return nil
	}

	return errors.New("Pattern mismatch")
}

// 实例化 InitViperConf
func InitViperConf() error {
	b , err := ioutil.ReadFile(ConfEnvPath)
	if err != nil {
		return err
	}
	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(bytes.NewBuffer(b))
	if ViperConfMap == nil {
		ViperConfMap = make(map[string]*viper.Viper)
	}

	// ViperConfMap[ConfEnv] // ConfEnv value is dev or prod
	ViperConfMap[ConfEnv] = v
	return nil
}

func GetStringConf(env,key string) string {
	// 对 env 环境变量设置对其拦截验证
	if env == ""|| key == "" {
		return ""
	}
	v, ok := ViperConfMap[env]
	if !ok {
		return ""
	}
	confString := v.GetString(key)
	return confString
}

func GetIntConf(env,key string) int {
	if env == ""|| key == "" {
		return 0
	}
	v := ViperConfMap[env]
	wt := v.GetInt(key)
	return wt
}

func GetStringSliceConf(env,key string)  []string{
	if env == "" || key == "" {
		return nil
	}
	v := ViperConfMap[env]
	wt := v.GetStringSlice(key)
	return wt

}

func ParseConfig(filepath string, conf interface{}) error {
	// 读取filepath文件，dev 模式目录为 dev_conf.toml
	data, err := ioutil.ReadFile(filepath)

	// 对其进行错误处理
	if err != nil {
		return fmt.Errorf("Read config fail, %v", err)
	}
	// new viper 对其进行操作！
	v:=viper.New()
	// 对后缀文件进行筛选
	v.SetConfigType("toml")
	// 把文件数据的buffer对其进行存储
	v.ReadConfig(bytes.NewBuffer(data))
	// 把文件数据反序列化存储在 conf.toml struct 里面
	if err =v.Unmarshal(conf); err != nil{
		return fmt.Errorf("Parse config fail, config:%v, err:%v", string(data), err)
	}

	return nil
}
