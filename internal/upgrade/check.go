package upgrade

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/miekg/dns"
)

func (u *Helper) CheckUpdate() (latest, required *semver.Version, eligible bool, applied bool, err error) {
	if u.latestVer.Valid() && u.requiredVer.Valid() {
		latest, required := u.latestVer.Get(), u.requiredVer.Get()
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

				keys := []string{
					"latest_major", "latest_minor", "latest_patch",
					"required_major", "required_minor", "required_patch",
				}
				var parts [6]string
				for i, key := range keys {
					value, ok := metadataMap[key].(string)
					if !ok {
						return nil, nil, false, false, fmt.Errorf("%s missing or invalid", key)
					}
					number, err := strconv.ParseInt(value, 10, 64)
					if err != nil || number < 0 {
						return nil, nil, false, false, fmt.Errorf("%s missing or invalid", key)
					}
					parts[i] = value
				}

				latest := semver.New(parts[0], parts[1], parts[2], "")
				if latest.GetMajor() == 0 && latest.GetMinor() == 0 && latest.GetPatch() == 0 {
					return nil, nil, false, false, errors.New("latest version is zero")
				}
				required := semver.New(parts[3], parts[4], parts[5], "")
				if latest.LessThan(required) {
					return nil, nil, false, false, errors.New("required version exceeds latest version")
				}
				u.latestVer.Set(latest)
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
