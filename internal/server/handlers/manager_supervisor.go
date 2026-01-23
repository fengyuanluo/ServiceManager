package handlers

import (
	"net/http"

	"github.com/fengyuanluo/go-backend-template/internal/server/middleware"
	"github.com/fengyuanluo/go-backend-template/internal/server/resp"
	"github.com/fengyuanluo/go-backend-template/internal/server/router"
	"github.com/fengyuanluo/go-backend-template/internal/service/manager"
	"github.com/gin-gonic/gin"
)

func init() {
	router.NewGroupRouter("/api/v1/supervisor").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/list", http.MethodGet).
				Handle(listSupervisor),
		)
}

func listSupervisor(c *gin.Context) {
	resp.Success(c, []manager.ServiceSummary{})
}
