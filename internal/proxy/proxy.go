package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pratikkumar2201/pk-proxy/internal/tree"
	"github.com/pratikkumar2201/pk-proxy/models"
)

type Server struct {
	Routes []*models.Route
	Tree   *tree.RadixTree
}

type Proxy struct {
	Servers map[string]*Server
	client  *http.Client
}

func NewProxy() *Proxy {
	return &Proxy{
		Servers: map[string]*Server{},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ----HELPERS----
func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.ProxyErr{
		Code:    code,
		Success: false,
		Message: msg,
	})
}

// -----SERVER METHODS-----
func (s *Server) matchRoute(path string) *tree.TrieNode {
	return s.Tree.FindRoute(path)
}

// -------PROXY METHODS-------
func (px *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	px.handleRequest(w, r)
}

func (px *Proxy) getServer(host string) (*Server, error) {
	server, ok := px.Servers[host]
	if !ok {
		return nil, fmt.Errorf("host not registerd")
	}
	return server, nil
}

func (px *Proxy) buildUpstreamRequest(r *http.Request, target string) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header.Clone()
	return req, nil
}

func copyResponse(w http.ResponseWriter, res *http.Response) {
	for k, v := range res.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.Header().Add("Server", "pk-proxy")
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func (px *Proxy) handleRequest(w http.ResponseWriter, r *http.Request) {
	server, serverErr := px.getServer(r.Host)
	if serverErr != nil {
		writeError(w, http.StatusNotFound, "host not registered")
		return
	}

	// find the longest matching route in tree -> tree node, which has value as upstream of that route
	routeNode := server.matchRoute(r.URL.Path)
	fmt.Println(routeNode.Value.Servers, "<<<<routeNoderouteNode")
	if routeNode == nil {
		writeError(w, http.StatusNotFound, "route not registered")
		return
	}

	target := routeNode.Value.Servers[0] + r.URL.RequestURI()

	req, err := px.buildUpstreamRequest(r, target)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to build request")
		return
	}

	res, err := px.client.Do(req)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to reach upstream")
		return
	}

	defer res.Body.Close()
	copyResponse(w, res)
}
