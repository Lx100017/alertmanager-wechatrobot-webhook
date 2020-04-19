package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jacobslei/alertmanager-wechatrobot-webhook/model"
	"github.com/jacobslei/alertmanager-wechatrobot-webhook/transformer"
)

var AccessToken = ""

//
// @Title Send send markdown message to enterprise WeCHAT by Robot
// @Author Jacobs Lei 2020-04-19
func Send(notification model.Notification, defaultRobot string) (err error) {

	markdown, robotURL, err := transformer.TransformToMarkdown(notification)

	if err != nil {
		return
	}

	data, err := json.Marshal(markdown)
	if err != nil {
		return
	}

	var wechatRobotURL string

	if robotURL != "" {
		wechatRobotURL = robotURL
	} else {
		wechatRobotURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + defaultRobot
	}

	req, err := http.NewRequest(
		"POST",
		wechatRobotURL,
		bytes.NewBuffer(data))

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	return
}

//
// @Title send Text message to WeChat Application
// @Author Jacobs Lei 2020-04-20
func SendApp(notification model.Notification, ToTag string, CorpId string, AgentId string, AppSecret string) (err error) {
	if AccessToken == "" {
		AccessTokenResponse, err := getAccessToken(CorpId, AppSecret)
		if err != nil {
			return err
		}
		var bytes, _ = json.MarshalIndent(AccessTokenResponse, " ", " ")

		fmt.Println("AccessTokenResponse  is :", string(bytes))

		AccessToken = AccessTokenResponse.AccessToken

	}

	fmt.Println("Access token is :", AccessToken)

	markdown, err := transformer.TransformToAppMessageText(notification, ToTag, AgentId)

	if err != nil {
		return err
	}

	wechatCommonResponse, err := sendMessageOriginal(markdown)
	if err != nil {
		return err
	}

	if wechatCommonResponse.ErrCode == 40014 {
		AccessTokenResponse, err := getAccessToken(CorpId, AppSecret)
		if err != nil {
			return err
		}
		AccessToken = AccessTokenResponse.AccessToken
		sendMessageOriginal(markdown)
	}
	return nil
}

func getAccessToken(CorpId string, AppSecret string) (accessToken *model.AccessTokenResponse, err error) {
	var GetAccessTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + CorpId + "&corpsecret=" + AppSecret
	//fmt.Println("GetAccessTokenUrl :", GetAccessTokenUrl)
	var AccessTokenResponseB []byte
	if response, err := http.Get(GetAccessTokenUrl); err != nil {
		return nil, err
	} else {
		if AccessTokenResponseB, err = ioutil.ReadAll(response.Body); err != nil {
			return nil, err
		}
		fmt.Println(string(AccessTokenResponseB))
	}
	var AccessTokenResponse = &model.AccessTokenResponse{}
	if err := json.Unmarshal(AccessTokenResponseB, AccessTokenResponse); err != nil {
		return nil, err
	}
	//var bytes, _ = json.MarshalIndent(AccessTokenResponse, " ", " ")
	//fmt.Println("AccessTokenResponse  is :",string(bytes) )
	return AccessTokenResponse, nil
}

func sendMessageOriginal(text *model.WeChatText) (wechatCommonReponse *model.WechatMessageCommonResponse, err error) {
	data, err := json.Marshal(text)
	if err != nil {
		return nil, err
	}

	var sendMessageURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + AccessToken

	req, err := http.NewRequest(
		"POST",
		sendMessageURL,
		bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var respByte []byte
	if respByte, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	wechatCommonReponse = &model.WechatMessageCommonResponse{}
	err = json.Unmarshal(respByte, wechatCommonReponse)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Println(fmt.Println("response Status:", resp.Status))
	log.Println(fmt.Println("response Headers:", resp.Header))
	return wechatCommonReponse, nil

}
