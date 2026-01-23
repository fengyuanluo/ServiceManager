package handlers

import (
	"net/http"

	"github.com/fengyuanluo/go-backend-template/internal/model"
	"github.com/fengyuanluo/go-backend-template/internal/op"
	"github.com/fengyuanluo/go-backend-template/internal/server/auth"
	"github.com/fengyuanluo/go-backend-template/internal/server/middleware"
	"github.com/fengyuanluo/go-backend-template/internal/server/resp"
	"github.com/fengyuanluo/go-backend-template/internal/server/router"
	"github.com/gin-gonic/gin"
)

func init() {
	router.NewGroupRouter("/api/v1/user").
		AddRoute(
			router.NewRoute("/login", http.MethodPost).
				Handle(login),
		)
	router.NewGroupRouter("/api/v1/user").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/change-password", http.MethodPost).
				Handle(changePassword),
		).
		AddRoute(
			router.NewRoute("/change-username", http.MethodPost).
				Handle(changeUsername),
		)
}

// @Summary 用户登录
// @Description 校验用户名密码并返回 JWT
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body model.UserLogin true "登录参数"
// @Success 200 {object} resp.ResponseStruct{data=model.UserLoginResponse} "成功"
// @Failure 400 {object} resp.ResponseStruct "请求参数错误"
// @Failure 401 {object} resp.ResponseStruct "认证失败"
// @Failure 500 {object} resp.ResponseStruct "服务器错误"
// @Router /v1/user/login [post]
func login(c *gin.Context) {
	var user model.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.Error(c, http.StatusBadRequest, resp.ErrInvalidJSON)
		return
	}
	if err := op.UserVerify(user.Username, user.Password); err != nil {
		resp.Error(c, http.StatusUnauthorized, resp.ErrUnauthorized)
		return
	}
	token, expire, err := auth.GenerateToken(user.Expire)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, resp.ErrInternalServer)
		return
	}
	resp.Success(c, model.UserLoginResponse{Token: token, ExpireAt: expire})
}

// @Summary 修改密码
// @Description 需要携带 Authorization 的 Bearer Token
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.UserChangePassword true "修改密码参数"
// @Success 200 {object} resp.ResponseStruct{data=string} "成功"
// @Failure 400 {object} resp.ResponseStruct "请求参数错误"
// @Failure 401 {object} resp.ResponseStruct "认证失败"
// @Failure 500 {object} resp.ResponseStruct "服务器错误"
// @Router /v1/user/change-password [post]
func changePassword(c *gin.Context) {
	var user model.UserChangePassword
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.Error(c, http.StatusBadRequest, resp.ErrInvalidJSON)
		return
	}
	if err := op.UserChangePassword(user.OldPassword, user.NewPassword); err != nil {
		resp.Error(c, http.StatusInternalServerError, resp.ErrDatabase)
		return
	}
	resp.Success(c, "password changed successfully")
}

// @Summary 修改用户名
// @Description 需要携带 Authorization 的 Bearer Token
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.UserChangeUsername true "修改用户名参数"
// @Success 200 {object} resp.ResponseStruct{data=string} "成功"
// @Failure 400 {object} resp.ResponseStruct "请求参数错误"
// @Failure 401 {object} resp.ResponseStruct "认证失败"
// @Failure 500 {object} resp.ResponseStruct "服务器错误"
// @Router /v1/user/change-username [post]
func changeUsername(c *gin.Context) {
	var user model.UserChangeUsername
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.Error(c, http.StatusBadRequest, resp.ErrInvalidJSON)
		return
	}
	if err := op.UserChangeUsername(user.NewUsername); err != nil {
		resp.Error(c, http.StatusInternalServerError, resp.ErrDatabase)
		return
	}
	resp.Success(c, "username changed successfully")
}
