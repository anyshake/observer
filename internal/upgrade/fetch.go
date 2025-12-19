package upgrade

import (
	"context"
	"time"

	"github.com/anyshake/observer/pkg/semver"
	"golang.org/x/sync/errgroup"
)

func (h *Helper) FetchRelease(ver *semver.Version, timeout time.Duration) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	toolchain := h.currentBuild.GetToolchainId()

	digestUrl, err := h.buildReleaseUrl(ver, toolchain, true)
	if err != nil {
		return nil, "", err
	}
	archiveUrl, err := h.buildReleaseUrl(ver, toolchain, false)
	if err != nil {
		return nil, "", err
	}

	var (
		digestData  []byte
		archiveData []byte
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		data, err := h.fetchDataFromUrl(ctx, digestUrl)
		if err != nil {
			return err
		}
		digestData = data
		return nil
	})

	g.Go(func() error {
		data, err := h.fetchDataFromUrl(ctx, archiveUrl)
		if err != nil {
			return err
		}
		archiveData = data
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, archiveUrl, err
	}

	digestMap := make(map[string]any)
	if err := h.unmarshallKvPair(string(digestData), "\n", &digestMap); err != nil {
		return nil, archiveUrl, err
	}

	releaseData, err := h.extractExecutableFromZip(archiveData, RELEASE_EXECUTABLE_NAME)
	if err != nil {
		return nil, archiveUrl, err
	}

	if err := h.verifyChecksum(releaseData, digestMap); err != nil {
		return nil, archiveUrl, err
	}

	return releaseData, archiveUrl, nil
}
