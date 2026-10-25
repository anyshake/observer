package upgrade

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/miekg/dns"
)

const (
	versionCheckTimeout = 20 * time.Second
	resolverTimeout     = 2 * time.Second
)

func (u *Helper) CheckUpdate(ctx context.Context) (latest, required *semver.Version, eligible bool, applied bool, err error) {
	return u.checkUpdate(ctx, false)
}

func (u *Helper) CheckUpdateFresh(ctx context.Context) (latest, required *semver.Version, eligible bool, applied bool, err error) {
	return u.checkUpdate(ctx, true)
}

func (u *Helper) checkUpdate(ctx context.Context, fresh bool) (latest, required *semver.Version, eligible bool, applied bool, err error) {
	if !fresh && u.latestVer.Valid() && u.requiredVer.Valid() {
		latest, required := u.latestVer.Get(), u.requiredVer.Get()
		eligible, applied := u.updateStatus(latest, required)
		return latest, required, eligible, applied, nil
	}

	// Leave time for later resolver batches when earlier endpoints are blocked.
	ctx, cancel := context.WithTimeout(ctx, versionCheckTimeout)
	defer cancel()

	client, err := u.getDNSClient()
	if err != nil {
		return nil, nil, false, false, fmt.Errorf("all resolvers failed: %w", err)
	}
	message := new(dns.Msg).SetQuestion(dns.Fqdn(u.GetVersionCheckDomain()), dns.TypeTXT)
	response, err := client.QueryContext(ctx, message)
	if err != nil {
		return nil, nil, false, false, err
	}
	var txt []string
	for _, answer := range response.Answer {
		if record, ok := answer.(*dns.TXT); ok && len(record.Txt) > 0 {
			txt = record.Txt
			break
		}
	}
	var metadataMap map[string]any
	if err := u.unmarshallKvPair(strings.Join(txt, ""), ";", &metadataMap); err != nil {
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

	latest = semver.New(parts[0], parts[1], parts[2], "")
	if latest.GetMajor() == 0 && latest.GetMinor() == 0 && latest.GetPatch() == 0 {
		return nil, nil, false, false, errors.New("latest version is zero")
	}
	required = semver.New(parts[3], parts[4], parts[5], "")
	if latest.LessThan(required) {
		return nil, nil, false, false, errors.New("required version exceeds latest version")
	}
	u.latestVer.Set(latest)
	u.requiredVer.Set(required)

	eligible, applied = u.updateStatus(latest, required)
	return latest, required, eligible, applied, nil
}

func (u *Helper) getDNSClient() (*dnsquery.Client, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.closed {
		return nil, dnsquery.ErrClosed
	}
	if u.resolverErr != nil {
		return nil, u.resolverErr
	}
	if u.dnsClient == nil {
		client, err := dnsquery.NewClient(u.resolvers, len(u.resolvers), resolverTimeout,
			dnsquery.WithLogger(logger.GetLogger("upgrade")))
		if err != nil {
			return nil, err
		}
		u.dnsClient = client
	}
	return u.dnsClient, nil
}
