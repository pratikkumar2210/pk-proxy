package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/pratikkumar2201/pk-proxy/internal/config"
	"github.com/pratikkumar2201/pk-proxy/internal/middleware"
	"github.com/pratikkumar2201/pk-proxy/internal/proxy"
)

func main() {
	store := proxy.NewProxy()

	loader := config.NewLoader(store)
	if err := loader.LoadConfigs(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Handle("/", store)
		err := http.ListenAndServe(":80", middleware.Chain(store, middleware.WithSecurityHeaders, middleware.WithLogger))
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
