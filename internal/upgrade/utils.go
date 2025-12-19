package upgrade

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/anyshake/observer/pkg/semver"
	"github.com/samber/lo"
)

func (h *Helper) unmarshallKvPair(metadata, sep string, v *map[string]any) error {
	if v == nil {
		return errors.New("nil map provided")
	}
	if *v == nil {
		*v = make(map[string]any)
	}

	pairs := strings.Split(metadata, sep)
	if len(pairs) == 0 {
		return nil
	}

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := kv[1]
		(*v)[key] = value
	}

	return nil
}

func (h *Helper) buildReleaseUrl(ver *semver.Version, toolchain string, digest bool) (string, error) {
	tpl, err := template.New("url").Parse(h.releaseFetchUrl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, map[string]string{
		"Version":       ver.String(),
		"ToolchainName": toolchain,
		"Extension":     lo.Ternary(digest, "dgst", "zip"),
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (h *Helper) fetchDataFromUrl(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch url %s: status %s", url, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (h *Helper) verifyChecksum(data []byte, expect map[string]any) error {
	algos := map[string]func([]byte) []byte{
		"MD5": func(b []byte) []byte {
			h := md5.Sum(b)
			return h[:]
		},
		"SHA1": func(b []byte) []byte {
			h := sha1.Sum(b)
			return h[:]
		},
		"SHA2-256": func(b []byte) []byte {
			h := sha256.Sum256(b)
			return h[:]
		},
		"SHA2-512": func(b []byte) []byte {
			h := sha512.Sum512(b)
			return h[:]
		},
	}

	for name, algo := range algos {
		exp, ok := expect[name]
		if !ok {
			continue
		}
		actual := hex.EncodeToString(algo(data))
		expVal, ok := exp.(string)
		if !ok {
			return fmt.Errorf("invalid expected %s checksum format", name)
		}
		if !strings.EqualFold(actual, expVal) {
			return fmt.Errorf("integrity check failed: at %s checksum, expect %s got %s", name, expVal, actual)
		}
	}

	return nil
}

func (h *Helper) extractExecutableFromZip(archive []byte, fileName string) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		name := path.Base(f.Name)
		if name != fileName {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		return io.ReadAll(rc)
	}

	return nil, fmt.Errorf("file %s not found in zip", fileName)
}

func (u *Helper) isEligibleForUpdate(latestVer, requiredVer *semver.Version) bool {
	if !u.currentVer.IsCompatible(latestVer) {
		return false
	}
	if u.currentVer.GreaterThanOrEqual(latestVer) || u.currentVer.LessThan(requiredVer) {
		return false
	}

	if u.appliedVer != nil {
		if u.appliedVer.GreaterThanOrEqual(latestVer) {
			return false
		}
	}

	return true
}
