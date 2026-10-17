package upgrade

import (
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/unibuild"
)

func NewHelper(currentExePath string, currentVer *semver.Version, currentBuild *unibuild.UniBuild) *Helper {
	const cacheTimeout = time.Hour
	resolvers, err := dnsquery.NewResolvers()
	return &Helper{
		versionCheckDomain: VERSION_CHECK_DOMAIN,
		releaseFetchUrl:    RELEASE_FETCH_URL_TEMPLATE,
		resolvers:          resolvers,
		resolverErr:        err,
		currentExePath:     currentExePath,
		currentVer:         currentVer,
		currentBuild:       currentBuild,
		latestVer:          cache.NewGeneric[*semver.Version](cacheTimeout),
		requiredVer:        cache.NewGeneric[*semver.Version](cacheTimeout),
	}
}
