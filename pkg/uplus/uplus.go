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

type Configure struct {
	Mode               string `env:"APP_MODE" envDefault:"debug"`
	AccessKey          string `env:"UPLUS_ACCESS_KEY" envDefault:""`
	SecretKey          string `env:"UPLUS_SECRET_KEY" envDefault:""`
	UniversalTelephone string `env:"UNIVERSAL_TELEPHONE" envDefault:""`
	SendFrom           string `env:"SEND_FROM" envDefault:"13332"`
}

type responseJson struct {
	Data struct {
		GrpId string `json:"grp_id"`
		MsgId string `json:"msg_id"`
	} `json:"data"`
	RDesc string `json:"rdesc"`
	RCode string `json:"rcode"`
}

var Config Configure

func Setup() {
	_ = env.Parse(&Config)
}

func SendSMS(phoneNumber string, message string) error {
	if Config.UniversalTelephone == "82-1000000000" {
		return nil
	}

	if Config.UniversalTelephone != "" {
		phoneNumber = Config.UniversalTelephone
	}

	url := "https://openapi.sms.uplus.co.kr:4443/v1/send"
	uuid := fmt.Sprintf("%08d", rand.Intn(100000000))
	phoneSplit := strings.Split(phoneNumber, "-")

	if phoneSplit[0] == "82" {
		phoneNumber = "0" + phoneSplit[1]
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10) + "000"
	requestBody, err := json.Marshal(map[string]string{
		"send_type": "S",
		"msg_type":  "S",
		"to":        phoneNumber,
		"from":      Config.SendFrom,
		"msg":       message,
		"country":   phoneSplit[0],
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
	request.Header.Set("api_key", Config.AccessKey)
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

	helper.Info(phoneNumber)
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
	h.Write([]byte(Config.AccessKey + timestamp + uuid + Config.SecretKey))
	return hex.EncodeToString(h.Sum(nil))
}
