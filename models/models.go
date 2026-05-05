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

type Server struct {
	Name   string
	Routes []*Route
}

type ProxyErr struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
