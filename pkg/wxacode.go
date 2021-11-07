package wxacode

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WXACode struct {
	appID     string
	appSecret string
}

type wxAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type QrCodeResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	B64Image string `json:"buffer,omitempty"`
}

func NewWxCodeClient(appID, appSecret string) *WXACode {
	return &WXACode{
		appID:     appID,
		appSecret: appSecret,
	}
}

func (p *WXACode) GetAccessToken(ctx context.Context) (*wxAccessTokenResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		p.appID, p.appSecret)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data wxAccessTokenResponse
	if err = json.Unmarshal(rawBody, &data); err == nil {
		return &data, nil
	}
	return nil, err
}

func (p *WXACode) GenerateQrCode(ctx context.Context, accessToken string, scene string) (*QrCodeResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s", accessToken)
	body := map[string]string{}
	if len(scene) > 0 {
		body["scene"] = scene
	}
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		var rawBody []byte
		rawBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var data QrCodeResponse
		if err = json.Unmarshal(rawBody, &data); err == nil {
			return &data, nil
		} else {
			base64Encoding := "data:image/jpeg;base64,"
			b64Image := base64.StdEncoding.EncodeToString(rawBody)
			return &QrCodeResponse{
				ErrCode:  0,
				ErrMsg:   "",
				B64Image: base64Encoding + b64Image,
			}, nil
		}
	}
	return nil, fmt.Errorf("generate failed with wx http status code: %d", resp.StatusCode)
}
