package upgrade

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/miekg/dns"
)

func (u *Helper) CheckUpdate() (latest, required *semver.Version, eligible bool, applied bool, err error) {
	if u.latestVer.Valid() && u.requiredVer.Valid() {
		latest, ok := u.latestVer.Get().(*semver.Version)
		if !ok {
			return nil, nil, false, false, errors.New("cached latest version has invalid type")
		}
		required, ok := u.requiredVer.Get().(*semver.Version)
		if !ok {
			return nil, nil, false, false, errors.New("cached required version has invalid type")
		}
		if u.appliedVer != nil {
			return latest, required, u.isEligibleForUpdate(latest, required), latest.Equal(u.appliedVer), nil
		}
		return latest, required, u.isEligibleForUpdate(latest, required), false, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		txt []string
		err error
	}
	resultCh := make(chan result, len(u.resolvers))

	for _, r := range u.resolvers {
		resolver := r
		go func() {
			dq, err := dnsquery.New(resolver.Server)
			if err != nil {
				resultCh <- result{nil, fmt.Errorf("resolver init failed: %w", err)}
				return
			}
			if err := dq.Open(); err != nil {
				resultCh <- result{nil, fmt.Errorf("resolver open failed: %w", err)}
				return
			}
			defer dq.Close()

			msg := (&dns.Msg{}).SetQuestion(fmt.Sprintf("%s.", u.versionCheckDomain), dns.TypeTXT)
			res, err := dq.Query(msg, 5*time.Second)
			if err != nil {
				resultCh <- result{nil, fmt.Errorf("query failed: %w", err)}
				return
			}
			if len(res.Answer) == 0 {
				resultCh <- result{nil, errors.New("empty answer")}
				return
			}

			for _, ans := range res.Answer {
				if txt, ok := ans.(*dns.TXT); ok && len(txt.Txt) > 0 {
					resultCh <- result{txt.Txt, nil}
					return
				}
			}

			resultCh <- result{nil, errors.New("no valid TXT record")}
		}()
	}

	for range u.resolvers {
		select {
		case r := <-resultCh:
			if r.err == nil && len(r.txt) > 0 {
				var metadataMap map[string]any
				if err := u.unmarshallKvPair(strings.Join(r.txt, ""), ";", &metadataMap); err != nil {
					return nil, nil, false, false, fmt.Errorf("metadata unmarshal failed: %w", err)
				}

				var (
					latestMajor   string
					latestMinor   string
					latestPatch   string
					requiredMajor string
					requiredMinor string
					requiredPatch string
				)
				if val, ok := metadataMap["latest_major"].(string); !ok {
					return nil, nil, false, false, errors.New("latest_major missing or invalid")
				} else {
					latestMajor = val
				}
				if val, ok := metadataMap["latest_minor"].(string); !ok {
					return nil, nil, false, false, errors.New("latest_minor missing or invalid")
				} else {
					latestMinor = val
				}
				if val, ok := metadataMap["latest_patch"].(string); !ok {
					return nil, nil, false, false, errors.New("latest_patch missing or invalid")
				} else {
					latestPatch = val
				}
				if val, ok := metadataMap["required_major"].(string); !ok {
					return nil, nil, false, false, errors.New("required_major missing or invalid")
				} else {
					requiredMajor = val
				}
				if val, ok := metadataMap["required_minor"].(string); !ok {
					return nil, nil, false, false, errors.New("required_minor missing or invalid")
				} else {
					requiredMinor = val
				}
				if val, ok := metadataMap["required_patch"].(string); !ok {
					return nil, nil, false, false, errors.New("required_patch missing or invalid")
				} else {
					requiredPatch = val
				}

				latest := semver.New(latestMajor, latestMinor, latestPatch, "")
				u.latestVer.Set(latest)

				required := semver.New(requiredMajor, requiredMinor, requiredPatch, "")
				u.requiredVer.Set(required)

				if u.appliedVer != nil {
					return latest, required, u.isEligibleForUpdate(latest, required), latest.Equal(u.appliedVer), nil
				}
				return latest, required, u.isEligibleForUpdate(latest, required), false, nil
			}
		case <-ctx.Done():
			return nil, nil, false, false, ctx.Err()
		}
	}

	return nil, nil, false, false, errors.New("all resolvers failed")
}
