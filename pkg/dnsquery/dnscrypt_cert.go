package dnsquery

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ameshkov/dnscrypt/v2"
	"github.com/ameshkov/dnscrypt/v2/xsecretbox"
	"github.com/jedisct1/go-dnsstamps"
	"github.com/miekg/dns"
	"golang.org/x/crypto/nacl/box"
)

// The library's Dial performs uncancelable certificate I/O. Fetch over a
// context-bound TCP connection, then use its certificate and crypto primitives.
func dialDNSCrypt(ctx context.Context, endpoint string, timeout time.Duration) (*dnscrypt.ResolverInfo, error) {
	stamp, err := dnsstamps.NewServerStampFromString(endpoint)
	if err != nil {
		return nil, err
	}
	if stamp.Proto != dnsstamps.StampProtoTypeDNSCrypt || len(stamp.ServerPk) != 32 {
		return nil, errors.New("endpoint must contain a DNSCrypt stamp with a valid public key")
	}
	request := new(dns.Msg).SetQuestion(dns.Fqdn(stamp.ProviderName), dns.TypeTXT)
	response, err := exchangeContext(ctx, &dns.Client{Net: "tcp", Timeout: timeout}, request, stamp.ServerAddrStr)
	if err != nil {
		return nil, err
	}
	if err := validateResponse(request, response); err != nil {
		return nil, err
	}
	var certificate *dnscrypt.Cert
	for _, answer := range response.Answer {
		txt, ok := answer.(*dns.TXT)
		if !ok || !strings.EqualFold(txt.Hdr.Name, request.Question[0].Name) || txt.Hdr.Class != dns.ClassINET {
			continue
		}
		data, err := certificateBytes(txt)
		if err != nil {
			continue
		}
		candidate := &dnscrypt.Cert{}
		if candidate.Deserialize(data) != nil || !candidate.VerifyDate() || !candidate.VerifySignature(stamp.ServerPk) {
			continue
		}
		if certificate == nil || candidate.Serial > certificate.Serial || (candidate.Serial == certificate.Serial && candidate.EsVersion > certificate.EsVersion) {
			certificate = candidate
		}
	}
	if certificate == nil {
		return nil, errors.New("no valid signed DNSCrypt certificate")
	}
	public, secret, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DNSCrypt key pair: %w", err)
	}
	resolver := &dnscrypt.ResolverInfo{
		SecretKey: *secret, PublicKey: *public, ServerPublicKey: stamp.ServerPk,
		ServerAddress: stamp.ServerAddrStr, ProviderName: stamp.ProviderName, ResolverCert: certificate,
	}
	switch certificate.EsVersion {
	case dnscrypt.XSalsa20Poly1305:
		box.Precompute(&resolver.SharedKey, &certificate.ResolverPk, secret)
	case dnscrypt.XChacha20Poly1305:
		resolver.SharedKey, err = xsecretbox.SharedKey(*secret, certificate.ResolverPk)
	default:
		return nil, errors.New("unsupported DNSCrypt crypto construction")
	}
	if err != nil {
		return nil, err
	}
	return resolver, nil
}

// Repacking lets miekg/dns decode TXT escapes and character-string boundaries.
// The certificate is binary data, not the escaped presentation in TXT.Txt.
func certificateBytes(txt *dns.TXT) ([]byte, error) {
	record := *txt
	record.Hdr.Name = "."
	wire := make([]byte, dns.Len(&record))
	n, err := dns.PackRR(&record, wire, 0, nil, false)
	if err != nil {
		return nil, err
	}
	const headerSize = 11 // Root label plus TYPE, CLASS, TTL and RDLENGTH.
	var data []byte
	for offset := headerSize; offset < n; {
		length := int(wire[offset])
		offset++
		if offset+length > n {
			return nil, errors.New("invalid DNSCrypt certificate TXT data")
		}
		data = append(data, wire[offset:offset+length]...)
		offset += length
	}
	return data, nil
}
