package upgrade

import (
	"context"
	"time"

	"github.com/anyshake/observer/pkg/semver"
	"golang.org/x/sync/errgroup"
)

func (h *Helper) FetchRelease(ctx context.Context, ver *semver.Version, timeout time.Duration) ([]byte, string, error) {
	return h.FetchReleaseWithProgress(ctx, ver, timeout, nil)
}

func (h *Helper) FetchReleaseWithProgress(ctx context.Context, ver *semver.Version, timeout time.Duration, progress func(read, total int64)) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
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
		data, err := h.fetchDataFromUrlWithProgress(ctx, archiveUrl, progress)
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
