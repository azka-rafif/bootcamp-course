package http

import (
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/container"
	"github.com/evermos/boilerplate-go/docs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Http struct {
	Config *configs.Config
	DB     *infras.MysqlConn
	Router *chi.Mux
}

func (h *Http) Shutdown(done chan os.Signal, svc container.ServiceRegistry) {
	<-done
	defer os.Exit(0)
	time.Sleep(time.Duration(5) * time.Second)
	log.Println("Clean up all resources...")
	svc.Shutdown()
	log.Println("Server shutdown properly...")
}

func (h *Http) GracefulShutdown(svc container.ServiceRegistry) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go h.Shutdown(done, svc)
}

func (h *Http) Serve() {
	log.Println("Running service on port : ", h.Config.Port)
	h.Router.Get("/health", h.HealthCheck)
	if h.Config.Env == shared.EnvirontmentDev {
		docs.SwaggerInfo.Title = shared.ServiceName
		docs.SwaggerInfo.Version = shared.ServiceVersion
		swaggerURL := fmt.Sprintf("%s/swagger/doc.json", h.Config.AppURL)
		h.Router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(swaggerURL)))
	}

	netHttp.ListenAndServe(":"+h.Config.Port, h.Router)
}

func (h *Http) HealthCheck(w netHttp.ResponseWriter, r *netHttp.Request) {
	if err := h.DB.Read.Ping(); err != nil {
		shared.Failed(w, r, shared.WriteFailed(500, "Server is unhealthy", nil))
		return
	}
	shared.Success(w, r, shared.WriteSuccess(200, "Server is alive", nil))
}
