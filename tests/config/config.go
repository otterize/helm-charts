package config

import (
	"github.com/spf13/viper"
	"strings"
)

const (
	KubeConfigPath        = "kube-config-path"
	KubeConfigPathDefault = ".kube/config"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetDefault(KubeConfigPath, KubeConfigPathDefault)
	viper.AutomaticEnv()
}
