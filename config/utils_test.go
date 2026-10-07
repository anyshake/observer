package config

import (
	"encoding/json"
	"testing"
)

func TestConfigValueParsers(t *testing.T) {
	t.Parallel()

	intCases := []any{json.Number("42"), "42", int64(42), float64(42)}
	for _, value := range intCases {
		got, err := GetConfigValInt64(value)
		if err != nil || got != 42 {
			t.Fatalf("GetConfigValInt64(%#v) = %d, %v", value, got, err)
		}
	}
	if _, err := GetConfigValInt64(true); err == nil {
		t.Fatal("bool accepted as int")
	}
	if _, err := GetConfigValInt64("nope"); err == nil {
		t.Fatal("invalid integer string accepted")
	}

	gotInts, err := GetConfigValInt64Array([]any{json.Number("1"), "2", int64(3), float64(4)})
	if err != nil || len(gotInts) != 4 || gotInts[0] != 1 || gotInts[3] != 4 {
		t.Fatalf("GetConfigValInt64Array() = %v, %v", gotInts, err)
	}
	emptyInts, err := GetConfigValInt64Array(nil)
	if err != nil || len(emptyInts) != 0 {
		t.Fatalf("nil int array = %v, %v", emptyInts, err)
	}
	if _, err := GetConfigValInt64Array("nope"); err == nil {
		t.Fatal("string accepted as int array")
	}
	if _, err := GetConfigValInt64Array([]any{"bad"}); err == nil {
		t.Fatal("invalid int array element accepted")
	}

	floatCases := []any{json.Number("1.5"), "1.5", float64(1.5)}
	for _, value := range floatCases {
		got, err := GetConfigValFloat64(value)
		if err != nil || got != 1.5 {
			t.Fatalf("GetConfigValFloat64(%#v) = %v, %v", value, got, err)
		}
	}
	if _, err := GetConfigValFloat64(true); err == nil {
		t.Fatal("bool accepted as float")
	}

	gotFloats, err := GetConfigValFloat64Array([]any{json.Number("1.5"), "2.5", float64(3.5)})
	if err != nil || len(gotFloats) != 3 || gotFloats[2] != 3.5 {
		t.Fatalf("GetConfigValFloat64Array() = %v, %v", gotFloats, err)
	}
	if _, err := GetConfigValFloat64Array(nil); err != nil {
		t.Fatal(err)
	}
	if _, err := GetConfigValFloat64Array(1); err == nil {
		t.Fatal("int accepted as float array")
	}
	if _, err := GetConfigValFloat64Array([]any{true}); err == nil {
		t.Fatal("invalid float array element accepted")
	}

	gotString, err := GetConfigValString("station")
	if err != nil || gotString != "station" {
		t.Fatalf("GetConfigValString() = %q, %v", gotString, err)
	}
	if _, err := GetConfigValString(1); err == nil {
		t.Fatal("int accepted as string")
	}

	gotStrings, err := GetConfigValStringArray([]any{"EHZ", "EHN"})
	if err != nil || len(gotStrings) != 2 || gotStrings[1] != "EHN" {
		t.Fatalf("GetConfigValStringArray() = %v, %v", gotStrings, err)
	}
	if _, err := GetConfigValStringArray(nil); err != nil {
		t.Fatal(err)
	}
	if _, err := GetConfigValStringArray("EHZ"); err == nil {
		t.Fatal("string accepted as string array")
	}
	if _, err := GetConfigValStringArray([]any{1}); err == nil {
		t.Fatal("invalid string array element accepted")
	}

	gotBool, err := GetConfigValBool(true)
	if err != nil || !gotBool {
		t.Fatalf("GetConfigValBool() = %v, %v", gotBool, err)
	}
	if _, err := GetConfigValBool("true"); err == nil {
		t.Fatal("string accepted as bool")
	}
}
