package upgrade

import (
	"sync"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/unibuild"
)

type Helper struct {
	mu                 sync.Mutex
	versionCheckDomain string
	releaseFetchUrl    string
	resolvers          dnsquery.Resolvers

	currentBuild *unibuild.UniBuild
	currentVer   *semver.Version

	latestVer   cache.AnyCache
	requiredVer cache.AnyCache

	appliedVer *semver.Version
}
