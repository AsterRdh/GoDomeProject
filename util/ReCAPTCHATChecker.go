package util

import (
	"awesomeProject/model"
	"fmt"
	"net/http"
)

func CheckReCAPTCHAToken(token string) (bool, error) {
	ReCAPTCHATokenKey := model.ReCaptchaTokenKey
	targetUrl := model.ReCaptchaURL + "?secret=" + ReCAPTCHATokenKey + "&response=" + token
	response, err := http.Post(targetUrl, "", nil)
	if err != nil {
		return false, err
	}
	fmt.Println(response)
	return true, err
}
