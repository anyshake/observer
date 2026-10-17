package dnsquery

import (
	"embed"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"gopkg.in/yaml.v2"
)

//go:embed resolvers.yaml
var _resolvers embed.FS

type Resolver struct {
	Name   string `yaml:"name"`
	Server string `yaml:"server"`
}

type Resolvers []Resolver

// NewResolvers loads only encrypted endpoints. Catalog failures never downgrade
// queries to public UDP resolvers or the system resolver.
func NewResolvers() (Resolvers, error) {
	data, err := _resolvers.ReadFile("resolvers.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read DNS resolver catalog: %w", err)
	}
	return parseResolvers(data)
}

func parseResolvers(data []byte) (Resolvers, error) {
	var resolvers Resolvers
	if err := yaml.UnmarshalStrict(data, &resolvers); err != nil {
		return nil, fmt.Errorf("failed to parse DNS resolver catalog: %w", err)
	}
	if len(resolvers) == 0 {
		return nil, fmt.Errorf("DNS resolver catalog is empty")
	}
	for _, resolver := range resolvers {
		endpoint, err := url.Parse(resolver.Server)
		if err != nil || (endpoint.Scheme != "https" && endpoint.Scheme != "tls" && endpoint.Scheme != "sdns") {
			return nil, fmt.Errorf("DNS resolver %q must use an encrypted endpoint", resolver.Name)
		}
		if _, err := New(resolver.Server); err != nil {
			return nil, fmt.Errorf("invalid DNS resolver %q: %w", resolver.Name, err)
		}
	}
	return resolvers, nil
}

func (r Resolvers) PickRandom() Resolver {
	if len(r) == 0 {
		return Resolver{}
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r[rnd.Intn(len(r))]
}
