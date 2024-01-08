package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type GoogleLoginResult struct {
	Access_token  string
	Expires_in    int
	Refresh_token string
	Scope         string
	Token_type    string
}

// "{\n  \
// 	"access_token\": \"ya29.a0AfB_byByxNJtjxa0jMhI7sG3CBo4_6b1VUEeOrJAzyJYrZxvca1cgFSqsTpBp1Ot_KubrHd65FMRpF8gC-lqmLNHbrTW5rMXOfbRorGG7Orgb6cF3_SpWaxNNVdUmGBniBYqXQ0_5gKprAWhDS32sBKSLxKQ25ItdOkiaCgYKAd8SARMSFQHGX2MicMA2BRvZa1LixebSiIPHJw0171\",
// 	\n  \"expires_in\": 3599,\n  \
// 	"refresh_token\": \"1//0e6ejmhYa7iwbCgYIARAAGA4SNwF-L9Ir5AWYn8iMsywJYdmhig35N2UAAeOj2IBp8s5Ahwj2EobiVkzTpiDIgenBJei1mD9FFXM\",\n  \
// 	"scope\": \"https://www.googleapis.com/auth/drive.appdata https://www.googleapis.com/auth/games_lite\",\n
// 	\"token_type\": \"Bearer\"\n}"

func GetGoogleLoginDatFromCode(code string) *GoogleLoginResult {
	_map := map[string]string{}
	_map["client_id"] = Conf.String("client_id")
	_map["client_secret"] = Conf.String("client_secret")
	_map["redirect_uri"] = Conf.String("redirect_uri")
	_map["grant_type"] = "authorization_code"
	_map["code"] = code
	return googleDatFromCode(_map)
}

func googleDatFromCode(param map[string]string) *GoogleLoginResult {
	LogDebug("google request", param)
	bytesData, err := json.Marshal(param)
	if ChkErr(err) {
		return nil
	}
	req, err := http.NewRequest("POST", Conf.String("check_url_login"), bytes.NewBuffer(bytesData))
	if ChkErr(err) {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := httpClient.Do(req)
	if ChkErr(err) {
		return nil
	}

	respBytes, err := io.ReadAll(response.Body)
	if ChkErr(err) {
		return nil
	}
	LogDebug("google response", string(respBytes))
	result := &GoogleLoginResult{}
	err = json.Unmarshal(respBytes, result)
	if ChkErr(err) {
		return nil
	}

	if result != nil && result.Access_token == "" {
		return nil
	}
	return result
}

var requestCodeUrl = Conf.String("request_code_url_pay")
var code = Conf.String("code_pay")
var PayAccessToken = ""
var PayRefreshToken = Conf.String("refresh_token_pay")
var PayScope = Conf.String("scope_pay")

func getGooglePayAccessToken() {
	_url := fmt.Sprintf(requestCodeUrl, Conf.String("redirect_uri_pay"), Conf.String("client_id_pay"))
	Log("=====请求Google Pay Code 地址,网页打开=======", _url)

	if PayRefreshToken == "" {
		_map := map[string]string{}
		_map["client_id"] = Conf.String("client_id_pay")
		_map["client_secret"] = Conf.String("client_secret_pay")
		_map["redirect_uri"] = Conf.String("redirect_uri_pay")
		_map["grant_type"] = "authorization_code"
		_code, _ := url.QueryUnescape(code)
		_map["code"] = _code

		_result := googleDatFromCode(_map)

		if _result == nil {
			Log("请求Pay Toke Fail !!!!!")
		} else {
			PayAccessToken = _result.Access_token
			PayRefreshToken = _result.Refresh_token
			PayScope = _result.Scope

			Log("PayAccessToken", PayAccessToken)
			Log("PayRefreshToken", PayRefreshToken)
			Log("PayScope", PayScope)
		}
	} else if !ServerInChina {
		_map := map[string]string{}
		_map["client_id"] = Conf.String("client_id_pay")
		_map["client_secret"] = Conf.String("client_secret_pay")
		_map["refresh_token"] = PayRefreshToken
		_map["grant_type"] = "refresh_token"

		_result := googleDatFromCode(_map)
		if _result == nil {
			Log("refresh token error")
		} else {
			PayAccessToken = _result.Access_token
		}
		time.Sleep(time.Minute * 60)
		getGooglePayAccessToken()
	}
}
