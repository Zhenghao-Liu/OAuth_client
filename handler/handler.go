package handler

import (
	"github.com/Zhenghao-Liu/OAuth_client/common"
	"github.com/Zhenghao-Liu/OAuth_client/service"
	"github.com/Zhenghao-Liu/OAuth_client/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func RegisterHandler(router *gin.Engine) {
	router.LoadHTMLGlob("html/*")

	router.GET("/welcome", func(ctx *gin.Context) {
		common.State = utils.GenString()
		tarUrl := common.OAuthPage +
			common.OAuthWelcome +
			"?app_id=" + url.QueryEscape(common.AppID) + "&" +
			"callback=" + common.Callback + "&" +
			"response_type=code&" +
			"state=" + common.State
		ctx.HTML(http.StatusOK, "welcome.html", gin.H{
			"url": tarUrl,
		})
	})

	router.GET("/cancel", func(ctx *gin.Context) {
		common.Code = ""
		common.Token = ""
		common.State = ""
		ctx.String(http.StatusOK, common.StatusSuccess)
	})

	router.GET("/oauth/redirect", func(ctx *gin.Context) {
		if err := service.SendCode(ctx); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		resourceResp, err := service.GetResource(ctx)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.HTML(http.StatusOK, "resource.html", gin.H{
			"结果":        "授权成功，获取资源",
			"resource1": common.Resource1 + resourceResp.Resource1,
			"resource2": common.Resource2 + resourceResp.Resource2,
			"resource3": common.Resource3 + resourceResp.Resource3,
		})
	})

	router.GET("/oauth/resource", func(ctx *gin.Context) {
		resourceResp, err := service.GetResource(ctx)
		if err != nil {
			if err = service.SendRefresh(ctx); err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			resourceResp, err = service.GetResource(ctx)
		}
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.HTML(http.StatusOK, "resource.html", gin.H{
			"resource1": common.Resource1 + resourceResp.Resource1,
			"resource2": common.Resource2 + resourceResp.Resource2,
			"resource3": common.Resource3 + resourceResp.Resource3,
		})
	})
}
