package graph_resolver

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	MaxRequestBytes       = 1 << 20
	RequestTimeout        = 30 * time.Second
	maxOperationCost      = 1000
	maxOperationFields    = 256
	maxOperationDepth     = 20
	maxConcurrentHistory  = 2
	maxHistoryRecords     = 72000
	maxHistorySamples     = 8_000_000
	maxHistoryStoredBytes = 64 << 20
)

func NewHandler(resolver *Resolver) *handler.Server {
	var costs ComplexityRoot
	costs.Query.GetSeisRecordsByTime = historyQueryCost
	// A waveform array can contain millions of samples; duplicating it with
	// nested aliases must also consume the operation's output budget.
	costs._channelData.Data = func(int) int { return 300 }
	listCost := func(childComplexity int) int {
		return min(maxOperationCost+1, 100+10*min(childComplexity, maxOperationCost+1))
	}
	costs.Query.GetEventsBySource = func(childComplexity int, _ string) int { return listCost(childComplexity) }
	costs.Query.GetMiniSeedFiles = listCost
	costs.Query.GetHelicorderFiles = listCost
	costs.Query.GetApplicationLogs = listCost
	costs.Query.GetSysUsers = listCost
	costs.Query.GetServiceConfigConstraint = listCost
	costs.Query.GetStationConfigConstraint = listCost

	server := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver, Complexity: costs}))
	server.SetParserTokenLimit(10000)
	// Bound fragment expansion and aliases before the complexity walker runs.
	server.Use(operationLimits{})
	server.Use(extension.FixedComplexityLimit(maxOperationCost))
	return server
}

func historyQueryCost(childComplexity int, startTime, endTime int64) int {
	window := time.UnixMilli(endTime).Sub(time.UnixMilli(startTime))
	if window < 0 || window > time.Hour {
		return maxOperationCost + 1
	}
	minutes := int((window + time.Minute - 1) / time.Minute)
	// Decoding happens even for timestamp-only selections. The base cost ensures
	// two history aliases cannot share one operation, including zero-length windows.
	return min(maxOperationCost+1, 501+2*minutes+min(childComplexity, maxOperationCost+1))
}

type operationLimits struct{}

func (operationLimits) ExtensionName() string                   { return "OperationLimits" }
func (operationLimits) Validate(graphql.ExecutableSchema) error { return nil }

func (operationLimits) MutateOperationContext(ctx context.Context, operation *graphql.OperationContext) *gqlerror.Error {
	remaining := maxOperationFields
	var visit func(ast.SelectionSet, int) bool
	visit = func(selections ast.SelectionSet, depth int) bool {
		if depth > maxOperationDepth || ctx.Err() != nil {
			return false
		}
		for _, selection := range selections {
			remaining--
			if remaining < 0 {
				return false
			}
			switch field := selection.(type) {
			case *ast.Field:
				if len(field.SelectionSet) > 0 && !visit(field.SelectionSet, depth+1) {
					return false
				}
			case *ast.FragmentSpread:
				if !visit(field.Definition.SelectionSet, depth+1) {
					return false
				}
			case *ast.InlineFragment:
				if !visit(field.SelectionSet, depth+1) {
					return false
				}
			}
		}
		return true
	}
	if !visit(operation.Operation.SelectionSet, 1) {
		return gqlerror.Errorf("operation exceeds field, depth or execution limits")
	}
	return nil
}

type historyQueryLimit struct {
	once  sync.Once
	slots chan struct{}
}

func (limit *historyQueryLimit) acquire(ctx context.Context) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	limit.once.Do(func() { limit.slots = make(chan struct{}, maxConcurrentHistory) })
	select {
	case limit.slots <- struct{}{}:
		return func() { <-limit.slots }, nil
	default:
		return nil, errors.New("too many concurrent history queries; try again later")
	}
}
