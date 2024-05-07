package conf

import (
	"github.com/mx5566/logm"
	"gopkg.in/yaml.v3"
	"net"
	"os"
)

type DB struct {
	Ip           string `yaml:"ip"`
	Name         string `yaml:"name"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	MaxIdleConns int32  `yaml:"maxIdleConns"`
	MaxOpenConns int32  `yaml:"maxOpenConns"`
	Port         uint16 `yaml:"port"`
}

type Server struct {
	Ip   string `yaml:"ip"`
	Port uint16 `yaml:"port"`
}

type ServiceEtcd struct {
	EndPoints []string `yaml:"endpoints"`
	GrantTime int64    `yaml:"granttime"`
}

type ModuleEtcd struct {
	EndPoints []string `yaml:"endpoints"`
	GrantTime int64    `yaml:"granttime"`
}

type MailBoxEtcd struct {
	EndPoints []string `yaml:"endpoints"`
	GrantTime int64    `yaml:"granttime"`
}

type Nats struct {
	EndPoints []string `yaml:"endpoints"`
}

type ModuleP struct {
	ModuleCount map[string]int64 `yaml:"module_count"`
}

func ReadConf(path string, data interface{}) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		logm.FatalfE("解析config.yaml读取错误: %s", err.Error())
		return false
	}

	err = yaml.Unmarshal(content, data)
	if err != nil {
		logm.FatalfE("解析config.yaml出错: %s", err.Error())
		return false
	}

	return true
}

func GetLanAddr(ip string) string {
	if ip == "0.0.0.0" {
		addrs, _ := net.InterfaceAddrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					return ip
				}
			}
		}
	}
	return ip
}
