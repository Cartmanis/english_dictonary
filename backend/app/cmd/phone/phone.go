package phone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SmsRes struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	GuideLink string `json:"guide_link"`
}

func SendSms(numberPhone, message string) error {
	phone := getCutPhone(numberPhone)
	fmt.Println(phone)
	url := getUrlString()
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	smsRes := &SmsRes{}
	if err := json.Unmarshal(data, smsRes); err != nil {
		return err
	}
	fmt.Println(smsRes.Message)
	fmt.Println(res.Status)
	return nil
}

func getCutPhone(phone string) string {
	p := strings.Replace(phone, "-", "", -1)
	return strings.Replace(p, " ", "", -1)
}

func getUrlString() string {
	return "http://192.168.0.17:8090/SendSMS?username=cartmanis&password=5eu7ve&phone=23231&message=777"
}
