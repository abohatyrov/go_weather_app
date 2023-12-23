package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
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

	totalAccesses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "web_user_access_total",
			Help: "Total number of user accesses to the website",
		},
		[]string{"user_ip"},
	)
)

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(pageViews)
	prometheus.MustRegister(totalAccesses)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIP := c.ClientIP()

		totalAccesses.WithLabelValues(userIP).Inc()

		totalRequests.WithLabelValues().Inc()

		pageViews.WithLabelValues(c.FullPath()).Inc()

		c.Next()
	}
}

func MetricsHandler(c *gin.Context) {
	log.Println("Handling metrics request")
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
