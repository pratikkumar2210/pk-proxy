package models

type Upstream struct {
	Name       string   `yaml:"name" json:"name"`
	LBStrategy string   `yaml:"lb-strategy" json:"lb_strategy"`
	Servers    []string `yaml:"servers" json:"servers"`
}

type Route struct {
	Path     string    `yaml:"path" json:"path"`
	Upstream *Upstream `yaml:"upstreams" json:"upstream"`
}

type ServerConfig struct {
	Name   string   `yaml:"server" json:"server"`
	Routes []*Route `yaml:"routes" json:"routes"`
}

type ProxyLogDirConfig struct {
	Error  string `yaml:"error" json:"error"`
	Access string `yaml:"access" json:"access"`
}

type ProxyLogConfig struct {
	Format string            `yaml:"format" json:"format"`
	Dir    ProxyLogDirConfig `yaml:"dir" json:"dir"`
}

type ProxyTimeoutConfig struct {
	ClientRead      string `yaml:"client_read" json:"client_read"`
	ClientWrite     string `yaml:"client_write" json:"client_write"`
	UpstreamConnect string `yaml:"upstream_connect" json:"upstream_connect"`
	UpstreamRead    string `yaml:"upstream_read" json:"upstream_read"`
}

type ProxyConfig struct {
	Log     ProxyLogConfig     `yaml:"log" json:"log"`
	Timeout ProxyTimeoutConfig `yaml:"timeout" json:"timeout"`
}

type Server struct {
	Name   string
	Routes []*Route
}

type ProxyErr struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
