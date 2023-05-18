package goma

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aead/cmac"
	"github.com/joho/godotenv"
)

// https://doc.candyhouse.co/ja/SesameAPI#%E3%82%B3%E3%83%9E%E3%83%B3%E3%83%89%E3%82%B3%E3%83%BC%E3%83%89%E3%83%AA%E3%82%B9%E3%83%88
var cmdNumber = map[string]uint{"toggle": 88, "lock": 82, "unlock": 83}

type Options struct {
	Account       string
	HistoryPage   uint
	HistoryNumber uint
}

type RequestBody struct {
	Cmd     uint   `json:"cmd"`
	History string `json:"history"`
	Sign    string `json:"sign"`
}

func State() string {
	return exec("state", nil)
}

func History(page uint, number uint) string {
	return exec("history", &Options{
		HistoryPage:   page,
		HistoryNumber: number,
	})
}

func Toggle(account string) string {
	return exec("toggle", &Options{
		Account: account,
	})
}

func Lock(account string) string {
	return exec("lock", &Options{
		Account: account,
	})
}

func Unlock(account string) string {
	return exec("unlock", &Options{
		Account: account,
	})
}

func exec(cmd string, opt *Options) string {
	if opt == nil {
		opt = &Options{}
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api_key, ok := os.LookupEnv("SESAME_API_KEY")
	if !ok {
		log.Fatal("Not Fonund SESAME_API_EY in .env")
	}
	base_url, ok := os.LookupEnv("SESAME_BASE_URL")
	if !ok {
		base_url = "https://app.candyhouse.co/api/sesame2/"
	}
	secret_key, ok := os.LookupEnv("SESAME_SECRET_KEY")
	if !ok {
		log.Fatal("Not Fonund SESAME_SECRET_KEY in .env")
	}
	uuid, ok := os.LookupEnv("SESAME_UUID")
	if !ok {
		log.Fatal("Not Fonund SESAME_UUID in .env")
	}
	url := fmt.Sprintf("%s%s", base_url, uuid)

	var req *http.Request
	var res *http.Response
	var method = "GET"

	if cmd == "state" {
		if req, err = http.NewRequest(method, url, nil); err != nil {
			log.Fatal(err)
		}
	} else if cmd == "history" {
		url = fmt.Sprintf("%s/history?page=%d&lg=%d", url, opt.HistoryPage, opt.HistoryNumber)
		if req, err = http.NewRequest(method, url, nil); err != nil {
			log.Fatal(err)
		}
	} else if cmd == "toggle" || cmd == "lock" || cmd == "unlock" {
		url = url + "/cmd"
		var sign = getEncryptKey(secret_key)
		var src = []byte(opt.Account)
		var history = base64.StdEncoding.EncodeToString(src)

		requestBody := RequestBody{
			Cmd:     cmdNumber[cmd],
			History: history,
			Sign:    sign,
		}
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatal(err)
		}
		method = "POST"
		if req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody)); err != nil {
			log.Fatal(err)
		}
	}

	req.Header.Set("x-api-key", api_key)
	var client = &http.Client{}
	if res, err = client.Do(req); err != nil {
		log.Fatal(err)
	} else if res.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("Cant access : HTTP Status %d : %s %s", res.StatusCode, method, url))
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return string(body)
}

func getEncryptKey(key string) string {
	_key, _ := hex.DecodeString(key)
	timestamp := time.Now().Unix()
	timestampBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestampBytes, uint64(timestamp))
	timestampBytes = timestampBytes[1:4]
	block, _ := aes.NewCipher(_key)
	cmac, _ := cmac.New(block)
	cmac.Write(timestampBytes)
	tag := cmac.Sum(nil)
	return hex.EncodeToString(tag)
}
