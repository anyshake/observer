package unibuild

import (
	"embed"
	"encoding/json"
)

//go:embed map.json
var _map embed.FS

func extractMapValue[T any](m map[string]any, key string) (T, bool) {
	var n T
	obj, ok := m[key]
	if !ok {
		return n, false
	}
	n, ok = obj.(T)
	if !ok {
		return n, false
	}
	return n, true
}

func getToolchainById(id string) *Toolchain {
	if id == "" {
		return nil
	}

	dataBytes, err := _map.ReadFile("map.json")
	if err != nil {
		return nil
	}

	var rawMap map[string]map[string]any
	if err = json.Unmarshal(dataBytes, &rawMap); err != nil {
		return nil
	}

	targetToolchain, ok := rawMap[id]
	if !ok {
		return nil
	}

	name, ok := extractMapValue[string](targetToolchain, "name")
	if !ok {
		return nil
	}
	gomips, ok := extractMapValue[string](targetToolchain, "gomips")
	if !ok {
		return nil
	}
	goarm, ok := extractMapValue[string](targetToolchain, "goarm")
	if !ok {
		return nil
	}
	goarch, ok := extractMapValue[string](targetToolchain, "goarch")
	if !ok {
		return nil
	}
	goos, ok := extractMapValue[string](targetToolchain, "goos")
	if !ok {
		return nil
	}

	return &Toolchain{
		GOOS:   goos,
		GOARCH: goarch,
		GOARM:  goarm,
		GOMIPS: gomips,
		Name:   name,
	}
}
