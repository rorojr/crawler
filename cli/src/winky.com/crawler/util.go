package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"sort"
	"time"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

/**
 * 四舍五入
 *
 */
func round(f float64, n int) float64 {
	pow := math.Pow10(n)
	return math.Trunc((f+0.5/pow)*pow) / pow
}

/**
 * 生成随机数
 *
 */
func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

func getImgUrl()string{
	rd := RandInt64(1, 4)
	return imgUrl[rd]
}

/**
 * 字符串时间格式化成时间戳
 *
 * return int
 */
func formatTime(dateStr string) (s int64, err error) {
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.Local)
	s = timestamp.Unix()
	return
}

/**
 * 返回毫秒
 *
 * return int
 */
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

/**
 * key 排序
 *
 * return int
 */
func sortStringKey(list map[string]string) (keys []string) {
	for k := range list {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

func httpSend(apiUrl string, params map[string]string, method string) []byte {

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
	quest, err := http.NewRequest(method, apiUrl, body)
	if err != nil {
		log.Println("http new request failed")
	}
	quest.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response, err := client.Do(quest)
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("get response %s failed: %v \n", response, err)
	}

	if response.StatusCode != 200 {
		log.Println("http antnest failed", apiUrl, response.StatusCode)
	}
	return data
}

func trace(str string) {
	if conf.Log.Debug {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(now + " " + str)
	}
}

func subStr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])
}

func byte2string(in [16]byte) []byte {
	return in[:16]
}