/**
* @author mnunez
 */

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/matiasnu/go-jopit-toolkit/goutils/apierrors"
)

var production bool = os.Getenv("GO_ENVIRONMENT") == "production"

func DefaultJopitRouter() *gin.Engine {
	return CustomJopitRouter(JopitRouterConfig{})
}

func CustomJopitRouter(conf JopitRouterConfig) *gin.Engine {
	router := gin.New()

	if conf.DisableCancellationOnClientDisconnect {
		router.Use(func(c *gin.Context) {
			c.Request = c.Request.WithContext(context.Background())
			c.Next()
		})
	}

	if conf.EnableResponseCompressionSupport {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	if !conf.DisableCommonApiFilter {
		router.Use(CommonAPiFilter(!conf.DisableCommonApiFilterErrorLog))
	}
	if !conf.DisablePprof {
		pprof.Register(router)
	}
	if !conf.DisableHeaderForwarding {
		router.Use(HeaderForwarding())
	}
	if !production {
		router.Use(gin.Logger())
	}

	router.NoRoute(noRouteHandler)
	return router
}

type JopitRouterConfig struct {
	DisableCommonApiFilter           bool
	DisablePprof                     bool
	EnableResponseCompressionSupport bool
	DisableHeaderForwarding          bool

	// DisableCommonApiFilterErrorLog tells the Common Api Filter to omit logging the URI in error handling.
	// This is useful when some params or query params are private data, like tokens.
	DisableCommonApiFilterErrorLog bool

	// DisableCancellationOnClientDisconnect tells the server to detach the
	// c.Request.Context() from the incoming TCP connection. If set to false
	// then the client closing the connection does not cancels the context.
	// The default behavior from Go is to cancel the request context if it can
	// ensure that there's no one on the other side to read the response.
	DisableCancellationOnClientDisconnect bool
}

func noRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, apierrors.NewNotFoundApiError(fmt.Sprintf("Resource %s not found.", c.Request.URL.Path)))
}

func AddResponseExpiration(time time.Duration, c *gin.Context) {
	var roundTime int = int(time.Seconds())
	c.Writer.Header()["Cache-Control"] = []string{fmt.Sprintf("max-age=%v,stale-while-revalidate=%v, stale-if-error=%v", roundTime, roundTime/2, roundTime*2)}
}
