package ecflow_watchman

type EcflowServerConfig struct {
	Owner          string `yaml:"owner"`
	Repo           string `yaml:"repo"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	ConnectTimeout int    `yaml:"connect_timeout"`
}
