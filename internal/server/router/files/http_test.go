package files

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/service"
	"github.com/gin-gonic/gin"
)

func TestFileDownloadRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	assets := &assetService{data: map[string]*service.AssetData{
		"trace.txt": {FileName: "trace.txt", ContentType: "text/plain", Data: []byte("wave")},
	}}
	router := gin.New()
	if err := Setup(router.Group("/"), map[string]service.IService{"archiver": assets}, func(ctx *gin.Context) {
		ctx.Next()
	}); err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/files", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("missing path status = %d", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodPost, "/files", strings.NewReader(`{"file_path":"trace.txt"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("token status = %d body = %s", recorder.Code, recorder.Body.String())
	}
	var issued struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &issued); err != nil || issued.Data == "" {
		t.Fatalf("token body = %s, %v", recorder.Body.String(), err)
	}

	download := func(query url.Values) *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/files?"+query.Encode(), nil))
		return recorder
	}
	if download(url.Values{}).Code != http.StatusBadRequest {
		t.Fatal("missing query accepted")
	}
	if download(url.Values{"namespace": {"archiver"}, "file_path": {"trace.txt"}, "token": {"1:bad"}}).Code != http.StatusUnauthorized {
		t.Fatal("bad token accepted")
	}
	if download(url.Values{"namespace": {"missing"}, "file_path": {"trace.txt"}, "token": {issued.Data}}).Code != http.StatusBadRequest {
		t.Fatal("missing service accepted")
	}
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodPost, "/files", strings.NewReader(`{"file_path":"missing.txt"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	var missingToken struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &missingToken); err != nil || missingToken.Data == "" {
		t.Fatalf("missing token body = %s, %v", recorder.Body.String(), err)
	}
	if download(url.Values{"namespace": {"archiver"}, "file_path": {"missing.txt"}, "token": {missingToken.Data}}).Code != http.StatusBadRequest {
		t.Fatal("missing asset was not rejected")
	}
	ok := download(url.Values{"namespace": {"archiver"}, "file_path": {"trace.txt"}, "token": {issued.Data}})
	if ok.Code != http.StatusOK || ok.Body.String() != "wave" {
		t.Fatalf("download = %d %q", ok.Code, ok.Body.String())
	}
}

type assetService struct {
	data map[string]*service.AssetData
}

func (s *assetService) GetStatus() *service.Status                { return &service.Status{} }
func (s *assetService) GetName() string                           { return "assets" }
func (s *assetService) GetDescription() string                    { return "assets" }
func (s *assetService) Init() error                               { return nil }
func (s *assetService) IsEnabled() bool                           { return true }
func (s *assetService) Start() error                              { return nil }
func (s *assetService) Stop() error                               { return nil }
func (s *assetService) Restart() error                            { return nil }
func (s *assetService) GetAssetList() ([]service.Asset, error)    { return nil, nil }
func (s *assetService) GetConfigConstraint() []config.IConstraint { return nil }
func (s *assetService) GetAssetData(id string) (*service.AssetData, error) {
	asset, ok := s.data[id]
	if !ok {
		return nil, errors.New("asset not found")
	}
	return asset, nil
}
