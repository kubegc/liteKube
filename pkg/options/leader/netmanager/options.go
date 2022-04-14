package netmanager

import (
	"fmt"
	"sort"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type PrintFunc func(format string, a ...interface{}) error

type NetManagerOptions struct {
	RegisterOptions *NetOptions `yaml:"register"`
	JoinOptions     *NetOptions `yaml:"join"`
	MarkToken       string      `yaml:"mark-token"`
}

type NetOptions struct {
	Address        string `yaml:"network-address"`
	SecurePort     int16  `yaml:"secure-port"`
	CACert         string `yaml:"ca-cert"`
	ClientCertFile string `yaml:"client-cert-file"`
	ClientkeyFile  string `yaml:"client-key-file"`
}

var DefaultRONO NetOptions = NetOptions{
	Address:    "127.0.0.1",
	SecurePort: 6440,
}

var DefaultJONO NetOptions = NetOptions{
	Address:    "127.0.0.1",
	SecurePort: 6441,
}

func NewRegisterOptions() *NetOptions {
	options := DefaultRONO
	return &options
}

func NewJoinOptions() *NetOptions {
	options := DefaultJONO
	return &options
}

var DefaultNMO NetManagerOptions = NetManagerOptions{
	RegisterOptions: NewRegisterOptions(),
	JoinOptions:     NewJoinOptions(),
}

func NewNetManagerOptions() *NetManagerOptions {
	options := DefaultNMO
	return &options
}

func (opt *NetManagerOptions) HelpSection() *help.Section {
	section := help.NewSection("network-manager", "network register and manager component for litekube", nil)
	section.AddTip("make-token", "string", "token to indicates a host. Do not modify it later.", DefaultNMO.MarkToken)

	registerSection := help.NewSection("register", "to register and query from manager", nil)
	registerSection.AddTip("network-address", "string", "server address.", DefaultRONO.Address)
	registerSection.AddTip("secure-port", "uint16", "serving port.", fmt.Sprintf("%d", DefaultRONO.SecurePort))
	registerSection.AddTip("ca-cert", "string", "SSL Certificate Authority file used to secure communication.", DefaultRONO.CACert)
	registerSection.AddTip("client-cert-file", "string", "SSL certification file used to secure communication.", DefaultRONO.ClientCertFile)
	registerSection.AddTip("client-key-file", "string", "SSL key file used to secure communication.", DefaultRONO.ClientkeyFile)

	joinSection := help.NewSection("join", "to be joined and managered", nil)
	joinSection.AddTip("network-address", "string", "server address.", DefaultJONO.Address)
	joinSection.AddTip("secure-port", "uint16", "serving port.", fmt.Sprintf("%d", DefaultJONO.SecurePort))
	joinSection.AddTip("ca-cert", "string", "SSL Certificate Authority file used to secure communication.", DefaultJONO.CACert)
	joinSection.AddTip("client-cert-file", "string", "SSL certification file used to secure communication.", DefaultJONO.ClientCertFile)
	joinSection.AddTip("client-key-file", "string", "SSL key file used to secure communication.", DefaultJONO.ClientkeyFile)

	section.AddSection(registerSection)
	section.AddSection(joinSection)
	return section
}

// print all flags
func (opt *NetManagerOptions) PrintFlags(prefix string, printFunc func(format string, a ...interface{}) error) error {
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
