package fdsnws

import (
	"embed"
	"net/http"

	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
)

var (
	//go:embed static/assets
	webpageAssets embed.FS
	//go:embed static/index.html
	fdsnwsRootPage string
	//go:embed static/dataselect/index.html
	dataSelectRootPage string
	//go:embed static/dataselect/1/index.html
	dataSelectVersionRootPage string
	//go:embed static/station/index.html
	stationRootPage string
	//go:embed static/station/1/index.html
	stationVersionRootPage string
)

func htmlPage(body string) (fdsnwsgo.PreparedResponse, fdsnwsgo.ProviderErrorKind, error) {
	return fdsnwsgo.PreparedResponse{
		Body:       body,
		StatusCode: http.StatusOK,
		MIMEType:   "text/html; charset=utf-8",
	}, fdsnwsgo.ProviderErrorNone, nil
}

func (*provider) FDSNWSRoot() (fdsnwsgo.PreparedResponse, fdsnwsgo.ProviderErrorKind, error) {
	return htmlPage(fdsnwsRootPage)
}

func (*provider) DataSelectRoot() (fdsnwsgo.PreparedResponse, fdsnwsgo.ProviderErrorKind, error) {
	return htmlPage(dataSelectRootPage)
}

func (*provider) DataSelectVersionRoot() (fdsnwsgo.PreparedResponse, fdsnwsgo.ProviderErrorKind, error) {
	return htmlPage(dataSelectVersionRootPage)
}

func (*provider) StationRoot() (fdsnwsgo.PreparedResponse, fdsnwsgo.ProviderErrorKind, error) {
	return htmlPage(stationRootPage)
}

func (*provider) StationVersionRoot() (fdsnwsgo.PreparedResponse, fdsnwsgo.ProviderErrorKind, error) {
	return htmlPage(stationVersionRootPage)
}
