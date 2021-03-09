package provider

import (
	"net/http"

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
	providerFactories = map[string]Factory{}
)

func RegisterFactory(name string, factory Factory) {
	providerFactories[name] = factory
}

func GetFactory(name string) (factory Factory) {
	factory = providerFactories[name]
	return
}

func RegisterProxyProvider(name string, provider ProxyProvider) {
	providers[name] = provider
}

func GetProxyProvider(name string) ProxyProvider {
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
