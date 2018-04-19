package main

import (
	"encoding/base64"
	"fmt"
	"github.com/panthesingh/goson"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"winky.com/x/des"
)

/*
 * 蚁巢push推送
 * @param userId int 用户ID
 * @param mobile string 手机号码
 * @param jPushType string 推送类型
 * @param jPushId string jPush_id
 * @param jPushData string 客户端解析数据
 * @param jsPushContent string 推送内容
 * @return bool
 */
func antUploadImage(fileSuffix string, fileContent []byte, fileMd5 string) (string, error) {

	path := fmt.Sprintf("news/%s", time.Now().Format("2006-01-02"))
	path = path + "/" + subStr(fileMd5, 0, 4) + "/" + subStr(fileMd5, 4, 8)

	apiUrl := conf.Ant.AntNestUrl + "/transfer/oss/contentupload"
	params := make(map[string]string)

	params["app_type"] = "1"
	params["source_path"] = path
	params["file_suffix"] = fileSuffix
	params["file_content"] = base64.StdEncoding.EncodeToString(fileContent)
	params["code_md5"] = fileMd5

	result := antSend(apiUrl, params)
	jsonObj, err := goson.Parse(result)
	if err != nil {
		return "", fmt.Errorf("ParseJson json %s failed: %v", result, err.Error())
	}
	code := jsonObj.Get("code").String()
	ossPath := jsonObj.Get("data").Get("path").String()
	if code != "000000" {
		return "", fmt.Errorf("ParseJson antnest upload images code = %s msg = %s", code, jsonObj.Get("msg").String())
	}
	return ossPath, nil
}

/*
 * 蚁巢request
 * @param apiUrl string 手机号码
 * @param params map 推送类型
 * @return byte
 */
func antSend(apiUrl string, params map[string]string) []byte {

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, address string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(network, address, 5*time.Second)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	value := url.Values{}
	for k, v := range params {
		value.Set(k, v)
	}

	body := ioutil.NopCloser(strings.NewReader(value.Encode()))
	quest, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		log.Println("http antnest failed")
	}

	key := fmt.Sprintf("%s_%s_%d", conf.Ant.AppKey, conf.Ant.SecretKey, time.Now().Unix())
	token := des.Base64des(key, conf.Ant.SecretKey)

	quest.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	quest.Header.Set("appkey", conf.Ant.AppKey)
	quest.Header.Set("token", token)

	response, err := client.Do(quest)
	if err != nil {
		log.Printf("antSend client request failed: %v \n", err)
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("antSend ReadAll %s failed: %v \n", data, err)
	}

	if response.StatusCode == 200 {
		return data
	} else {
		log.Println("http antnest failed", apiUrl, response.StatusCode)
	}
	return data
}
