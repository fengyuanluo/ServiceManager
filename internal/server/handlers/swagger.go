package handlers

import (
	"net/http"

	_ "github.com/fengyuanluo/go-backend-template/docs"
	"github.com/fengyuanluo/go-backend-template/internal/server/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	router.NewGroupRouter("/api/swagger").
		AddRoute(
			router.NewRoute("/*any", http.MethodGet).
				Handle(ginSwagger.WrapHandler(swaggerFiles.Handler)),
		)
}
