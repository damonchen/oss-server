package provider

import (
	"net/http"
	"sync"

	"github.com/op/go-logging"

	"github.com/damonchen/oss-server/internal/config"
)

var log = logging.MustGetLogger("provider")

type ProxyProvider interface {
	Handle(w http.ResponseWriter, req *http.Request)
	Name() string
}

// Factory provider factory
type Factory interface {
	Create(cfg *config.Configuration) []ProxyProvider
}

var (
	proxyProviders    = map[string]ProxyProvider{}
	mutex             = sync.Mutex{}
	providerFactories = map[string]Factory{}
)

func RegisterFactory(name string, factory Factory) {
	providerFactories[name] = factory
}

func GetFactory(name string) Factory {
	return providerFactories[name]
}

func GetSupportProviders() []string {
	var pvds []string
	for key := range providerFactories {
		pvds = append(pvds, key)
	}
	return pvds
}

func RegisterProxyProvider(providers []ProxyProvider) {
	mutex.Lock()
	defer mutex.Unlock()

	// TODO: 重名问题
	for _, provider := range providers {
		proxyProviders[provider.Name()] = provider
	}
}

func GetProxyProvider(name string) ProxyProvider {
	mutex.Lock()
	defer mutex.Unlock()
	return proxyProviders[name]
}

func IsSupportProvider(name string) bool {
	for key := range providerFactories {
		if name == key {
			return true
		}
	}
	return false
}

func IsSupportProxy(name string) bool {
	for key := range proxyProviders {
		if name == key {
			return true
		}
	}
	return false
}
