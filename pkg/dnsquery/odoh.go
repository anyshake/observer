package dnsquery

// import (
// 	"bytes"
// 	"crypto/tls"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strings"

// 	"github.com/labstack/gommon/log"
// 	"github.com/miekg/dns"
// 	"github.com/sthorne/odoh-go"
// )

// const ODoHContentType = "application/oblivious-dns-message"

// // buildURL adds HTTPS to argument s if it doesn't contain a protocol and appends defaultPath if no path is already specified
// func buildURL(s, defaultPath string) *url.URL {
// 	if //goland:noinspection HttpUrlsUsage
// 	!strings.HasPrefix(s, "https://") && !strings.HasPrefix(s, "http://") {
// 		s = "https://" + s
// 	}
// 	u, err := url.Parse(s)
// 	if err != nil {
// 		log.Fatalf("failed to parse url: %v", err)
// 	}
// 	if u.Path == "" {
// 		u.Path = defaultPath
// 	}
// 	return u
// }

// // ODoH makes a DNS query over ODoH
// type ODoH struct {
// 	Common    // Server is the target
// 	Proxy     string
// 	TLSConfig *tls.Config

// 	conn *http.Client
// }

// func (o *ODoH) Exchange(m *dns.Msg) (*dns.Msg, error) {
// 	// Query ODoH configs on target
// 	req, err := http.NewRequest(
// 		http.MethodGet,
// 		buildURL(strings.TrimSuffix(o.Server, "/dns-query"), "/.well-known/odohconfigs").String(),
// 		nil,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("new target configs request: %s", err)
// 	}

// 	if o.conn == nil || !o.ReuseConn {
// 		o.conn = &http.Client{
// 			Transport: &http.Transport{
// 				TLSClientConfig: o.TLSConfig,
// 			},
// 		}
// 	}
// 	resp, err := o.conn.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("do target configs request: %s", err)
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	odohConfigs, err := odoh.UnmarshalObliviousDoHConfigs(bodyBytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("unmarshal target configs: %s", err)
// 	}

// 	if len(odohConfigs.Configs) == 0 {
// 		return nil, errors.New("target provided no valid ODoH configs")
// 	}
// 	log.Debugf("[odoh] retrieved %d ODoH configs", len(odohConfigs.Configs))

// 	packedDnsQuery, err := m.Pack()
// 	if err != nil {
// 		return nil, err
// 	}

// 	firstODoHConfig := odohConfigs.Configs[0]
// 	log.Debugf("[odoh] using first ODoH config: %+v", firstODoHConfig)
// 	odnsMessage, queryContext, err := firstODoHConfig.Contents.EncryptQuery(odoh.CreateObliviousDNSQuery(packedDnsQuery, 0))
// 	if err != nil {
// 		return nil, fmt.Errorf("encrypt query: %s", err)
// 	}

// 	t := buildURL(o.Server, "/dns-query")
// 	p := buildURL(o.Proxy, "/proxy")
// 	qry := p.Query()
// 	if qry.Get("targethost") == "" {
// 		qry.Set("targethost", t.Host)
// 	}
// 	if qry.Get("targetpath") == "" {
// 		qry.Set("targetpath", t.Path)
// 	}
// 	p.RawQuery = qry.Encode()

// 	log.Debugf("POST %s %+v", p, odnsMessage)
// 	req, err = http.NewRequest(http.MethodPost, p.String(), bytes.NewBuffer(odnsMessage.Marshal()))
// 	if err != nil {
// 		return nil, fmt.Errorf("create new request: %s", err)
// 	}
// 	req.Header.Set("Content-Type", ODoHContentType)
// 	req.Header.Set("Accept", ODoHContentType)

// 	resp, err = o.conn.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("do request: %s", err)
// 	}
// 	contentType := resp.Header.Get("Content-Type")
// 	if contentType != ODoHContentType {
// 		return nil, fmt.Errorf("%s responded with an invalid Content-Type header %s, expected %s", req.URL, contentType, ODoHContentType)
// 	}

// 	bodyBytes, err = io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("read response body: %s", err)
// 	}
// 	odohMessage, err := odoh.UnmarshalDNSMessage(bodyBytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("odoh unmarshal: %s", err)
// 	}

// 	decryptedResponse, err := queryContext.OpenAnswer(odohMessage)
// 	if err != nil {
// 		return nil, fmt.Errorf("open answer: %s", err)
// 	}

// 	msg := &dns.Msg{}
// 	err = msg.Unpack(decryptedResponse)
// 	if err != nil {
// 		err = fmt.Errorf("unpack message: %s", err)
// 	}
// 	return msg, err
// }

// func (o *ODoH) Close() error {
// 	o.conn.CloseIdleConnections()
// 	return nil
// }
