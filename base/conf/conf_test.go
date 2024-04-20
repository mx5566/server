package conf

import (
	"testing"
)

type Config1 struct {
	Server      `yaml:"gate"`
	ModuleEtcd  `yaml:"moduleetcd"`
	ServiceEtcd `yaml:"etcd"`
	Nats        `yaml:"nats"`
	ModuleP     `yaml:"module"`
}

func Test1(t *testing.T) {

	config := &Config1{}

	ReadConf("../../release/config.yaml", config)

	t.Log(config)
}
