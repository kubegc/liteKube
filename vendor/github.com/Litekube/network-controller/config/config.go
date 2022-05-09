/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: wanna <wananzjx@163.com>
 *
 */
package config

import (
	"errors"
	"github.com/Litekube/network-controller/contant"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

// server.yml / client.yml

// Server Config
type ServerConfig struct {
	Ip             string `yaml:"ip"`
	Port           int    `yaml:"port"`
	NetworkCertDir string `yaml:"networkCertDir"`
	BootstrapPort  int    `yaml:"bootstrapPort"`
	GrpcPort       int    `yaml:"grpcPort"`
	GrpcCertDir    string `yaml:"grpcCertDir"`

	NetworkAddr string `yaml:"networkAddr"`
	//DbPath          string `yaml:"dbPath"`
	LogDir          string `yaml:"logDir"`
	WorkDir         string `yaml:"workDir"`
	TlsDir          string `yaml:"tlsDir"`
	Debug           bool   `yaml:"debug"`
	MTU             int    `yaml:"mtu"`
	Interconnection bool   `yaml:"interconnection"`

	NetworkCAFile         string
	NetworkCAKeyFile      string
	NetworkServerCertFile string
	NetworkServerKeyFile  string

	GrpcCAFile         string
	GrpcCAKeyFile      string
	GrpcServerCertFile string
	GrpcServerKeyFile  string
}

// Client Config
type ClientConfig struct {
	NetworkCertDir  string `yaml:"networkCertDir"`
	ServerAddr      string `yaml:"serverAddr"`
	Port            int    `yaml:"port"`
	LogDir          string `yaml:"logDir"`
	WorkDir         string `yaml:"workDir"`
	Debug           bool   `yaml:"debug"`
	MTU             int    `yaml:"mut"`
	Token           string `yaml:"token"`
	RedirectGateway bool   `yaml:"redirectGateway"`

	CAFile         string
	ClientCertFile string
	ClientKeyFile  string
}

type NetworkConfig struct {
	Mode   string       `yaml:"mode"`
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

// return server/client config
func ParseConfig(filename string) (interface{}, error) {
	cfg := &NetworkConfig{}

	File, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("fail to read file: %v", err)
	}
	err = yaml.Unmarshal(File, &cfg)
	if err != nil {
		return nil, err
	}

	switch cfg.Mode {
	case "server":
		cfg.Server.NetworkCAFile = filepath.Join(cfg.Server.NetworkCertDir, contant.CAFile)
		cfg.Server.NetworkCAKeyFile = filepath.Join(cfg.Server.NetworkCertDir, contant.CAKeyFile)
		cfg.Server.NetworkServerCertFile = filepath.Join(cfg.Server.NetworkCertDir, contant.ServerCertFile)
		cfg.Server.NetworkServerKeyFile = filepath.Join(cfg.Server.NetworkCertDir, contant.ServerKeyFile)

		cfg.Server.GrpcCAFile = filepath.Join(cfg.Server.GrpcCertDir, contant.CAFile)
		cfg.Server.GrpcCAKeyFile = filepath.Join(cfg.Server.GrpcCertDir, contant.CAKeyFile)
		cfg.Server.GrpcServerCertFile = filepath.Join(cfg.Server.GrpcCertDir, contant.ServerCertFile)
		cfg.Server.GrpcServerKeyFile = filepath.Join(cfg.Server.GrpcCertDir, contant.ServerKeyFile)
		return cfg.Server, nil
	case "client":
		cfg.Client.CAFile = filepath.Join(cfg.Client.NetworkCertDir, contant.CAFile)
		cfg.Client.ClientCertFile = filepath.Join(cfg.Client.NetworkCertDir, contant.ClientCertFile)
		cfg.Client.ClientKeyFile = filepath.Join(cfg.Client.NetworkCertDir, contant.ClientKeyFile)
		return cfg.Client, nil
	default:
		return nil, errors.New("Wrong config data")
	}
}
