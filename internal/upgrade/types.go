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
	resolverErr        error
	dnsClient          *dnsquery.Client

	currentExePath string
	currentBuild   *unibuild.UniBuild
	currentVer     *semver.Version

	latestVer   cache.GenericCache[*semver.Version]
	requiredVer cache.GenericCache[*semver.Version]

	appliedVer *semver.Version
}
