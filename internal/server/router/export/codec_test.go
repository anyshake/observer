package export

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/bclswl0827/mseedio"
)

func TestForEachContinuousChannelRecord(t *testing.T) {
	start := time.Date(2024, 3, 15, 13, 5, 9, 0, time.UTC)
	good := syntheticRecords(start, 2, 100, []explorer.ChannelData{{
		ChannelCode: "EHZ", ChannelId: 1, ByteSize: 4, DataType: "int32", Data: []int32{1, -2, 3},
	}})

	var seen int
	got, err := forEachContinuousChannelRecord(iterOf(good), "EHZ", func(_ model.SeisRecord, samples []int32) error {
		seen += len(samples)
		return nil
	})
	if err != nil || got.Count != 2 || got.SampleRate != 100 || seen != 6 {
		t.Fatalf("range = %+v, seen = %d, err = %v", got, seen, err)
	}

	skipped, err := forEachContinuousChannelRecord(iterOf(good), "EHE", func(model.SeisRecord, []int32) error {
		t.Fatal("missing channel invoked the callback")
		return nil
	})
	if err != nil || skipped.Count != 2 {
		t.Fatalf("skipped range = %+v, err = %v", skipped, err)
	}

	empty, err := forEachContinuousChannelRecord(iterOf(nil), "EHZ", func(model.SeisRecord, []int32) error {
		t.Fatal("empty iterator invoked the callback")
		return nil
	})
	if err != nil || empty.Count != 0 {
		t.Fatalf("empty range = %+v, err = %v", empty, err)
	}

	if _, err := forEachContinuousChannelRecord(func(func(model.SeisRecord) error) error {
		return errors.New("iterator failed")
	}, "EHZ", nil); err == nil || !strings.Contains(err.Error(), "iterator failed") {
		t.Fatalf("iterator error = %v", err)
	}

	broken := good[0]
	broken.ChannelData = []byte{1}
	if _, err := forEachContinuousChannelRecord(iterOf([]model.SeisRecord{broken}), "EHZ", nil); err == nil {
		t.Fatal("corrupt channel data was accepted")
	}

	if _, err := forEachContinuousChannelRecord(iterOf(good), "EHZ", func(model.SeisRecord, []int32) error {
		return errors.New("stop")
	}); err == nil || err.Error() != "stop" {
		t.Fatalf("callback error = %v", err)
	}

	jitter := syntheticRecords(start, 2, 100, mustChannels(t, good[0]))
	jitter[1].RecordTime += 500
	if _, err := forEachContinuousChannelRecord(iterOf(jitter), "EHZ", func(model.SeisRecord, []int32) error { return nil }); err == nil || !strings.Contains(err.Error(), "jitter") {
		t.Fatalf("jitter error = %v", err)
	}

	rates := syntheticRecords(start, 2, 100, mustChannels(t, good[0]))
	rates[1].SampleRate = 50
	if _, err := forEachContinuousChannelRecord(iterOf(rates), "EHZ", func(model.SeisRecord, []int32) error { return nil }); err == nil || !strings.Contains(err.Error(), "sample rate") {
		t.Fatalf("sample rate error = %v", err)
	}
}

func TestTxtEncoder(t *testing.T) {
	start := time.Date(2024, 3, 15, 13, 5, 9, 123000000, time.UTC)
	channels := []explorer.ChannelData{{
		ChannelCode: "EHZ", ChannelId: 1, ByteSize: 4, DataType: "int32", Data: []int32{10, -4},
	}}
	records := syntheticRecords(start, 2, 2, channels)
	encoder := &seismicDataEncoderTxtImpl{}

	if encoder.GetName() != "TXT" {
		t.Fatalf("name = %s", encoder.GetName())
	}
	if data, err := encoder.Encode(iterOf(nil), "EHZ"); err != nil || data != nil {
		t.Fatalf("empty encode = %q, %v", data, err)
	}

	data, err := encoder.Encode(iterOf(records), "EHZ")
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, sample := range []string{" 10\n", " -4\n"} {
		if !strings.Contains(text, sample) {
			t.Fatalf("txt %q missing %q", text, sample)
		}
	}
	if _, err := encoder.Encode(iterOf([]model.SeisRecord{{ChannelData: []byte{9}}}), "EHZ"); err == nil {
		t.Fatal("corrupt txt record accepted")
	}

	handler := openSettings(t, "station", "location", "network")
	named := txtEncoder(handler)
	filename, err := named.GetFileName(start, "EHZ")
	if err != nil {
		t.Fatal(err)
	}
	if filename != "2024.075.13.05.09.0123.SHAKE.AS.00.EHZ.D.txt" {
		t.Fatalf("filename = %s", filename)
	}

	if _, err := txtEncoder(openSettings(t)).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing station code accepted")
	}
	if _, err := txtEncoder(openSettings(t, "station")).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing location code accepted")
	}
	if _, err := txtEncoder(openSettings(t, "station", "location")).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing network code accepted")
	}
}

func TestSacEncoder(t *testing.T) {
	start := time.Date(2024, 3, 15, 13, 5, 9, 123000000, time.UTC)
	channels := []explorer.ChannelData{{
		ChannelCode: "EHZ", ChannelId: 2, ByteSize: 4, DataType: "int32", Data: []int32{5, 6, 7, 8},
	}}
	records := syntheticRecords(start, 2, 4, channels)
	encoder := sacEncoder(openSettings(t, "station", "location", "network"))
	if encoder.GetName() != "SAC" {
		t.Fatalf("name = %s", encoder.GetName())
	}
	if data, err := encoder.Encode(iterOf(nil), "EHZ"); err != nil || data != nil {
		t.Fatalf("empty encode = %d bytes, %v", len(data), err)
	}
	data, err := encoder.Encode(iterOf(records), "EHZ")
	if err != nil || len(data) == 0 {
		t.Fatalf("sac encode = %d bytes, %v", len(data), err)
	}
	if data, err := encoder.Encode(iterOf(records), "EHE"); err != nil || len(data) == 0 {
		t.Fatalf("missing channel encode = %d bytes, %v", len(data), err)
	}
	if _, err := encoder.Encode(func(func(model.SeisRecord) error) error { return errors.New("read failed") }, "EHZ"); err == nil {
		t.Fatal("iterator error ignored")
	}

	filename, err := encoder.GetFileName(start, "EHZ")
	if err != nil {
		t.Fatal(err)
	}
	if filename != "2024.075.13.05.09.0123.SHAKE.AS.SHAKE.EHZ.D.sac" {
		t.Fatalf("filename = %s", filename)
	}
	if _, err := sacEncoder(openSettings(t)).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing station code accepted")
	}
	if _, err := sacEncoder(openSettings(t, "station")).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing network code accepted")
	}
	if _, err := sacEncoder(openSettings(t)).Encode(iterOf(records), "EHZ"); err == nil {
		t.Fatal("encode without station code accepted")
	}
	if _, err := sacEncoder(openSettings(t, "station")).Encode(iterOf(records), "EHZ"); err == nil {
		t.Fatal("encode without location code accepted")
	}
	if _, err := sacEncoder(openSettings(t, "station", "location")).Encode(iterOf(records), "EHZ"); err == nil {
		t.Fatal("encode without network code accepted")
	}
}

func TestMseedEncoder(t *testing.T) {
	start := time.Date(2024, 3, 15, 13, 5, 9, 123000000, time.UTC)
	samples := make([]int32, 16)
	for i := range samples {
		samples[i] = int32(i*3 - 7)
	}
	channels := []explorer.ChannelData{{
		ChannelCode: "EHZ", ChannelId: 1, ByteSize: 4, DataType: "int32", Data: samples,
	}}
	records := syntheticRecords(start, 2, 100, channels)
	handler := openSettings(t, "station", "location", "network")

	for _, tc := range []struct {
		name   string
		encode int
	}{
		{"MiniSEED (INT32)", mseedio.INT32},
		{"MiniSEED (STEIM1)", mseedio.STEIM1},
		{"MiniSEED (STEIM2)", mseedio.STEIM2},
	} {
		encoder := mseedEncoder(handler, tc.name, tc.encode)
		if encoder.GetName() != tc.name {
			t.Fatalf("name = %s", encoder.GetName())
		}
		if data, err := encoder.Encode(iterOf(nil), "EHZ"); err != nil || data != nil {
			t.Fatalf("%s empty encode = %d bytes, %v", tc.name, len(data), err)
		}
		data, err := encoder.Encode(iterOf(records), "EHZ")
		if err != nil || len(data) == 0 {
			t.Fatalf("%s encode = %d bytes, %v", tc.name, len(data), err)
		}
		if data, err := encoder.Encode(iterOf(records), "EHE"); err != nil || len(data) != 0 {
			t.Fatalf("%s missing channel = %d bytes, %v", tc.name, len(data), err)
		}
		filename, err := encoder.GetFileName(start, "EHZ")
		if err != nil {
			t.Fatal(err)
		}
		if filename != "2024.075.13.05.09.0123.SHAKE.AS.00.EHZ.D.mseed" {
			t.Fatalf("filename = %s", filename)
		}
	}

	if _, err := mseedEncoder(handler, "bad", 999).Encode(iterOf(records), "EHZ"); err == nil {
		t.Fatal("invalid encoding accepted")
	}
	if _, err := mseedEncoder(handler, "int32", mseedio.INT32).Encode(func(func(model.SeisRecord) error) error {
		return errors.New("read failed")
	}, "EHZ"); err == nil {
		t.Fatal("iterator error ignored")
	}
	broken := records[0]
	broken.ChannelData = []byte{2}
	if _, err := mseedEncoder(handler, "int32", mseedio.INT32).Encode(iterOf([]model.SeisRecord{broken}), "EHZ"); err == nil {
		t.Fatal("corrupt mseed record accepted")
	}

	if _, err := mseedEncoder(openSettings(t), "int32", mseedio.INT32).Encode(iterOf(records), "EHZ"); err == nil {
		t.Fatal("missing station code accepted")
	}
	if _, err := mseedEncoder(openSettings(t, "station"), "int32", mseedio.INT32).Encode(iterOf(records), "EHZ"); err == nil {
		t.Fatal("missing location code accepted")
	}
	if _, err := mseedEncoder(openSettings(t, "station", "location"), "int32", mseedio.INT32).Encode(iterOf(records), "EHZ"); err == nil {
		t.Fatal("missing network code accepted")
	}
	if _, err := mseedEncoder(openSettings(t, "station"), "int32", mseedio.INT32).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing location code accepted")
	}
	if _, err := mseedEncoder(openSettings(t, "station", "location"), "int32", mseedio.INT32).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing network code accepted")
	}
	if _, err := mseedEncoder(openSettings(t), "int32", mseedio.INT32).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("filename without station code accepted")
	}
}

func TestWavEncoder(t *testing.T) {
	start := time.Date(2024, 3, 15, 13, 5, 9, 123000000, time.UTC)
	samples := []int32{1000, -1000, 500, -250, 125, -60, 30, -15}
	channels := []explorer.ChannelData{{
		ChannelCode: "EHZ", ChannelId: 1, ByteSize: 4, DataType: "int32", Data: samples,
	}}
	encoder := &seismicDataEncoderWavImpl{outputSampleRate: 8000}
	if encoder.GetName() != "WAV (Audio)" {
		t.Fatalf("name = %s", encoder.GetName())
	}
	if data, err := encoder.Encode(iterOf(nil), "EHZ"); err != nil || data != nil {
		t.Fatalf("empty encode = %d bytes, %v", len(data), err)
	}

	one := syntheticRecords(start, 1, 8, channels)
	if _, err := encoder.Encode(iterOf(one), "EHZ"); err == nil || !strings.Contains(err.Error(), "time difference") {
		t.Fatalf("single record error = %v", err)
	}

	records := syntheticRecords(start, 2, 8, channels)
	data, err := encoder.Encode(iterOf(records), "EHZ")
	if err != nil || len(data) < 12 || string(data[:4]) != "RIFF" {
		t.Fatalf("wav encode = %d bytes, %v", len(data), err)
	}

	zeros := syntheticRecords(start, 2, 4, []explorer.ChannelData{{
		ChannelCode: "EHZ", Data: []int32{0, 0, 0, 0},
	}})
	if data, err := encoder.Encode(iterOf(zeros), "EHZ"); err != nil || len(data) < 12 {
		t.Fatalf("silent wav = %d bytes, %v", len(data), err)
	}
	if _, err := encoder.Encode(iterOf(records), "EHE"); err != nil {
		t.Fatalf("missing channel error = %v", err)
	}
	if _, err := encoder.Encode(func(func(model.SeisRecord) error) error { return errors.New("read failed") }, "EHZ"); err == nil {
		t.Fatal("iterator error ignored")
	}

	same := encoder.linearInterpolate([]int16{1, 2, 3}, 100, 100)
	if len(same) != 3 || same[2] != 3 {
		t.Fatalf("same rate = %v", same)
	}
	unchanged := encoder.linearInterpolate([]int16{4, 5}, 0, 100)
	if len(unchanged) != 2 || unchanged[0] != 4 {
		t.Fatalf("zero rate = %v", unchanged)
	}
	stretched := encoder.linearInterpolate([]int16{0, 1000}, 1, 4)
	if len(stretched) != 8 || stretched[0] != 0 || stretched[len(stretched)-1] == 0 {
		t.Fatalf("stretched = %v", stretched)
	}

	normalized := encoder.normalizeToInt16([]int32{-2, 4})
	if len(normalized) != 2 || normalized[1] <= 0 || normalized[0] >= 0 {
		t.Fatalf("normalized = %v", normalized)
	}
	if encoder.normalizeToInt16([]int32{0, 0}) != nil {
		t.Fatal("zero samples were normalized")
	}
	kernel := encoder.getLowPassFilter(200, 10, 5)
	if len(kernel) != 5 {
		t.Fatalf("kernel = %v", kernel)
	}
	filtered := encoder.applyFilter([]int16{1, 2, 3}, kernel)
	if len(filtered) != 3 {
		t.Fatalf("filtered = %v", filtered)
	}

	handler := openSettings(t, "station", "location", "network")
	named := wavEncoder(handler, 8000)
	filename, err := named.GetFileName(start, "EHZ")
	if err != nil {
		t.Fatal(err)
	}
	if filename != "2024.075.13.05.09.0123.SHAKE.AS.00.EHZ.D.wav" {
		t.Fatalf("filename = %s", filename)
	}
	if _, err := wavEncoder(openSettings(t), 8000).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing station code accepted")
	}
	if _, err := wavEncoder(openSettings(t, "station"), 8000).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing location code accepted")
	}
	if _, err := wavEncoder(openSettings(t, "station", "location"), 8000).GetFileName(start, "EHZ"); err == nil {
		t.Fatal("missing network code accepted")
	}
}

func syntheticRecords(start time.Time, count, rate int, channels []explorer.ChannelData) []model.SeisRecord {
	records := make([]model.SeisRecord, count)
	for i := 0; i < count; i++ {
		cloned := make([]explorer.ChannelData, len(channels))
		for n := range channels {
			cloned[n] = channels[n]
			cloned[n].Data = append([]int32(nil), channels[n].Data...)
		}
		if err := records[i].Encode(start.Add(time.Duration(i)*time.Second), rate, cloned); err != nil {
			panic(err)
		}
	}
	return records
}

func iterOf(records []model.SeisRecord) seismicRecordIterator {
	return func(callback func(model.SeisRecord) error) error {
		for _, record := range records {
			if err := callback(record); err != nil {
				return err
			}
		}
		return nil
	}
}

func openSettings(t *testing.T, codes ...string) *action.Handler {
	t.Helper()
	_, handler := testsupport.OpenDAO(t)
	known := map[string]config.IConstraint{
		"station":  &config.StationStationCodeConfigConstraintImpl{},
		"location": &config.StationLocationCodeConfigConstraintImpl{},
		"network":  &config.StationNetworkCodeConfigConstraintImpl{},
	}
	for _, code := range codes {
		if err := known[code].Init(handler); err != nil {
			t.Fatal(err)
		}
	}
	return handler
}

func txtEncoder(handler *action.Handler) *seismicDataEncoderTxtImpl {
	return &seismicDataEncoderTxtImpl{
		actionHandler:      handler,
		stationCodeConfig:  config.StationStationCodeConfigConstraintImpl{},
		locationCodeConfig: config.StationLocationCodeConfigConstraintImpl{},
		networkCodeConfig:  config.StationNetworkCodeConfigConstraintImpl{},
	}
}

func sacEncoder(handler *action.Handler) *seismicDataEncoderSacImpl {
	return &seismicDataEncoderSacImpl{
		actionHandler:      handler,
		stationCodeConfig:  config.StationStationCodeConfigConstraintImpl{},
		locationCodeConfig: config.StationLocationCodeConfigConstraintImpl{},
		networkCodeConfig:  config.StationNetworkCodeConfigConstraintImpl{},
	}
}

func mseedEncoder(handler *action.Handler, name string, encodeType int) *seismicDataEncoderMseedImpl {
	return &seismicDataEncoderMseedImpl{
		name:               name,
		encodeType:         encodeType,
		actionHandler:      handler,
		stationCodeConfig:  config.StationStationCodeConfigConstraintImpl{},
		locationCodeConfig: config.StationLocationCodeConfigConstraintImpl{},
		networkCodeConfig:  config.StationNetworkCodeConfigConstraintImpl{},
	}
}

func wavEncoder(handler *action.Handler, rate int) *seismicDataEncoderWavImpl {
	return &seismicDataEncoderWavImpl{
		outputSampleRate:   rate,
		actionHandler:      handler,
		stationCodeConfig:  config.StationStationCodeConfigConstraintImpl{},
		locationCodeConfig: config.StationLocationCodeConfigConstraintImpl{},
		networkCodeConfig:  config.StationNetworkCodeConfigConstraintImpl{},
	}
}

func mustChannels(t *testing.T, record model.SeisRecord) []explorer.ChannelData {
	t.Helper()
	_, _, channels, err := record.Decode()
	if err != nil {
		t.Fatal(err)
	}
	return channels
}
