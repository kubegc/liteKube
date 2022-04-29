package config

type TLSConfig struct {
	CAFile         string
	CAKeyFile      string
	ServerCertFile string
	ServerKeyFile  string
	ClientCertFile string
	ClientKeyFile  string
}
