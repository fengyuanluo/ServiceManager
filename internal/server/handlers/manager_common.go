package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fengyuanluo/go-backend-template/internal/server/middleware"
	"github.com/fengyuanluo/go-backend-template/internal/server/resp"
	"github.com/fengyuanluo/go-backend-template/internal/server/router"
	"github.com/fengyuanluo/go-backend-template/internal/service/manager"
	"github.com/gin-gonic/gin"
)

var managerExecSource = manager.ExecSource{}
var managerEnvSource = manager.EnvSource{ReadFile: os.ReadFile}

type batchRequest struct {
	Items []manager.BatchItem `json:"items"`
	User  string              `json:"user,omitempty"`
}

func registerManagerRoutes(key string) {
	router.NewGroupRouter("/api/v1/" + key).
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/list", http.MethodGet).
				Handle(listHandler(key)),
		).
		AddRoute(
			router.NewRoute("/detail", http.MethodGet).
				Handle(detailHandler(key)),
		).
		AddRoute(
			router.NewRoute("/logs", http.MethodGet).
				Handle(logsHandler(key)),
		).
		AddRoute(
			router.NewRoute("/action", http.MethodPost).
				Handle(actionHandler(key)),
		).
		AddRoute(
			router.NewRoute("/batch", http.MethodPost).
				Handle(batchHandler(key)),
		).
		AddRoute(
			router.NewRoute("/install", http.MethodPost).
				Handle(installHandler(key)),
		)
}

func listHandler(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := strings.TrimSpace(c.Query("user"))
		services, err := manager.ListServices(key, managerExecSource, user)
		if err != nil {
			resp.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		resp.Success(c, services)
	}
}

func detailHandler(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.TrimSpace(c.Query("name"))
		if name == "" {
			resp.Error(c, http.StatusBadRequest, resp.ErrBadRequest)
			return
		}
		user := strings.TrimSpace(c.Query("user"))
		detail, err := manager.GetServiceDetail(key, managerExecSource, name, user)
		if err != nil {
			resp.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		resp.Success(c, detail)
	}
}

func logsHandler(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.TrimSpace(c.Query("name"))
		if name == "" {
			resp.Error(c, http.StatusBadRequest, resp.ErrBadRequest)
			return
		}
		lines, _ := strconv.Atoi(c.DefaultQuery("lines", "200"))
		keyword := strings.TrimSpace(c.Query("keyword"))
		user := strings.TrimSpace(c.Query("user"))
		query := manager.LogQuery{Lines: lines, Keyword: keyword, User: user}
		logs, err := manager.GetServiceLogs(key, managerExecSource, name, query)
		if err != nil {
			resp.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		resp.Success(c, logs)
	}
}

func actionHandler(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload manager.ActionRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			resp.Error(c, http.StatusBadRequest, resp.ErrInvalidJSON)
			return
		}
		if strings.TrimSpace(payload.User) == "" {
			payload.User = strings.TrimSpace(c.Query("user"))
		}
		action := strings.ToLower(strings.TrimSpace(payload.Action))
		if payload.Name == "" && action != "autostart" && action != "enable" {
			resp.Error(c, http.StatusBadRequest, resp.ErrBadRequest)
			return
		}
		if err := manager.RunServiceAction(key, managerExecSource, payload); err != nil {
			resp.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		resp.Success(c, "ok")
	}
}

func batchHandler(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload batchRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			resp.Error(c, http.StatusBadRequest, resp.ErrInvalidJSON)
			return
		}
		if len(payload.Items) == 0 {
			resp.Error(c, http.StatusBadRequest, resp.ErrBadRequest)
			return
		}
		results := manager.RunServiceBatch(key, managerExecSource, payload.User, payload.Items)
		resp.Success(c, results)
	}
}

func installHandler(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		env := manager.DetectEnvironment(managerEnvSource)
		result := manager.InstallManager(key, env, managerExecSource)
		resp.Success(c, result)
	}
}
