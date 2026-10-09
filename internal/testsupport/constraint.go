package testsupport

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

func ExerciseConstraint(t *testing.T, handler *action.Handler, constraint config.IConstraint) {
	t.Helper()
	name := constraint.GetNamespace() + "/" + constraint.GetKey()
	if constraint.GetName() == "" || constraint.GetKey() == "" || constraint.GetNamespace() == "" || constraint.GetDescription() == "" {
		t.Fatalf("%s is missing metadata", name)
	}
	_ = constraint.GetVersion()
	_ = constraint.IsRequired()
	_ = constraint.GetType()
	_ = constraint.GetOptions()

	if _, err := constraint.Get(handler); err == nil {
		t.Fatalf("%s Get before Init succeeded", name)
	}
	if err := constraint.Set(handler, struct{ Invalid bool }{}); err == nil {
		t.Fatalf("%s accepted a struct", name)
	}
	if err := constraint.Init(handler); err != nil {
		t.Fatalf("%s Init: %v", name, err)
	}
	if err := constraint.Init(handler); err != nil {
		t.Fatalf("%s second Init: %v", name, err)
	}
	if _, err := constraint.Get(handler); err != nil {
		t.Fatalf("%s Get: %v", name, err)
	}

	var setErr error
	accepted := false
	for _, candidate := range constraintCandidates(constraint) {
		setErr = constraint.Set(handler, candidate)
		if setErr != nil {
			continue
		}
		if _, err := constraint.Get(handler); err != nil {
			t.Fatalf("%s Get after Set(%#v): %v", name, candidate, err)
		}
		accepted = true
		break
	}
	if !accepted {
		t.Fatalf("%s rejected every candidate, last error: %v", name, setErr)
	}

	for _, candidate := range invalidCandidates(constraint) {
		_ = constraint.Set(handler, candidate)
	}
	if err := constraint.Restore(handler); err != nil {
		t.Fatalf("%s Restore: %v", name, err)
	}
	if _, err := constraint.Get(handler); err != nil {
		t.Fatalf("%s Get after Restore: %v", name, err)
	}
}

func constraintCandidates(constraint config.IConstraint) []any {
	candidates := []any{constraint.GetDefaultValue()}
	switch value := constraint.GetDefaultValue().(type) {
	case []string:
		boxed := make([]any, len(value))
		for i, item := range value {
			boxed[i] = item
		}
		candidates = append(candidates, boxed, []any{"EHZ"}, []any{"80/tcp"}, []any{"123/udp"})
	case []int64:
		boxed := make([]any, len(value))
		for i, item := range value {
			boxed[i] = float64(item)
		}
		candidates = append(candidates, boxed, []any{float64(1)})
	case []float64:
		boxed := make([]any, len(value))
		for i, item := range value {
			boxed[i] = item
		}
		candidates = append(candidates, boxed, []any{float64(1)})
	case int:
		candidates = append(candidates, int64(value), float64(value), json.Number(strconv.Itoa(value)), float64(1), json.Number("1"))
	case int64:
		candidates = append(candidates, float64(value), json.Number(strconv.FormatInt(value, 10)), float64(1))
	case float64:
		candidates = append(candidates, json.Number(strconv.FormatFloat(value, 'f', -1, 64)), float64(1))
	case string:
		candidates = append(candidates, "test", "Ab1", "127.0.0.1", "https://example.com", "tcp", "sta/lta")
		if value == "" {
			candidates = append(candidates, "value")
		}
	case bool:
		candidates = append(candidates, !value, true, false)
	}
	for key, option := range constraint.GetOptions() {
		candidates = append(candidates, key, option)
	}
	return candidates
}

func invalidCandidates(constraint config.IConstraint) []any {
	values := []any{nil, "", strings.Repeat("A", 300), float64(-1), "___not_an_option___", []any{true}}
	switch constraint.GetType() {
	case action.Bool:
		values = append(values, "nope")
	case action.Int, action.Float:
		values = append(values, "nope")
	case action.String:
		values = append(values, 1)
	case action.StringArray, action.IntArray, action.FloatArray:
		values = append(values, "nope")
	}
	return values
}
