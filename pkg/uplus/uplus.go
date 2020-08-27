package uplus

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/kimtaek/gamora/pkg/helper"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type config struct {
	Mode               string `env:"APP_MODE" envDefault:"debug"`
	AccessKey          string `env:"UPLUS_ACCESS_KEY" envDefault:""`
	SecretKey          string `env:"UPLUS_SECRET_KEY" envDefault:""`
	UniversalTelephone string `env:"UNIVERSAL_TELEPHONE" envDefault:""`
	SendFrom           string `env:"SEND_FROM" envDefault:"13332"`
}

type responseData struct {
	GrpId string `json:"grp_id"`
	MsgId string `json:"msg_id"`
}

type responseJson struct {
	Data  responseData
	RDesc string `json:"rdesc"`
	RCode string `json:"rcode"`
}

var c config

func Setup() {
	_ = env.Parse(&c)
}

func SendSMS(phoneNumber string, message string) error {
	if c.UniversalTelephone == "82-1000000000" {
		return nil
	}

	if c.UniversalTelephone != "" {
		phoneNumber = c.UniversalTelephone
	}

	url := "https://openapi.sms.uplus.co.kr:4443/v1/send"
	uuid := fmt.Sprintf("%08d", rand.Intn(100000000))

	timestamp := strconv.FormatInt(time.Now().Unix(), 10) + "000"
	koreanPhoneNumber := "0" + strings.Split(phoneNumber, "-")[1]

	requestBody, err := json.Marshal(map[string]string{
		"send_type": "S",
		"msg_type":  "S",
		"to":        koreanPhoneNumber,
		"from":      c.SendFrom,
		"msg":       message,
		"country":   "82",
		"subject":   "",
		"device_id": "",
		"datetime":  "",
	})

	if err != nil {
		helper.Error(err.Error())
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("api_key", c.AccessKey)
	request.Header.Set("algorithm", "1")
	request.Header.Set("hash_hmac", hash(timestamp, uuid))
	request.Header.Set("cret_txt", uuid)
	request.Header.Set("timestamp", timestamp)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		helper.Error(err.Error())
		return err
	}

	var data responseJson
	err = json.NewDecoder(response.Body).Decode(&data)
	defer response.Body.Close()

	helper.Info(koreanPhoneNumber)
	helper.Info(data.RDesc)

	if err != nil {
		return err
	}

	if data.RCode == "1000" {
		return nil
	}
	return errors.New("error")
}

func hash(timestamp string, uuid string) string {
	h := sha1.New()
	h.Write([]byte(c.AccessKey + timestamp + uuid + c.SecretKey))
	return hex.EncodeToString(h.Sum(nil))
}
