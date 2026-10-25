package graph_resolver

import (
	"context"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/server/middleware/auth_jwt"
	graph_model "github.com/anyshake/observer/internal/server/router/graph/model"
	"github.com/anyshake/observer/internal/service"
	service_helicorder "github.com/anyshake/observer/internal/service/helicorder"
	service_miniseed "github.com/anyshake/observer/internal/service/miniseed"
	"github.com/anyshake/observer/internal/upgrade"
	"github.com/anyshake/observer/pkg/jobtracker"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/seisevent"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/unibuild"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ContextKey string

type Resolver struct {
	AppContext               context.Context
	serviceControl           serviceControl
	historyQueries           historyQueryLimit
	RestartChan              chan struct{}
	ExecutablePath           string
	CurrentVersion           *semver.Version
	CurrentBuild             *unibuild.UniBuild
	UpgradeHelper            *upgrade.Helper
	HardwareDev              hardware.IHardware
	TimeSource               *timesource.Source
	ActionHandler            *action.Handler
	LogBuffer                *ringbuf.Buffer[string]
	StationConfigConstraints []config.IConstraint
	ServiceMap               map[string]service.IService
	SeisEventSource          map[string]seisevent.IDataSource
	DataPurgeJob             *jobtracker.Tracker
}

func (r *Resolver) getCurrentUserId(ctx context.Context) string {
	user := ctx.Value(ContextKey("user_status"))
	if user == nil {
		return ""
	}

	userStatusMap, ok := user.(map[string]any)
	if !ok {
		return ""
	}

	user, ok = userStatusMap[auth_jwt.UserIdKey]
	if !ok {
		return ""
	}

	return user.(string)
}

func (r *Resolver) checkIsAdmin(ctx context.Context) bool {
	user := ctx.Value(ContextKey("user_status"))
	if user == nil {
		return false
	}

	userStatusMap, ok := user.(map[string]any)
	if !ok {
		return false
	}

	isAdmin, ok := userStatusMap[auth_jwt.IsAdminKey]
	if !ok {
		return false
	}

	return isAdmin.(bool)
}

func (r *Resolver) getMiniSeedService() (*service_miniseed.MiniSeedServiceImpl, error) {
	serviceObj, ok := r.ServiceMap[service_miniseed.ID]
	if !ok {
		return nil, fmt.Errorf("service was not found, maybe it was excluded from building")
	}

	typedServiceObj, ok := serviceObj.(*service_miniseed.MiniSeedServiceImpl)
	if !ok {
		return nil, fmt.Errorf("service %s has unexpected type", service_miniseed.ID)
	}

	return typedServiceObj, nil
}

func (r *Resolver) getHelicorderService() (*service_helicorder.HelicorderServiceImpl, error) {
	serviceObj, ok := r.ServiceMap[service_helicorder.ID]
	if !ok {
		return nil, fmt.Errorf("service was not found, maybe it was excluded from building")
	}

	typedServiceObj, ok := serviceObj.(*service_helicorder.HelicorderServiceImpl)
	if !ok {
		return nil, fmt.Errorf("service %s has unexpected type", service_helicorder.ID)
	}

	return typedServiceObj, nil
}

func (r *Resolver) toFDSNAccessKey(record model.FDSNAccessKey) *graph_model.FdsnAccessKey {
	optionalTime := func(value int64) *int64 {
		if value == 0 {
			return nil
		}
		return &value
	}
	return &graph_model.FdsnAccessKey{
		KeyID: record.KeyId, KeyPrefix: record.KeyPrefix, Remark: record.Remark,
		CreatedAt: record.CreatedAt, ExpiresAt: optionalTime(record.ExpiresAt),
		RevokedAt: optionalTime(record.RevokedAt), LastRequestedAt: optionalTime(record.LastRequestedAt),
		UserIP: record.UserIp, UserAgent: record.UserAgent,
	}
}

func (r *Resolver) toPurgeDataJobResponse(jobStatus *jobtracker.Job) *graph_model.PurgeDataJob {
	var (
		jobStatusEnum      graph_model.JobStatus
		startedAtInt64Ptr  *int64
		finishedAtInt64Ptr *int64
		errorStringPtr     *string
	)

	switch jobStatus.Status {
	case jobtracker.JobStatusIdle:
		jobStatusEnum = graph_model.JobStatusIdle
	case jobtracker.JobStatusRunning:
		jobStatusEnum = graph_model.JobStatusRunning
	case jobtracker.JobStatusSucceeded:
		jobStatusEnum = graph_model.JobStatusSucceeded
	case jobtracker.JobStatusFailed:
		jobStatusEnum = graph_model.JobStatusFailed
	default:
		jobStatusEnum = graph_model.JobStatusIdle
	}

	if jobStatus.StartedAt != nil {
		t := jobStatus.StartedAt.UnixMilli()
		startedAtInt64Ptr = &t
	}
	if jobStatus.FinishedAt != nil {
		t := jobStatus.FinishedAt.UnixMilli()
		finishedAtInt64Ptr = &t
	}
	if jobStatus.Error != nil {
		err := jobStatus.Error.Error()
		errorStringPtr = &err
	}

	return &graph_model.PurgeDataJob{
		ID:         jobStatus.ID,
		Kind:       jobStatus.Kind,
		Status:     jobStatusEnum,
		StartedAt:  startedAtInt64Ptr,
		FinishedAt: finishedAtInt64Ptr,
		Error:      errorStringPtr,
	}
}

func LoadOrCreatePurgeDataJob(resolver *Resolver) *jobtracker.Tracker {
	if resolver == nil {
		return nil
	}
	if resolver.DataPurgeJob == nil {
		ctx := resolver.AppContext
		if ctx == nil {
			ctx = context.Background()
		}
		resolver.DataPurgeJob = jobtracker.New(ctx, "data_purge_job")
	}

	return resolver.DataPurgeJob
}
