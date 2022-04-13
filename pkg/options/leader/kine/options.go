package kine

import (
	"fmt"
	"sort"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type PrintFunc func(format string, a ...interface{}) error

type KineOptions struct {
	BindAddress    string `yaml:"bind-address"`
	SecurePort     int16  `yaml:"secure-port"`
	CACert         string `yaml:"ca-cert"`
	ServerCertFile string `yaml:"server-cert-file"`
	ServerkeyFile  string `yaml:"server-key-file"`
}

var defaultKO KineOptions = KineOptions{
	BindAddress: "127.0.0.1",
	SecurePort:  2379,
}

func NewKineOptions() *KineOptions {
	options := defaultKO
	return &options
}

func (opt *KineOptions) HelpSection() *help.Section {
	section := help.NewSection("kine", "lite-Database for litekube", nil)

	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port.", defaultKO.BindAddress)
	section.AddTip("secure-port", "int16", "The port on which to serve HTTPS with authentication and authorization. It cannot be switched off with 0.", fmt.Sprintf("%d", defaultKO.SecurePort))
	section.AddTip("ca-cert", "string", "SSL Certificate Authority file used to secure kine communication.", defaultKO.CACert)
	section.AddTip("server-cert-file", "string", "SSL certification file used to secure kine communication.", defaultKO.ServerCertFile)
	section.AddTip("server-key-file", "string", "SSL key file used to secure etcd communication.", defaultKO.ServerkeyFile)
	return section
}

// print all flags
func (opt *KineOptions) PrintFlags(prefix string, printFunc func(format string, a ...interface{}) error) error {
	// print flags
	flags, err := common.StructToMap(opt)
	if err != nil {
		return err
	}
	printMap(flags, prefix, printFunc)
	return nil
}

func printMap(m map[string]string, prefix string, printFunc PrintFunc) {
	if m == nil {
		return
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		printFunc("--%s-%s=%s", prefix, key, m[key])
	}
}
