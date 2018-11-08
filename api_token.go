package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/*用go写的一个接口测试实例，Post方法，修改请求头*/

const (
	authorization        string = "4a152d2f31888e24c608ec643b09d7d0c"
	api_key              string = "b6b8b216-3d3d-48a5-8c6d-df07e5587f43"
	api_secret           string = "ea3aae23145dd8ba39949f2394a01862"
	content_type         string = "application/json"
	accept               string = "application/json"
	method               string = "POST"
	server_api_token_url string = "url"
	server_api_order_url string = "url"
)

/*MD5转换*/
func make_sign(sign_str string) string {
	fmt.Printf("%s\n", sign_str)
	data := []byte(sign_str)
	has := md5.Sum(data)
	md5_string := fmt.Sprintf("%x", has)
	fmt.Printf("%s\n", md5_string)
	fmt.Printf("%s\n", md5_string[0:28])
	return md5_string[0:28]

}

/*发送http_post请求*/
func http_post(data *map[string]interface{}, url string) {

	client := &http.Client{}

	data_json, json_err := json.Marshal(data)

	if json_err != nil {
		fmt.Println(json_err.Error())
		return
	}

	/*构造请求体*/
	request, request_err := http.NewRequest(method, url, strings.NewReader(string(data_json)))
	if request_err != nil {
		fmt.Println(request_err.Error())
		return
	}
	defer request.Body.Close()
	request.Header.Add("Content-Type", content_type)
	request.Header.Add("Accept", accept)
	request.Header.Add("Authorization", authorization)

	/*发起请求*/
	response, client_err := client.Do(request)
	if client_err != nil {
		fmt.Println(request_err.Error())
		return
	}
	defer response.Body.Close()

	body, body_err := ioutil.ReadAll(response.Body)
	if body_err != nil {
		fmt.Println(body_err.Error())
	}

	fmt.Println(string(body))
}

func main() {

	/*测试获取apiToken*/
	/*拼装参数*/
	data := make(map[string]interface{})
	data["api_privilege"] = string("Accounts,Order,Withdraw")
	data["google_auth_code"] = string("")
	data["label"] = string("order")
	data["security_pwd"] = string("a123456")
	data["sms_code"] = string("aaa")
	data["sms_code_email"] = string("")
	data["sms_code_mobile"] = string("")
	data["trusted_ip"] = string("")
	http_post(&data, server_api_token_url)

	/*测试下单*/
	data_order := make(map[string]interface{})
	data_order["api_key"] = string(api_key)
	uid, _ := uuid.NewV4()
	data_order["o_no"] = uid.String()
	data_order["o_price_type"] = string("limit")
	data_order["o_type"] = string("buy")
	data_order["price"] = float32(0.188)
	data_order["sign_type"] = string("MD5")
	data_order["symbol"] = string("BTCCBTC")
	data_order["timestamp"] = time.Now().Unix()
	data_order["volume"] = float32(22)

	buffer := bytes.Buffer{}
	buffer.WriteString("api_key=dfb85f1d-6671-4baf-94cf-445056aa4b22")
	buffer.WriteString(api_key)
	buffer.WriteString("&o_no=1541235652&")
	buffer.WriteString("o_price_type=limit&o_type=buy&price=0.188&sign_type=MD5&symbol=BTCCBTC")
	buffer.WriteString("&timestamp=")
	buffer.WriteString(string(time.Now().Unix()))
	buffer.WriteString("&volume=22&apiSecret=")
	buffer.WriteString(api_secret)
	sign_str := buffer.String()
	data_order["sign"] = string(make_sign(sign_str))

	http_post(&data_order, server_api_order_url)
}
