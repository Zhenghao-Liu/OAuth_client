package handler

import (
	"encoding/json"
	"errors"
	"github.com/Zhenghao-Liu/OAuth_client/common"
	"github.com/Zhenghao-Liu/OAuth_client/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type TokenResp struct {
	Token   string `json:"access_token"`
	Refresh string `json:"refresh_token"`
	Final   string `json:"result"`
}

type ResourceResp struct {
	Resource1 string `json:"resource1"`
	Resource2 string `json:"resource2"`
	Resource3 string `json:"resource3"`
	Final     string `json:"result"`
}

func RegisterHandler(router *gin.Engine) {

	router.GET("/welcome", func(ctx *gin.Context) {
		common.State = utils.GenString()
		tarUrl := common.OAuthPage +
			common.OAuthWelcome +
			"?app_id=g3*sj1qdrr%40sdhm-nes%26qn5shpgwy2c%3Dd3z6%5E7ymr1zzuew%3Dpd))%40(qwcd%3Ducjqq&" +
			"redirect_url=" + common.Callback + "&" +
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
		if err := sendCode(ctx); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		resourceResp, err := getResource(ctx)
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
		resourceResp, err := getResource(ctx)
		if err != nil {
			if err = sendRefresh(ctx); err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			resourceResp, err = getResource(ctx)
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

func sendCode(ctx *gin.Context) error {
	code, _ := url.QueryUnescape(ctx.DefaultQuery("code", ""))
	common.Code = code
	state := ctx.DefaultQuery("state", "")
	respErr := ctx.DefaultQuery("error", "")
	if respErr != "" {
		return errors.New(respErr)
	} else if code == "" || state == "" || state != common.State {
		return errors.New(common.LogInErr)
	}
	client := &http.Client{}
	data := url.Values{
		"grant_type": {common.AuthorizationCode},
		"callback":   {common.Callback},
	}
	req, err := http.NewRequest(http.MethodPost, common.OAuthPage+common.OAuthToken, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("app_id", common.AppID)
	req.Header.Set("app_secret", common.AppSecret)
	req.Header.Set("code", common.Code)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	tokenResp := &TokenResp{}
	if err = json.Unmarshal(body, &tokenResp); err != nil {
		return err
	} else if tokenResp.Final != common.StatusSuccess {
		return err
	}
	common.Token = tokenResp.Token
	common.Refresh = tokenResp.Refresh
	return nil
}

func getResource(ctx *gin.Context) (*ResourceResp, error) {
	client := &http.Client{}
	data := url.Values{
		"token_type": {common.TokenType},
	}
	req, err := http.NewRequest(http.MethodPost, common.OAuthPage+common.OAuthResource, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("app_id", common.AppID)
	req.Header.Set("app_secret", common.AppSecret)
	req.Header.Set("access_token", common.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resourceResp := &ResourceResp{}
	if err = json.Unmarshal(body, &resourceResp); err != nil {
		return nil, err
	} else if resourceResp.Final != common.StatusSuccess {
		return nil, errors.New(resourceResp.Final)
	}
	return resourceResp, nil
}

func sendRefresh(ctx *gin.Context) error {
	client := &http.Client{}
	data := url.Values{
		"grant_type": {common.RefreshToken},
		"callback":   {common.Callback},
	}
	req, err := http.NewRequest(http.MethodPost, common.OAuthPage+common.OAuthRefresh, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("app_id", common.AppID)
	req.Header.Set("app_secret", common.AppSecret)
	req.Header.Set("refresh_token", common.Refresh)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	tokenResp := &TokenResp{}
	if err = json.Unmarshal(body, &tokenResp); err != nil {
		return err
	} else if tokenResp.Final != common.StatusSuccess {
		return err
	}
	common.Token = tokenResp.Token
	common.Refresh = tokenResp.Refresh
	return nil
}
