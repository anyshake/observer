package metadata

import "testing"

func TestSensitivitiesFromModelTemplates(t *testing.T) {
	t.Parallel()
	type want struct {
		gain  float64
		unit  string
		poles int
	}
	cases := map[string][]want{
		"E-C111G": {
			{85899345880, UNIT_VELOCITY, 3}, {85899345880, UNIT_VELOCITY, 3}, {85899345880, UNIT_VELOCITY, 3},
			{1670.7030433430377, UNIT_ACCELERATION, 0}, {1670.7030433430377, UNIT_ACCELERATION, 0}, {1670.7030433430377, UNIT_ACCELERATION, 0},
		},
		"E-C131G": {
			{85899345880, UNIT_VELOCITY, 3}, {85899345880, UNIT_VELOCITY, 3}, {85899345880, UNIT_VELOCITY, 3},
			{26104.735052234966, UNIT_ACCELERATION, 0}, {26104.735052234966, UNIT_ACCELERATION, 0}, {26104.735052234966, UNIT_ACCELERATION, 0},
		},
		"E-C121G": {
			{85899345880, UNIT_VELOCITY, 2}, {85899345880, UNIT_VELOCITY, 2}, {85899345880, UNIT_VELOCITY, 2},
			{1670.7030433430377, UNIT_ACCELERATION, 0}, {1670.7030433430377, UNIT_ACCELERATION, 0}, {1670.7030433430377, UNIT_ACCELERATION, 0},
		},
		"E-D001": {
			{343597383520, UNIT_VELOCITY, 3}, {343597383520, UNIT_VELOCITY, 3}, {343597383520, UNIT_VELOCITY, 3},
		},
	}
	for model, expected := range cases {
		got, err := Sensitivities(model)
		if err != nil {
			t.Fatalf("%s: %v", model, err)
		}
		if len(got) != len(expected) {
			t.Fatalf("%s: got %d channels, want %d", model, len(got), len(expected))
		}
		for i, want := range expected {
			if got[i] == nil || got[i].Gain != want.gain || got[i].Unit != want.unit || len(got[i].Poles) != want.poles {
				t.Errorf("%s channel %d = %+v, want %+v", model, i+1, got[i], want)
			}
		}
	}
	if _, err := Sensitivities("missing-model"); err == nil {
		t.Fatal("missing model returned sensitivities")
	}
}

func TestSensitivityResponsePoles(t *testing.T) {
	t.Parallel()
	got, err := Sensitivities("E-C121G")
	if err != nil {
		t.Fatal(err)
	}
	velocity := got[0]
	if velocity.Frequency != 4.5 {
		t.Fatalf("velocity sensitivity frequency = %g, want 4.5", velocity.Frequency)
	}
	if velocity.Poles[0] != complex(-2.2211060060879837, 2.221776881419294) {
		t.Fatalf("first pole = %v", velocity.Poles[0])
	}
	if got[3].Frequency != 0 {
		t.Fatalf("accelerometer sensitivity frequency = %g, want 0 for None", got[3].Frequency)
	}
}

func TestNormalizeUnit(t *testing.T) {
	t.Parallel()
	for name, want := range map[string]string{"m/s": UNIT_VELOCITY, "M/S**2": UNIT_ACCELERATION, "m/s2": UNIT_ACCELERATION, "COUNTS": ""} {
		if got := NormalizeUnit(name); got != want {
			t.Errorf("NormalizeUnit(%q) = %q, want %q", name, got, want)
		}
	}
}
