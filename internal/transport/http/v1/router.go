package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase"

	_ "gitlab.com/fluxx1on_group/event_message_service/docs"
)

var RequestsTotalCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_total",
		Help: "HTTP Responses",
	},
	[]string{"method", "endpoint", "code"},
)

/*
# HELP requests_total HTTP Failures
# TYPE requests_total counter
requests_total{code="500",endpoint="/v1/client",method="PUT"} 1
*/

func pushMetric(method string, endpoint string, code int) {
	RequestsTotalCounter.With(
		prometheus.Labels{
			"method":   method,
			"endpoint": endpoint,
			"code":     strconv.Itoa(code),
		},
	).Inc()
}

const basePath string = "/v1"

// NewRouter -.
// Swagger spec:
// @title       Go Mailing Service
// @description Closed API
// @version     1.0

// @contact.name   Nikolai Kaliga
// @contact.email  nick.kaliga@ya.ru

// @host      	localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, client usecase.Client, mailing usecase.Mailing) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	handler.GET("/docs/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	prometheus.MustRegister(RequestsTotalCounter)

	// Routers
	h := handler.Group(basePath)
	{
		newClientRoutes(h, client)
		newMailingRoutes(h, mailing)
	}
}
