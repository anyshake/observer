package export

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/metadata"
	"github.com/gin-gonic/gin"
)

func TestSetupListsFormatsAndRejectsUnauthorizedCalls(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, handler := testsupport.OpenDAO(t)
	engine := gin.New()
	Setup(engine.Group("/api"), handler, hardwareWithChannels("EHZ", "EHE"), func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	})

	denied := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/export", nil)
	engine.ServeHTTP(denied, req)
	if denied.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized status = %d", denied.Code)
	}

	engine = gin.New()
	Setup(engine.Group("/api"), handler, hardwareWithChannels("EHZ", "EHE"), func(ctx *gin.Context) { ctx.Next() })
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/export", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d body = %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Message string `json:"message"`
		Data    struct {
			DataFormat  map[string]string `json:"data_format"`
			ChannelCode []string          `json:"channel_code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Message != "data formats" || len(body.Data.DataFormat) != 6 {
		t.Fatalf("formats = %#v", body)
	}
	for _, format := range []string{"sac", "txt", "wav", "mseed_int32", "mseed_steim1", "mseed_steim2"} {
		if body.Data.DataFormat[format] == "" {
			t.Fatalf("missing format %s in %#v", format, body.Data.DataFormat)
		}
	}
	if len(body.Data.ChannelCode) != 2 || body.Data.ChannelCode[0] != "EHZ" {
		t.Fatalf("channels = %#v", body.Data.ChannelCode)
	}
}

func TestSetupPostEncodesStoredRecords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, handler := testsupport.OpenDAO(t)
	initExportCodes(t, handler)
	start := time.Date(2024, 3, 15, 13, 5, 9, 123000000, time.UTC)
	samples := make([]int32, 16)
	for i := range samples {
		samples[i] = int32((i - 8) * 20)
	}
	insertRecords(t, db, syntheticRecords(start, 2, 16, []explorer.ChannelData{{
		ChannelCode: "EHZ", ChannelId: 1, ByteSize: 4, DataType: "int32", Data: samples,
	}}))

	engine := gin.New()
	Setup(engine.Group("/api"), handler, hardwareWithChannels("EHZ"), func(ctx *gin.Context) { ctx.Next() })

	for _, format := range []string{"txt", "sac", "wav", "mseed_int32", "mseed_steim1", "mseed_steim2"} {
		rec := postExport(engine, map[string]any{
			"start_time":   start.UnixMilli(),
			"end_time":     start.Add(time.Second).UnixMilli(),
			"channel_code": "EHZ",
			"data_format":  format,
		})
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d body = %s", format, rec.Code, rec.Body.String())
		}
		if rec.Body.Len() == 0 {
			t.Fatalf("%s response was empty", format)
		}
		disposition := rec.Header().Get("Content-Disposition")
		if format == "txt" && disposition != "attachment; filename=2024.075.13.05.09.0123.SHAKE.AS.00.EHZ.D.txt" {
			t.Fatalf("txt disposition = %s", disposition)
		}
		if format == "wav" && !bytes.HasPrefix(rec.Body.Bytes(), []byte("RIFF")) {
			t.Fatalf("wav prefix = %q", rec.Body.Bytes()[:4])
		}
	}
}

func TestSetupPostRejectsBadRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, handler := testsupport.OpenDAO(t)
	initExportCodes(t, handler)
	start := time.Date(2024, 3, 15, 13, 5, 9, 0, time.UTC)
	insertRecords(t, db, syntheticRecords(start, 2, 8, []explorer.ChannelData{{
		ChannelCode: "EHE", ChannelId: 2, ByteSize: 4, DataType: "int32", Data: []int32{1, 2, 3, 4},
	}}))
	jittered := syntheticRecords(start.Add(time.Minute), 2, 8, []explorer.ChannelData{{
		ChannelCode: "EHZ", Data: []int32{1, 2, 3, 4},
	}})
	jittered[1].RecordTime = jittered[0].RecordTime + 500
	insertRecords(t, db, jittered)

	engine := gin.New()
	Setup(engine.Group("/api"), handler, hardwareWithChannels("EHZ", "EHE"), func(ctx *gin.Context) { ctx.Next() })

	cases := []struct {
		name   string
		body   any
		status int
		text   string
	}{
		{name: "invalid body", body: "{", status: http.StatusBadRequest, text: "not valid"},
		{name: "missing fields", body: map[string]any{"start_time": 1}, status: http.StatusBadRequest, text: "not valid"},
		{name: "unknown channel", body: exportBody(start, "EHN", "txt"), status: http.StatusBadRequest, text: "was not found"},
		{name: "unknown format", body: exportBody(start, "EHZ", "csv"), status: http.StatusBadRequest, text: "unknown data format"},
		{name: "empty range", body: exportBody(start.Add(time.Hour), "EHZ", "txt"), status: http.StatusNotFound, text: "no seis records"},
		{name: "reversed range", body: map[string]any{
			"start_time": start.Add(time.Second).UnixMilli(), "end_time": start.UnixMilli(),
			"channel_code": "EHZ", "data_format": "txt",
		}, status: http.StatusInternalServerError, text: "failed to encode"},
		{name: "jitter", body: exportBody(start.Add(time.Minute), "EHZ", "wav"), status: http.StatusInternalServerError, text: "jitter"},
		{name: "channel absent from records", body: exportBody(start, "EHZ", "mseed_int32"), status: http.StatusBadRequest, text: "unknown data type"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := postExport(engine, tc.body)
			if rec.Code != tc.status || !bytes.Contains(rec.Body.Bytes(), []byte(tc.text)) {
				t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestSetupPostReportsFilenameErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, handler := testsupport.OpenDAO(t)
	start := time.Date(2024, 3, 15, 13, 5, 9, 0, time.UTC)
	insertRecords(t, db, syntheticRecords(start, 2, 8, []explorer.ChannelData{{
		ChannelCode: "EHZ", Data: []int32{3, -3, 6, -6},
	}}))
	engine := gin.New()
	Setup(engine.Group("/api"), handler, hardwareWithChannels("EHZ"), func(ctx *gin.Context) { ctx.Next() })

	rec := postExport(engine, exportBody(start, "EHZ", "wav"))
	if rec.Code != http.StatusInternalServerError || !bytes.Contains(rec.Body.Bytes(), []byte("failed to get file name")) {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}
}

func initExportCodes(t *testing.T, handler *action.Handler) {
	t.Helper()
	constraints := []config.IConstraint{
		&config.StationStationCodeConfigConstraintImpl{},
		&config.StationLocationCodeConfigConstraintImpl{},
		&config.StationNetworkCodeConfigConstraintImpl{},
	}
	for _, constraint := range constraints {
		if err := constraint.Init(handler); err != nil {
			t.Fatal(err)
		}
	}
}

func insertRecords(t *testing.T, db *dao.DAO, records []model.SeisRecord) {
	t.Helper()
	seen := map[string]struct{}{}
	for _, record := range records {
		day := time.UnixMilli(record.RecordTime).UTC().YearDay()
		table := fmt.Sprintf("%sseis_records_%d", db.GetPrefix(), day%model.SEIS_RECORD_SHARDS)
		if _, ok := seen[table]; !ok {
			if err := db.Database.Table(table).AutoMigrate(&model.SeisRecord{}); err != nil {
				t.Fatal(err)
			}
			seen[table] = struct{}{}
		}
	}
	if err := action.NewHandler(db).SeisRecordsCreate(records...); err != nil {
		t.Fatal(err)
	}
}

func exportBody(start time.Time, channel, format string) map[string]any {
	return map[string]any{
		"start_time":   start.UnixMilli(),
		"end_time":     start.Add(time.Second).UnixMilli(),
		"channel_code": channel,
		"data_format":  format,
	}
}

func postExport(engine *gin.Engine, body any) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	switch typed := body.(type) {
	case string:
		reader = bytes.NewReader([]byte(typed))
	default:
		payload, err := json.Marshal(typed)
		if err != nil {
			panic(err)
		}
		reader = bytes.NewReader(payload)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/export", reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

type exportHardware struct {
	channels []string
}

func hardwareWithChannels(channels ...string) hardware.IHardware {
	return exportHardware{channels: channels}
}

func (exportHardware) Open(context.Context) (context.Context, context.CancelFunc, error) {
	return context.Background(), func() {}, nil
}

func (h exportHardware) GetConfig() explorer.DeviceConfig {
	cfg := explorer.DeviceConfig{}
	cfg.SetChannelCodes(append([]string(nil), h.channels...))
	cfg.SetSampleRate(100)
	cfg.SetProtocol("v1")
	cfg.SetModel("E-C111G")
	return cfg
}

func (exportHardware) GetStatus() explorer.DeviceStatus { return explorer.DeviceStatus{} }
func (exportHardware) Close() error                     { return nil }
func (exportHardware) Flush() error                     { return nil }
func (exportHardware) Subscribe(string, explorer.EventHandler) error {
	return nil
}
func (exportHardware) Unsubscribe(string) error { return nil }
func (exportHardware) SubscribeRealtime(string, explorer.EventHandler) error {
	return nil
}
func (exportHardware) UnsubscribeRealtime(string) error { return nil }
func (exportHardware) GetCoordinates(bool) (float64, float64, float64, error) {
	return 0, 0, 0, nil
}
func (exportHardware) GetTemperature() (float64, error) { return 0, nil }
func (exportHardware) GetDeviceId() string              { return "device" }
func (exportHardware) GetMetadata(string, string, string, string, string, string, string, bool) (*metadata.Render, error) {
	return nil, nil
}
