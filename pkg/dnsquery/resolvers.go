package dnsquery

import (
	"embed"
	"math/rand"
	"time"
)

//go:embed resolvers.yaml
var _resolvers embed.FS

type Resolver struct {
	Name   string `yaml:"name"`
	Server string `yaml:"server"`
}

type Resolvers []Resolver

func (r Resolvers) PickRandom() Resolver {
	if len(r) == 0 {
		return Resolver{}
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r[rnd.Intn(len(r))]
}
