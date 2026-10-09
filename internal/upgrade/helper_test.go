package upgrade

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/unibuild"
	"github.com/miekg/dns"
)

func TestNewHelperAndOptions(t *testing.T) {
	exe := filepath.Join(t.TempDir(), "observer")
	version := semver.New("1", "2", "3", "")
	build := unibuild.New("test-toolchain", "stable", "abc", "100")
	helper := NewHelper(exe, version, build)
	if helper.GetVersionCheckDomain() != VERSION_CHECK_DOMAIN || helper.GetReleaseFetchUrl() != RELEASE_FETCH_URL_TEMPLATE {
		t.Fatalf("defaults = %s %s", helper.GetVersionCheckDomain(), helper.GetReleaseFetchUrl())
	}
	if helper.currentExePath != exe || helper.currentVer != version || helper.currentBuild != build {
		t.Fatal("constructor did not keep the current build identity")
	}
	if len(helper.resolvers) == 0 {
		t.Fatal("resolver catalog is empty")
	}

	helper.SetVersionCheckDomain("updates.example.test")
	helper.SetReleaseFetchUrl("https://example.test/{{.Version}}.{{.Extension}}")
	if helper.GetVersionCheckDomain() != "updates.example.test" || helper.GetReleaseFetchUrl() != "https://example.test/{{.Version}}.{{.Extension}}" {
		t.Fatal("setters did not update the helper")
	}

	done := make(chan struct{})
	go func() {
		for i := 0; i < 20; i++ {
			helper.SetVersionCheckDomain("a")
			_ = helper.GetReleaseFetchUrl()
		}
		close(done)
	}()
	for i := 0; i < 20; i++ {
		helper.SetReleaseFetchUrl("b")
		_ = helper.GetVersionCheckDomain()
	}
	<-done
}

func TestCheckUpdateUsesCachedMetadata(t *testing.T) {
	current := semver.New("1", "2", "0", "")
	latest := semver.New("1", "4", "0", "")
	required := semver.New("1", "1", "0", "")
	helper := testHelper(t, current)
	helper.latestVer.Set(latest)
	helper.requiredVer.Set(required)

	gotLatest, gotRequired, eligible, applied, err := helper.CheckUpdate()
	if err != nil || !gotLatest.Equal(latest) || !gotRequired.Equal(required) || !eligible || applied {
		t.Fatalf("cached update = %v %v %v %v %v", gotLatest, gotRequired, eligible, applied, err)
	}

	helper.appliedVer = latest
	_, _, eligible, applied, err = helper.CheckUpdate()
	if err != nil || eligible || !applied {
		t.Fatalf("applied latest = eligible %v applied %v err %v", eligible, applied, err)
	}

	helper.currentVer = semver.New("1", "4", "0", "")
	helper.appliedVer = nil
	_, _, eligible, applied, err = helper.CheckUpdate()
	if err != nil || eligible || applied {
		t.Fatalf("current is latest = eligible %v applied %v err %v", eligible, applied, err)
	}

	helper.currentVer = semver.New("1", "0", "0", "")
	helper.appliedVer = semver.New("1", "2", "0", "")
	_, _, eligible, _, err = helper.CheckUpdate()
	if err != nil || eligible {
		t.Fatalf("below required minimum was eligible, err = %v", err)
	}

	helper.currentVer = semver.New("1", "2", "0", "rc.1")
	_, _, eligible, _, err = helper.CheckUpdate()
	if err != nil || eligible {
		t.Fatalf("prerelease was eligible, err = %v", err)
	}

	helper.currentVer = semver.New("0", "0", "0", "")
	_, _, eligible, _, err = helper.CheckUpdate()
	if err != nil || eligible {
		t.Fatalf("custom version was eligible, err = %v", err)
	}

	helper.currentVer = current
	helper.appliedVer = semver.New("1", "4", "1", "")
	_, _, eligible, applied, err = helper.CheckUpdate()
	if err != nil || eligible || applied {
		t.Fatalf("newer applied version = eligible %v applied %v err %v", eligible, applied, err)
	}
}

func TestCheckUpdateFromLocalDNS(t *testing.T) {
	helper := testHelper(t, semver.New("1", "2", "0", ""))
	txt := "latest_major=1;latest_minor=4;latest_patch=0;required_major=1;required_minor=1;required_patch=0"
	address := serveDNS(t, dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		reply := new(dns.Msg)
		reply.SetReply(r)
		if len(r.Question) > 0 {
			name := r.Question[0].Name
			reply.Answer = []dns.RR{
				&dns.A{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 30}, A: net.ParseIP("127.0.0.1").To4()},
				&dns.TXT{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 30}, Txt: []string{txt}},
			}
		}
		_ = w.WriteMsg(reply)
	}))
	helper.resolvers = dnsquery.Resolvers{{Name: "bad", Server: "tcp://nope"}, {Name: "local", Server: address}}

	latest, required, eligible, applied, err := helper.CheckUpdate()
	if err != nil || latest.String() != "v1.4.0" || required.String() != "v1.1.0" || !eligible || applied {
		t.Fatalf("dns update = %v %v %v %v %v", latest, required, eligible, applied, err)
	}

	appliedHelper := helperWithTXT(t, txt)
	appliedHelper.appliedVer = semver.New("1", "4", "0", "")
	latest, required, eligible, applied, err = appliedHelper.CheckUpdate()
	if err != nil || latest.String() != "v1.4.0" || required.String() != "v1.1.0" || eligible || !applied {
		t.Fatalf("applied dns update = %v %v %v %v %v", latest, required, eligible, applied, err)
	}

	helper.resolvers = nil
	latest, required, eligible, applied, err = helper.CheckUpdate()
	if err != nil || latest.String() != "v1.4.0" || !eligible || applied {
		t.Fatalf("cached after dns = %v %v %v %v %v", latest, required, eligible, applied, err)
	}
}

func TestCheckUpdateRejectsBadMetadata(t *testing.T) {
	fields := []string{"latest_major", "latest_minor", "latest_patch", "required_major", "required_minor", "required_patch"}
	base := map[string]string{
		"latest_major": "1", "latest_minor": "4", "latest_patch": "0",
		"required_major": "1", "required_minor": "1", "required_patch": "0",
	}
	for _, missing := range fields {
		t.Run(missing, func(t *testing.T) {
			pairs := make([]string, 0, len(base)-1)
			for key, value := range base {
				if key == missing {
					continue
				}
				pairs = append(pairs, key+"="+value)
			}
			helper := helperWithTXT(t, strings.Join(pairs, ";"))
			if _, _, _, _, err := helper.CheckUpdate(); err == nil || !strings.Contains(err.Error(), missing) {
				t.Fatalf("error = %v", err)
			}
		})
	}

	t.Run("empty txt", func(t *testing.T) {
		helper := helperWithTXT(t, "")
		if _, _, _, _, err := helper.CheckUpdate(); err == nil {
			t.Fatal("empty txt accepted")
		}
	})
	t.Run("empty answer", func(t *testing.T) {
		helper := testHelper(t, semver.New("1", "2", "0", ""))
		helper.resolvers = dnsquery.Resolvers{{Server: serveDNS(t, dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
			reply := new(dns.Msg)
			reply.SetReply(r)
			_ = w.WriteMsg(reply)
		}))}}
		if _, _, _, _, err := helper.CheckUpdate(); err == nil || !strings.Contains(err.Error(), "all resolvers failed") {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("no txt", func(t *testing.T) {
		helper := testHelper(t, semver.New("1", "2", "0", ""))
		helper.resolvers = dnsquery.Resolvers{{Server: serveDNS(t, dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
			reply := new(dns.Msg)
			reply.SetReply(r)
			if len(r.Question) > 0 {
				reply.Answer = []dns.RR{&dns.A{
					Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1},
					A:   net.ParseIP("127.0.0.1").To4(),
				}}
			}
			_ = w.WriteMsg(reply)
		}))}}
		if _, _, _, _, err := helper.CheckUpdate(); err == nil || !strings.Contains(err.Error(), "all resolvers failed") {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("blank txt record", func(t *testing.T) {
		helper := testHelper(t, semver.New("1", "2", "0", ""))
		helper.resolvers = dnsquery.Resolvers{{Server: serveDNS(t, dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
			reply := new(dns.Msg)
			reply.SetReply(r)
			if len(r.Question) > 0 {
				reply.Answer = []dns.RR{&dns.TXT{
					Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 1},
					Txt: []string{},
				}}
			}
			_ = w.WriteMsg(reply)
		}))}}
		if _, _, _, _, err := helper.CheckUpdate(); err == nil || !strings.Contains(err.Error(), "all resolvers failed") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestCheckUpdateResolverFailures(t *testing.T) {
	helper := testHelper(t, semver.New("1", "2", "0", ""))
	helper.resolvers = dnsquery.Resolvers{{Server: "tcp://nope"}, {Server: ""}}
	if _, _, _, _, err := helper.CheckUpdate(); err == nil || !strings.Contains(err.Error(), "all resolvers failed") {
		t.Fatalf("init failures error = %v", err)
	}

	helper.resolvers = dnsquery.Resolvers{{Server: "sdns://not-a-stamp"}}
	if _, _, _, _, err := helper.CheckUpdate(); err == nil || !strings.Contains(err.Error(), "all resolvers failed") {
		t.Fatalf("open failure error = %v", err)
	}

	packet, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = packet.Close() })
	helper.resolvers = dnsquery.Resolvers{{Server: "udp://" + packet.LocalAddr().String()}}
	_, _, _, _, err = helper.CheckUpdate()
	if err == nil || (!errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "all resolvers failed")) {
		t.Fatalf("silent resolver error = %v", err)
	}
}

func TestReleaseURLAndFetchHelpers(t *testing.T) {
	helper := testHelper(t, semver.New("1", "2", "3", ""))
	helper.releaseFetchUrl = "https://downloads.example.test/{{.Version}}/{{.ToolchainName}}.{{.Extension}}"
	got, err := helper.buildReleaseUrl(helper.currentVer, "test-toolchain", true)
	if err != nil || got != "https://downloads.example.test/v1.2.3/test-toolchain.dgst" {
		t.Fatalf("digest url = %s, %v", got, err)
	}
	got, err = helper.buildReleaseUrl(helper.currentVer, "test-toolchain", false)
	if err != nil || !strings.HasSuffix(got, ".zip") {
		t.Fatalf("archive url = %s, %v", got, err)
	}
	helper.releaseFetchUrl = "{{"
	if _, err := helper.buildReleaseUrl(helper.currentVer, "test-toolchain", false); err == nil {
		t.Fatal("broken template accepted")
	}
	helper.releaseFetchUrl = "{{index .Version 99}}"
	if _, err := helper.buildReleaseUrl(helper.currentVer, "test-toolchain", false); err == nil {
		t.Fatal("template execution error ignored")
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			_, _ = w.Write([]byte("payload"))
		case "/hang":
			<-r.Context().Done()
		default:
			http.Error(w, "nope", http.StatusBadGateway)
		}
	}))
	t.Cleanup(srv.Close)

	body, err := helper.fetchDataFromUrl(context.Background(), srv.URL+"/ok")
	if err != nil || string(body) != "payload" {
		t.Fatalf("fetch = %q, %v", body, err)
	}
	if _, err := helper.fetchDataFromUrl(context.Background(), srv.URL+"/missing"); err == nil {
		t.Fatal("bad status accepted")
	}
	if _, err := helper.fetchDataFromUrl(context.Background(), "://bad"); err == nil {
		t.Fatal("invalid url accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := helper.fetchDataFromUrl(ctx, srv.URL+"/hang"); err == nil {
		t.Fatal("canceled fetch accepted")
	}
}

func TestFetchReleaseFromTestServer(t *testing.T) {
	payload := []byte("observer-binary-v1")
	archive := zipArchive(t, map[string][]byte{"dist/observer": payload, "README": []byte("notes")})
	digest := checksumText(payload, true)
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		switch {
		case strings.HasSuffix(r.URL.Path, ".dgst"):
			writeResponse(w, []byte(digest))
		case strings.HasSuffix(r.URL.Path, ".zip"):
			writeResponse(w, archive)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	helper := testHelper(t, semver.New("1", "2", "3", "rc.1"))
	helper.SetReleaseFetchUrl(srv.URL + "/{{.Version}}/{{.ToolchainName}}.{{.Extension}}")
	version := semver.New("9", "8", "7", "")
	data, archiveURL, err := helper.FetchRelease(version, 5*time.Second)
	if count := hits.Load(); err != nil || !bytes.Equal(data, payload) || count != 2 {
		t.Fatalf("fetch = %q hits %d url %s err %v", data, count, archiveURL, err)
	}
	if !strings.Contains(archiveURL, "/v9.8.7/test-toolchain.zip") {
		t.Fatalf("archive url = %s", archiveURL)
	}

	helper.SetReleaseFetchUrl("{{")
	if _, url, err := helper.FetchRelease(version, time.Second); err == nil || url != "" {
		t.Fatalf("bad template = %s, %v", url, err)
	}

	helper.SetReleaseFetchUrl(srv.URL + "/nope")
	if _, url, err := helper.FetchRelease(version, time.Second); err == nil || !strings.Contains(url, "/nope") {
		t.Fatalf("missing release = %s, %v", url, err)
	}
}

func TestFetchReleaseRejectsBadArtifacts(t *testing.T) {
	payload := []byte("observer-binary-v2")
	archive := zipArchive(t, map[string][]byte{"bin/observer": payload})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "bad-digest"):
			if strings.HasSuffix(r.URL.Path, ".dgst") {
				_, _ = w.Write([]byte("SHA2-256=deadbeef\n"))
			} else {
				writeResponse(w, archive)
			}
		case strings.Contains(r.URL.Path, "bad-zip"):
			if strings.HasSuffix(r.URL.Path, ".dgst") {
				_, _ = w.Write([]byte("NOTE=ignored\n"))
			} else {
				_, _ = w.Write([]byte("not a zip"))
			}
		case strings.Contains(r.URL.Path, "missing-file"):
			if strings.HasSuffix(r.URL.Path, ".dgst") {
				writeResponse(w, []byte(checksumText([]byte("other"), false)))
			} else {
				writeResponse(w, zipArchive(t, map[string][]byte{"README": []byte("only")}))
			}
		case strings.Contains(r.URL.Path, "status"):
			http.Error(w, "busy", http.StatusServiceUnavailable)
		default:
			<-r.Context().Done()
		}
	}))
	t.Cleanup(srv.Close)

	version := semver.New("1", "0", "0", "")
	for _, tc := range []struct {
		name    string
		pattern string
		timeout time.Duration
		text    string
	}{
		{name: "checksum", pattern: srv.URL + "/bad-digest/file.{{.Extension}}", timeout: time.Second, text: "integrity check"},
		{name: "zip", pattern: srv.URL + "/bad-zip/file.{{.Extension}}", timeout: time.Second, text: "zip"},
		{name: "missing", pattern: srv.URL + "/missing-file/file.{{.Extension}}", timeout: time.Second, text: "not found"},
		{name: "status", pattern: srv.URL + "/status/file.{{.Extension}}", timeout: time.Second, text: "status"},
		{name: "timeout", pattern: srv.URL + "/hang/file.{{.Extension}}", timeout: 20 * time.Millisecond, text: "context"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			helper := testHelper(t, version)
			helper.SetReleaseFetchUrl(tc.pattern)
			_, archiveURL, err := helper.FetchRelease(version, tc.timeout)
			if err == nil || !strings.Contains(strings.ToLower(err.Error()), tc.text) || archiveURL == "" {
				t.Fatalf("url %s error %v", archiveURL, err)
			}
		})
	}
}

func TestUtilityHelpers(t *testing.T) {
	helper := testHelper(t, semver.New("1", "0", "0", ""))
	if err := helper.unmarshallKvPair("a=b", ";", nil); err == nil {
		t.Fatal("nil map accepted")
	}
	var decoded map[string]any
	if err := helper.unmarshallKvPair("a=b;skip;c=d=e;;", ";", &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["a"] != "b" || decoded["c"] != "d=e" || len(decoded) != 2 {
		t.Fatalf("decoded = %#v", decoded)
	}
	if err := helper.unmarshallKvPair("", "\n", &decoded); err != nil || decoded["a"] != "b" {
		t.Fatalf("empty metadata cleared the map: %#v %v", decoded, err)
	}
	var created map[string]any
	if err := helper.unmarshallKvPair("k=v", ";", &created); err != nil || created["k"] != "v" {
		t.Fatalf("created = %#v, %v", created, err)
	}

	payload := []byte("checksum-me")
	if err := helper.verifyChecksum(payload, map[string]any{}); err != nil {
		t.Fatal(err)
	}
	if err := helper.verifyChecksum(payload, map[string]any{"MD5": strings.ToUpper(checksumPayloadMD5)}); err != nil {
		t.Fatal(err)
	}
	if err := helper.verifyChecksum(payload, map[string]any{"SHA1": "00", "SHA2-256": hex.EncodeToString(sha256Sum(payload))}); err == nil {
		t.Fatal("mismatched checksum accepted")
	}
	if err := helper.verifyChecksum(payload, map[string]any{"SHA2-512": 12}); err == nil {
		t.Fatal("non-string checksum accepted")
	}
	expect := map[string]any{
		"MD5":      checksumPayloadMD5,
		"SHA1":     checksumPayloadSHA1,
		"SHA2-256": hex.EncodeToString(sha256Sum(payload)),
		"SHA2-512": hex.EncodeToString(sha512Sum(payload)),
		"EXTRA":    "ignored",
	}
	if err := helper.verifyChecksum(payload, expect); err != nil {
		t.Fatal(err)
	}

	extracted, err := helper.extractExecutableFromZip(zipArchive(t, map[string][]byte{"nested/observer": payload}), RELEASE_EXECUTABLE_NAME)
	if err != nil || !bytes.Equal(extracted, payload) {
		t.Fatalf("extracted = %q, %v", extracted, err)
	}
	if _, err := helper.extractExecutableFromZip([]byte("nope"), RELEASE_EXECUTABLE_NAME); err == nil {
		t.Fatal("invalid zip accepted")
	}
	if _, err := helper.extractExecutableFromZip(zipArchive(t, map[string][]byte{"other": payload}), RELEASE_EXECUTABLE_NAME); err == nil {
		t.Fatal("missing executable accepted")
	}

	current := semver.New("1", "2", "0", "")
	latest := semver.New("1", "3", "0", "")
	required := semver.New("1", "1", "0", "")
	helper.currentVer = current
	if !helper.isEligibleForUpdate(latest, required) {
		t.Fatal("compatible upgrade was rejected")
	}
	helper.appliedVer = latest
	if helper.isEligibleForUpdate(latest, required) {
		t.Fatal("already applied latest was eligible")
	}
	helper.appliedVer = semver.New("1", "2", "5", "")
	if !helper.isEligibleForUpdate(latest, required) {
		t.Fatal("older applied version blocked a newer release")
	}
	helper.currentVer = semver.New("2", "0", "0", "")
	if helper.isEligibleForUpdate(latest, required) {
		t.Fatal("different major version was eligible")
	}
}

func TestApplyUpgrade(t *testing.T) {
	helper := testHelper(t, semver.New("1", "2", "0", ""))
	if err := helper.ApplyUpgrade(nil, []byte("x")); err == nil {
		t.Fatal("nil version accepted")
	}

	dir := t.TempDir()
	exe := filepath.Join(dir, "observer")
	if err := os.WriteFile(exe, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(exe)
	if err != nil {
		t.Fatal(err)
	}
	helper.currentExePath = exe
	version := semver.New("1", "3", "0", "")
	if err := helper.ApplyUpgrade(version, []byte("new")); err != nil {
		t.Fatal(err)
	}
	body, err := readFile(exe)
	if err != nil || string(body) != "new" {
		t.Fatalf("replaced = %q, %v", body, err)
	}
	info, err := os.Stat(exe)
	if err != nil || info.Mode().Perm() != before.Mode().Perm() {
		t.Fatalf("mode = %v, want %v, err %v", info.Mode().Perm(), before.Mode().Perm(), err)
	}
	if helper.appliedVer == nil || !helper.appliedVer.Equal(version) {
		t.Fatal("applied version was not recorded")
	}
	if err := helper.ApplyUpgrade(version, []byte("ignored")); err != nil {
		t.Fatal(err)
	}
	body, err = readFile(exe)
	if err != nil || string(body) != "new" {
		t.Fatalf("repeat apply changed the file to %q, %v", body, err)
	}

	next := semver.New("1", "3", "1", "")
	if err := helper.ApplyUpgrade(next, []byte("newer")); err != nil {
		t.Fatal(err)
	}
	body, err = readFile(exe)
	if err != nil || string(body) != "newer" || !helper.appliedVer.Equal(next) {
		t.Fatalf("second apply = %q %v %v", body, helper.appliedVer, err)
	}
	backups, err := filepath.Glob(filepath.Join(dir, "observer.*"))
	if err != nil || len(backups) == 0 {
		t.Fatalf("backups = %v, %v", backups, err)
	}

	locked := t.TempDir()
	if err := os.Chmod(locked, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o600) })
	helper.currentExePath = filepath.Join(locked, "observer")
	helper.appliedVer = nil
	if err := helper.ApplyUpgrade(version, []byte("nope")); err == nil {
		t.Fatal("read-only directory accepted an upgrade")
	}

	missingDir := t.TempDir()
	helper.currentExePath = filepath.Join(missingDir, "observer")
	if err := helper.ApplyUpgrade(semver.New("1", "4", "0", ""), []byte("orphan")); err == nil {
		t.Fatal("missing executable was replaced")
	}
	leftovers, err := os.ReadDir(missingDir)
	if err != nil || len(leftovers) != 0 {
		t.Fatalf("temporary upgrade files = %v, %v", leftovers, err)
	}
}

func testHelper(t *testing.T, current *semver.Version) *Helper {
	t.Helper()
	return &Helper{
		versionCheckDomain: "updates.example.test",
		releaseFetchUrl:    RELEASE_FETCH_URL_TEMPLATE,
		currentExePath:     filepath.Join(t.TempDir(), "observer"),
		currentBuild:       unibuild.New("test-toolchain", "stable", "abc123", "1700000000"),
		currentVer:         current,
		latestVer:          cache.NewGeneric[*semver.Version](time.Hour),
		requiredVer:        cache.NewGeneric[*semver.Version](time.Hour),
	}
}

func helperWithTXT(t *testing.T, txt string) *Helper {
	t.Helper()
	helper := testHelper(t, semver.New("1", "2", "0", ""))
	helper.resolvers = dnsquery.Resolvers{{Name: "local", Server: serveDNS(t, dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		reply := new(dns.Msg)
		reply.SetReply(r)
		if len(r.Question) > 0 {
			reply.Answer = []dns.RR{&dns.TXT{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 30},
				Txt: []string{txt},
			}}
		}
		_ = w.WriteMsg(reply)
	}))}}
	return helper
}

func serveDNS(t *testing.T, handler dns.Handler) string {
	t.Helper()
	packet, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &dns.Server{PacketConn: packet, Handler: handler}
	go func() { _ = server.ActivateAndServe() }()
	t.Cleanup(func() { _ = server.Shutdown() })
	return "udp://" + packet.LocalAddr().String()
}

func zipArchive(t *testing.T, files map[string][]byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	for name, body := range files {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write(body); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

const (
	checksumPayloadMD5  = "8a666fd414abefc1bc4c706fe0d2e3f1"
	checksumPayloadSHA1 = "983317dfc918d601faacc29f45039d36ac21fcaa"
)

func writeResponse(w http.ResponseWriter, body []byte) {
	_, _ = io.Copy(w, bytes.NewReader(body))
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func checksumText(payload []byte, upper bool) string {
	encode := func(sum []byte) string {
		text := hex.EncodeToString(sum)
		if upper {
			return strings.ToUpper(text)
		}
		return text
	}
	return strings.Join([]string{
		"SHA2-256=" + encode(sha256Sum(payload)),
		"SHA2-512=" + encode(sha512Sum(payload)),
		"NOTE=ignored",
	}, "\n")
}

func sha256Sum(payload []byte) []byte {
	sum := sha256.Sum256(payload)
	return sum[:]
}

func sha512Sum(payload []byte) []byte {
	sum := sha512.Sum512(payload)
	return sum[:]
}
