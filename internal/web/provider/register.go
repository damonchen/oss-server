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
}

// Factory provider factory
type Factory interface {
	Create(cfg *config.Configuration) ProxyProvider
}

var (
	providers         = map[string]ProxyProvider{}
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

func RegisterProxyProvider(name string, provider ProxyProvider) {
	mutex.Lock()
	defer mutex.Unlock()

	providers[name] = provider
}

func GetProxyProvider(name string) ProxyProvider {
	mutex.Lock()
	defer mutex.Unlock()
	return providers[name]
}

func IsSupportProvider(name string) bool {
	for key := range providerFactories {
		if name == key {
			return true
		}
	}
	return false
}
