package api

import (
	"fmt"
	"net/http"

	"github.com/Depado/ginprom"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/services"
	"github.com/rs/zerolog"
)

// API is the struct that holds the router that handles incoming http Requests
type API struct {
	router   *gin.Engine
	log      *config.Logger
	config   *config.Config
	services *services.Services
}

// New creates a new API server that contains all of the supported endpoints of the adapult
// service.
func New(conf *config.Config, svcs *services.Services) *API {
	// cLog.Logger = cLog.Logger.With().Str("component", "api").Logger()
	api := &API{
		config: conf,
		// log:      cLog,
		services: svcs,
	}

	r := gin.New()

	if conf.Environment == "development" {
		// use colorized output during local development
		r.Use(logger.SetLogger())
	} else {
		r.Use(func(c *gin.Context) {
			skipPaths := map[string]bool{
				"/health":  true,
				"/metrics": true,
			}

			// Skip logging for the paths in the skipPaths map
			if _, ok := skipPaths[c.Request.URL.Path]; ok {
				c.Next()
				return
			}

			// c.Next() calls the next handlers in the chain, and the status code is set somewhere in these handlers.
			// This allows the request to be processed by the rest of the handlers before the logging middleware
			c.Next()

			// Skip logging for successful HTTP 200 responses
			if c.Writer.Status() == http.StatusOK {
				return
			}

			logger.SetLogger(logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
				return l.Output(gin.DefaultWriter).With().Logger()
			}))(c)
		})
	}
	r.Use(gin.Recovery())

	// This condition has been implemented to avoid that the ginprom
	// configuration loads more than once inside the unit tests
	// We got the next error message after implement ginprom library
	// `panic: duplicate metrics collector registration attempted`
	//
	if conf.Environment == "production" || conf.Environment == "staging" {
		// expose metrics
		p := ginprom.New(
			ginprom.Engine(r),
			// ginprom.Token(conf.MonitoringToken),
			ginprom.Subsystem("adapult"),
			ginprom.Path("/metrics"),
		)
		r.Use(p.Instrument())
	}

	// Routes

	r.GET("/health", api.Health())
	r.POST("/users/position", api.RecordPosition)
	r.GET("/ws", api.ws)

	api.router = r

	return api
}

// NotImplemented is used as placeholder for endpoints that are not implemented
func NotImplemented(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "route found but not implemented"})
}

// Start starts the REST API
func (api *API) Start() {
	err := api.router.Run(fmt.Sprintf(":%d", api.config.HTTP.Port))
	if err != nil {
		api.log.Fatal().Msgf("Could not start API server: %v", err)
	}
}

// Router returns the gin router instance
func (api *API) Router() *gin.Engine {
	return api.router
}

// Health returns the status of the application
func (api *API) Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}
