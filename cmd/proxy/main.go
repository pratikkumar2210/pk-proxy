package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/pratikkumar2201/pk-proxy/internal/config"
	"github.com/pratikkumar2201/pk-proxy/internal/middleware"
	"github.com/pratikkumar2201/pk-proxy/internal/proxy"
)

// func createDirectoryIfNotExists(path string) error {
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		return os.Mkdir(path, 0755)
// 	}
// 	return nil
// }

// func createLogDir() {
// 	createDirectoryIfNotExists("/Users/pratikkumar/Desktop/dumps/pk-proxy/logs")
// }

func main() {
	// createLogDir()

	store := proxy.NewProxy()

	loader := config.NewLoader(store)
	if err := loader.Load(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Handle("/", store)
		err := http.ListenAndServe(":80", middleware.Chain(store, middleware.WithSecurityHeaders))
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
