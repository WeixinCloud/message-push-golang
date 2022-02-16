package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const weixinAPI = "http://api.weixin.qq.com/cgi-bin/message/custom/send"

// IndexHandler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := getIndex()
		if err != nil {
			fmt.Fprintf(w, "internel error: %+v", err)
			fmt.Println(err)
			return
		}
		fmt.Fprint(w, data)
		return
	} else if r.Method == http.MethodPost {
		// 没有x-wx-source头的，不是微信的来源，不处理
		sourceKey := http.CanonicalHeaderKey("x-wx-source")
		source := r.Header[sourceKey]
		if len(source) <= 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid request source")
			return
		}

		openIDKey := http.CanonicalHeaderKey("x-wx-openid")
		openID := r.Header[openIDKey]
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "internel error: %+v", err)
			fmt.Println(err)
			return
		}

		if len(openID) <= 0 {
			fmt.Fprint(w, "x-wx-openid is null")
			return
		}

		rsp, err := sendMsg(openID[0], string(body))
		if err != nil {
			fmt.Fprintf(w, "internel error: %+v", err)
			fmt.Println(err)
		}
		fmt.Fprintf(w, "sendMsg rsp: %s", rsp)
	} else {
		fmt.Fprintf(w, "method %s not allow", r.Method)
	}
}

func sendMsg(openID, msgText string) (string, error) {
	text := map[string]string{"content": fmt.Sprintf("云托管接收消息推送成功，内容如下: %s\n", msgText)}
	req := struct {
		Touser  string            `json:"touser"`
		Msgtype string            `json:"msgtype"`
		Text    map[string]string `json:"text"`
	}{Touser: openID, Msgtype: "text", Text: text}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	fmt.Printf("req :[%s]\n", string(reqBytes))

	resp, err := http.Post(weixinAPI, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("resp :[%s]\n", string(bodyContent))

	return string(bodyContent), nil
}

// getIndex 获取主页
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./index.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
