package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pratikkumar2201/pk-proxy/internal/proxy"
	"github.com/pratikkumar2201/pk-proxy/internal/tree"
	"github.com/pratikkumar2201/pk-proxy/models"
	"gopkg.in/yaml.v3"
)

const (
	rootDir = "/Users/pratikkumar/Desktop/dumps/pk-proxy/internal/config/samples"
)

// -----CONFIG STORE-----
type Loader struct {
	proxy *proxy.Proxy
}

func NewLoader(proxy *proxy.Proxy) *Loader {
	return &Loader{
		proxy: proxy,
	}
}

func (ld *Loader) loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config *models.ServerConfig
	yamlErr := yaml.Unmarshal(data, &config)
	if yamlErr != nil {
		return yamlErr
	}

	ld.loadServer(config)
	return nil
}

func (ld *Loader) loadServer(config *models.ServerConfig) {
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

	ld.proxy.Servers[config.Name] = &proxy.Server{
		Routes: routes,
		Tree:   tree.NewRadixTree(),
	}

	fmt.Printf("Loaded Server %s \n", config.Name)
	ld.buildRouteTree(config)
}

func (ld *Loader) buildRouteTree(config *models.ServerConfig) {
	server, _ := ld.proxy.Servers[config.Name]

	for _, route := range config.Routes {
		fmt.Println("Inserting", route.Path)
		server.Tree.Insert(route.Path, route.Upstream)
	}
	fmt.Printf("Built Route Tree For Server %s \n", config.Name)
	server.Tree.Print(server.Tree.Root, "")
}

func (ld *Loader) LoadConfigs() error {
	enteries, err := os.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range enteries {
		if !e.IsDir() {
			fmt.Printf("Loading Server %s \n", e.Name())
			path := filepath.Join(rootDir, e.Name())
			err := ld.loadConfig(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
