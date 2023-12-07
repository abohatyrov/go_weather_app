package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_total_requests",
			Help: "Total number of requests to the web server",
		},
		[]string{},
	)

	pageViews = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "web_page_views_total",
			Help: "Total number of page views",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(pageViews)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		totalRequests.WithLabelValues().Inc()

		pageViews.WithLabelValues(c.FullPath()).Inc()

		c.Next()
	}
}

func MetricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
