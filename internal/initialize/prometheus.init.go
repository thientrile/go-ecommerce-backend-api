package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-ecommerce-backend-api.com/global"
)

var (
	OpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func InitPrometheus(r *gin.Engine) *gin.Engine {
	// prometheus.MustRegister(OpsProcessed)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	global.OpsProcessed = OpsProcessed
	return r
}
