package netmanager

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type NetManagerOptions struct {
	RegisterOptions *NetOptions `yaml:"register"`
	JoinOptions     *NetOptions `yaml:"join"`
	Token           string      `yaml:"token"`
	NodeToken       string      `yaml:"node-token"` // read from tls/Token/node.token, value not path
}

type NetOptions struct {
	Address        string `yaml:"network-address"`
	SecurePort     uint16 `yaml:"secure-port"`
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
	Token:           "local",
	RegisterOptions: NewRegisterOptions(),
	JoinOptions:     NewJoinOptions(),
}

func NewNetManagerOptions() *NetManagerOptions {
	options := DefaultNMO
	return &options
}

func (opt *NetManagerOptions) HelpSection() *help.Section {
	section := help.NewSection("network-manager", "network register and manager component for litekube", nil)
	section.AddTip("token", "string", "token value to add hosts to network and auto load certificates and node-token. It will be ignored if you provide any certificates or value by --node-token", DefaultNMO.Token)
	section.AddTip("node-token", "string", "[Not recommended] node-token value to mark one node to network, need to be given together with join and register certificates. Instead you can only specifies a valid --token.", DefaultNMO.NodeToken)

	registerSection := help.NewSection("register", "to register and query from manager. certificates need to be given together with --node-token. Or you can only ", nil)
	registerSection.AddTip("network-address", "string", "server address.", DefaultRONO.Address)
	registerSection.AddTip("secure-port", "uint16", "serving port.", fmt.Sprintf("%d", DefaultRONO.SecurePort))
	registerSection.AddTip("ca-cert", "string", "[Not recommended] SSL Certificate Authority file used to secure communication.", DefaultRONO.CACert)
	registerSection.AddTip("client-cert-file", "string", "[Not recommended] SSL certification file used to secure communication.", DefaultRONO.ClientCertFile)
	registerSection.AddTip("client-key-file", "string", "[Not recommended] SSL key file used to secure communication.", DefaultRONO.ClientkeyFile)

	joinSection := help.NewSection("join", "to be joined and managered. certificates need to be given together with --node-token, or you can only specifies a valid --token", nil)
	joinSection.AddTip("network-address", "string", "server address.", DefaultJONO.Address)
	joinSection.AddTip("secure-port", "uint16", "serving port.", fmt.Sprintf("%d", DefaultJONO.SecurePort))
	joinSection.AddTip("ca-cert", "string", "[Not recommended] SSL Certificate Authority file used to secure communication.", DefaultJONO.CACert)
	joinSection.AddTip("client-cert-file", "string", "[Not recommended] SSL certification file used to secure communication.", DefaultJONO.ClientCertFile)
	joinSection.AddTip("client-key-file", "string", "[Not recommended] SSL key file used to secure communication.", DefaultJONO.ClientkeyFile)

	section.AddSection(registerSection)
	section.AddSection(joinSection)
	return section
}

// print all flags
func (opt *NetManagerOptions) PrintFlags(prefix string, printFunc func(format string, a ...interface{}) error) error {
	// print flags
	globalFlags, err := common.StructToMapNoRecursion(opt)
	if err != nil {
		return err
	}
	common.PrintMap(globalFlags, prefix, printFunc)

	joinFags, err := common.StructToMap(opt.JoinOptions)
	if err != nil {
		return err
	}
	common.PrintMap(joinFags, prefix+"-join", printFunc)

	registerFlags, err := common.StructToMap(opt.RegisterOptions)
	if err != nil {
		return err
	}
	common.PrintMap(registerFlags, prefix+"-register", printFunc)
	return nil
}
