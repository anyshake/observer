package graph_resolver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/anyshake/observer/internal/dao/action"
)

func runLimitedQuery(t *testing.T, query string) (int32, string) {
	t.Helper()
	server := NewHandler(&Resolver{})
	var calls atomic.Int32
	server.AroundFields(func(context.Context, graphql.Resolver) (any, error) {
		calls.Add(1)
		return nil, errors.New("resolver reached")
	})
	body, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": map[string]any{"startTime": 0, "endTime": 3600000},
	})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)
	return calls.Load(), recorder.Body.String()
}

func TestGraphQLRejectsHistoryAliasesBeforeResolvers(t *testing.T) {
	for _, end := range []int{0, 3600000} {
		query := fmt.Sprintf(`query {
			first: getSeisRecordsByTime(startTime: 0, endTime: %d) { timestamp }
			...History
		} fragment History on Query {
			second: getSeisRecordsByTime(startTime: 0, endTime: %d) { timestamp }
		}`, end, end)
		calls, response := runLimitedQuery(t, query)
		if calls != 0 || !strings.Contains(response, "COMPLEXITY_LIMIT_EXCEEDED") {
			t.Fatalf("history alias budget: resolver calls = %d, response = %s", calls, response)
		}
	}
}

func TestGraphQLBoundsAliasesAndFragmentExpansion(t *testing.T) {
	var aliases strings.Builder
	aliases.WriteString("query {")
	for i := range maxOperationFields + 1 {
		fmt.Fprintf(&aliases, " f%d: getSoftwareVersion", i)
	}
	aliases.WriteString("}")
	var fragments strings.Builder
	fragments.WriteString("query { ...F0 } ")
	for i := range 10 {
		fmt.Fprintf(&fragments, "fragment F%d on Query { ...F%d ...F%d } ", i, i+1, i+1)
	}
	fragments.WriteString("fragment F10 on Query { getSoftwareVersion }")
	for _, query := range []string{aliases.String(), fragments.String()} {
		calls, response := runLimitedQuery(t, query)
		if calls != 0 || !strings.Contains(response, "operation exceeds") {
			t.Fatalf("operation limits: resolver calls = %d, response = %s", calls, response)
		}
	}
}

func TestGraphQLRejectsDuplicateWaveformOutput(t *testing.T) {
	query := `query {
		getSeisRecordsByTime(startTime: 0, endTime: 0) {
			channelData { first: data second: data }
		}
	}`
	calls, response := runLimitedQuery(t, query)
	if calls != 0 || !strings.Contains(response, "COMPLEXITY_LIMIT_EXCEEDED") {
		t.Fatalf("waveform output budget: resolver calls = %d, response = %s", calls, response)
	}
}

func TestGraphQLAdmitsDashboardQueries(t *testing.T) {
	for _, name := range []string{
		"Home/getHomeData.graphql",
		"Settings/about/getUpgradeStatus.graphql",
		"History/getSeismicRecords.graphql",
		"Settings/service/getServiceData.graphql",
		"Download/getFileListData.graphql",
	} {
		t.Run(name, func(t *testing.T) {
			query, err := os.ReadFile("../../../../web/src/src/views/" + name)
			if err != nil {
				t.Fatal(err)
			}
			calls, response := runLimitedQuery(t, string(query))
			if calls == 0 || strings.Contains(response, "COMPLEXITY_LIMIT_EXCEEDED") || strings.Contains(response, "operation exceeds") {
				t.Fatalf("normal query rejected: resolver calls = %d, response = %s", calls, response)
			}
		})
	}
}

func TestHistoryQueriesReleaseCapacityOnFailure(t *testing.T) {
	resolver := &Resolver{ActionHandler: &action.Handler{}}
	query := &queryResolver{resolver}
	for range maxConcurrentHistory + 1 {
		// A closed database fails after admission and must release the slot.
		if _, err := query.GetSeisRecordsByTime(t.Context(), 0, 0); err == nil || !strings.Contains(err.Error(), "database is not opened") {
			t.Fatalf("failed history request leaked capacity: %v", err)
		}
	}
	var releases []func()
	for range maxConcurrentHistory {
		release, err := resolver.historyQueries.acquire(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		releases = append(releases, release)
	}
	defer func() {
		for _, release := range releases {
			release()
		}
	}()
	if _, err := query.GetSeisRecordsByTime(t.Context(), 0, 0); err == nil || !strings.Contains(err.Error(), "concurrent history") {
		t.Fatalf("history concurrency limit did not reject before database access: %v", err)
	}
}
