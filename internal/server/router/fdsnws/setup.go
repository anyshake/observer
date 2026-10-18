package fdsnws

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
	"github.com/bclswl0827/fdsnws-go/service/station"
	"github.com/gin-gonic/gin"
)

const (
	bearerChallenge = `Bearer realm="fdsnws"`
	basePath        = "/fdsnws"
)

func Setup(engine *gin.Engine, handler *action.Handler, device hardware.IHardware, now func() time.Time) error {
	modules := []fdsnwsgo.Module{fdsnwsgo.ModuleStation, fdsnwsgo.ModuleDataSelect}
	dataProvider := &provider{handler: handler, device: device, now: now}
	stationService := station.New(dataProvider)
	dataSelectService := newDataSelectService(dataProvider)
	server, err := fdsnwsgo.NewFDSNWS(
		dataProvider,
		fdsnwsgo.WithBasePath(basePath),
		fdsnwsgo.WithModules(modules...),
		fdsnwsgo.WithService(stationService),
		fdsnwsgo.WithService(dataSelectService),
	)
	if err != nil {
		return err
	}

	routes := engine.Group(basePath, serviceMiddleware(handler, server), accessKeyMiddleware(handler, now, server, stationService, dataSelectService))

	assets, err := fs.Sub(webpageAssets, "static/assets")
	if err != nil {
		return err
	}
	routes.StaticFS("/assets", &gin.OnlyFilesFS{FileSystem: http.FS(assets)})

	serve := gin.WrapH(server)
	routes.Any("", serve)
	routes.Any("/", serve)
	for _, module := range modules {
		modulePath := "/" + string(module)
		routes.Any(modulePath, serve)
		routes.Any(modulePath+"/*path", serve)
	}
	return nil
}
