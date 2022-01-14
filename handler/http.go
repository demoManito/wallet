package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	goloader "wallet/pkg/core/loader"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "requests_total",
			Subsystem: "http",
			Help:      "statistics for http requests",
		},
		[]string{"code", "method"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "request_duration_seconds",
			Subsystem: "http",
			Help:      "http request duration",
		},
		[]string{"code", "method"},
	)
)

type Server struct {
	httpSrv         *http.Server
	shutdownTimeout int
}

func NewServer() *Server {
	h := GetHandler()
	goloader.LoadingAll(goloader.DefaultLoader())
	loader := goloader.DefaultLoader()
	loader.Loading(h)

	router := gin.New()
	routerV1 := router.Group("/api/v1")
	{
		routerWallet := routerV1.Group("/wallet")
		routerWallet.POST("/draw_money", h.HandlerV1.DrawMoney)
		routerWallet.POST("/save_money", h.HandlerV1.SaveMoney)
		routerWallet.POST("/transfer", h.HandlerV1.Transfer)
		routerWallet.POST("/status", h.HandlerV1.Status)
	}
	routerV2 := router.Group("/api/v2")
	{
		routerWallet := routerV2.Group("/wallet")
		routerWallet.POST("/draw_money", h.HandlerV2.DrawMoney)
		routerWallet.POST("/save_money", h.HandlerV2.SaveMoney)
		routerWallet.POST("/transfer", h.HandlerV2.Transfer)
	}
	var handler http.Handler = router
	handler = promhttp.InstrumentHandlerCounter(HttpRequestsTotal, promhttp.InstrumentHandlerDuration(HttpRequestDuration, handler))
	return &Server{
		httpSrv: &http.Server{
			Addr:    h.Config.Http.Port,
			Handler: handler,
		},
		shutdownTimeout: 10,
	}
}

func (s *Server) Run() error {
	go func() {
		logrus.Infof("HTTP server listen: %s", s.httpSrv.Addr)
		if err := s.httpSrv.ListenAndServe(); err != nil {
			logrus.WithError(err).Errorf("start http server failed")
		}
	}()
	return nil
}
