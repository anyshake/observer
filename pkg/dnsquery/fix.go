package dnsquery

import "net/url"

func fixPortNumber(urlObj *url.URL, defaultPort string) *url.URL {
	if urlObj.Port() == "" {
		var newUrl url.URL
		newUrl.Scheme = urlObj.Scheme
		newUrl.Host = urlObj.Host + ":" + defaultPort
		urlObj = &newUrl
	}

	return urlObj
}
