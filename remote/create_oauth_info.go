package remote

import (
	"encoding/json"
	"fmt"
	"github.com/Zhenghao-Liu/OAuth_client/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CreateOAuthResp struct {
	AppID     string `json:"app_id"`
	Appsecret string `json:"app_secret"`
	Final     string `json:"result"`
}

func CreateOAuthInfo() {
	resp, err := http.PostForm(common.OAuthPage+common.OAuthCreate, url.Values{
		"app_name":    {"lzh毕设专用"},
		"homepage":    {common.HomePage},
		"description": {"毕设专用"},
		"callback":    {common.Callback},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	createOauthResp := &CreateOAuthResp{}
	if err = json.Unmarshal(body, &createOauthResp); err != nil {
		panic(err)
	} else if createOauthResp.Final != common.StatusSuccess {
		panic(createOauthResp.Final)
	}
	fmt.Println("app_id: ", common.AppID)
	fmt.Println("app_secret: ", common.AppSecret)
	panic("请修改const中的常量，并注释该函数，后重新运行")
}
