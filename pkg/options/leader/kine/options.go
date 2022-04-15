package kine

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type KineOptions struct {
	BindAddress    string `yaml:"bind-address"`
	SecurePort     uint16 `yaml:"secure-port"`
	CACert         string `yaml:"ca-cert"`
	ServerCertFile string `yaml:"server-cert-file"`
	ServerkeyFile  string `yaml:"server-key-file"`
}

var DefaultKO KineOptions = KineOptions{
	BindAddress: "127.0.0.1",
	SecurePort:  2379,
}

func NewKineOptions() *KineOptions {
	options := DefaultKO
	return &options
}

func (opt *KineOptions) HelpSection() *help.Section {
	section := help.NewSection("kine", "lite-Database for litekube", nil)

	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port.", DefaultKO.BindAddress)
	section.AddTip("secure-port", "uint16", "The port on which to serve HTTPS with authentication and authorization. It cannot be switched off with 0.", fmt.Sprintf("%d", DefaultKO.SecurePort))
	section.AddTip("ca-cert", "string", "SSL Certificate Authority file used to secure kine communication.", DefaultKO.CACert)
	section.AddTip("server-cert-file", "string", "SSL certification file used to secure kine communication.", DefaultKO.ServerCertFile)
	section.AddTip("server-key-file", "string", "SSL key file used to secure etcd communication.", DefaultKO.ServerkeyFile)
	return section
}

// print all flags
func (opt *KineOptions) PrintFlags(prefix string, printFunc func(format string, a ...interface{}) error) error {
	// print flags
	flags, err := common.StructToMap(opt)
	if err != nil {
		return err
	}
	common.PrintMap(flags, prefix, printFunc)
	return nil
}
