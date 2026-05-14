package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pratikkumar2201/pk-proxy/internal/proxy"
	"github.com/pratikkumar2201/pk-proxy/internal/tree"
	"github.com/pratikkumar2201/pk-proxy/models"
	"github.com/pratikkumar2201/pk-proxy/pkg/logger"
	"gopkg.in/yaml.v3"
)

const (
	rootDir  = "/Users/pratikkumar/Desktop/dumps/pk-proxy/internal/config/samples"
	confPath = "/Users/pratikkumar/Desktop/dumps/pk-proxy/pk-proxy.yaml"
)

type Loader struct {
	proxy *proxy.Proxy
}

func NewLoader(proxy *proxy.Proxy) *Loader {
	return &Loader{
		proxy: proxy,
	}
}

func (ld *Loader) loadProxyConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("load proxy config: %w", err)
	}
	var config *models.ProxyConfig
	yamlErr := yaml.Unmarshal(data, &config)
	if yamlErr != nil {
		return fmt.Errorf("load proxy config: %w", yamlErr)
	}
	ld.proxy.Config = config
	return nil
}

func (ld *Loader) loadServerConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config *models.ServerConfig
	yamlErr := yaml.Unmarshal(data, &config)
	if yamlErr != nil {
		return yamlErr
	}

	ld.loadServer(path, config)
	return nil
}

func (ld *Loader) loadServer(path string, config *models.ServerConfig) {
	var routes []*models.Route
	for _, route := range config.Routes {
		routes = append(routes, &models.Route{
			Path: route.Path,
			Upstream: &models.Upstream{
				Name:       route.Upstream.Name,
				LBStrategy: route.Upstream.LBStrategy,
				Servers:    route.Upstream.Servers,
			},
		})
	}

	// create log file
	_, logFErr := os.Create(filepath.Join("/Users/pratikkumar/Desktop/dumps/pk-proxy/logs", config.Name))
	if logFErr != nil {
		log.Default().Printf("error creating log file %s", config.Name)
	}

	ld.proxy.Servers[config.Name] = &proxy.Server{
		Routes: routes,
		Tree:   tree.NewRadixTree(),
		Logger: logger.NewLogger(path),
	}
	ld.buildRouteTree(config)
}

func (ld *Loader) buildRouteTree(config *models.ServerConfig) {
	server, _ := ld.proxy.Servers[config.Name]

	for _, route := range config.Routes {
		// fmt.Println("Inserting", route.Path)
		server.Tree.Insert(route.Path, route.Upstream)
	}
	fmt.Printf("Built Route Tree For Server %s \n", config.Name)
	// server.Tree.Print(server.Tree.Root, "")
	log.Default().Printf("loaded server %s \n", config.Name)
}

func (ld *Loader) Load() error {
	err := ld.loadProxyConfig(confPath)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	enteries, err := os.ReadDir(rootDir)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	for _, e := range enteries {
		if !e.IsDir() {
			log.Default().Printf("loading server %s \n", e.Name())
			path := filepath.Join(rootDir, e.Name())
			err := ld.loadServerConfig(path)
			if err != nil {
				log.Default().Printf("error loading config %s %s", e.Name(), err.Error())
			}
		}
	}
	return nil
}
