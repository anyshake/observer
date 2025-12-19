package upgrade

import (
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/unibuild"
)

func NewHelper(currentVer *semver.Version, currentBuild *unibuild.UniBuild) *Helper {
	const cacheTimeout = time.Hour
	return &Helper{
		versionCheckDomain: VERSION_CHECK_DOMAIN,
		releaseFetchUrl:    RELEASE_FETCH_URL_TEMPLATE,
		resolvers:          dnsquery.NewResolvers(),
		currentVer:         currentVer,
		currentBuild:       currentBuild,
		latestVer:          cache.New(cacheTimeout),
		requiredVer:        cache.New(cacheTimeout),
	}
}
