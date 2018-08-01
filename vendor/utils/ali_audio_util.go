package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	ttsURL   = "http://nlsapi.aliyun.com/speak?"
	akId     = "LTAIK1CAPABIaXad"
	akSecret = "eAkz2skeYzYQQvBP9So6E9cdXezC0U"
)

type TtsRequest struct {
	EncodeType            string `json:"encode_type"`
	VoiceName             string `json:"voice_name"`
	Volume                int    `json:"volume"`
	SampleRate            int    `json:"sample_rate"`
	SpeechRate            int    `json:"speech_rate"`
	PitchRate             int    `json:"pitch_rate"`
	TtsNus                int    `json:"tts_nus"`
	BackgroundMusicId     int    `json:"background_music_id"`
	BackgroundMusicOffset int    `json:"background_music_offset"`
	BackgroundMusicVolume int    `json:"background_music_volume"`
}

func AudioGenerate(text string) error {
	tts := TtsRequest{EncodeType: "wav",
		VoiceName: "xiaoyun",
		Volume: 50,
		SampleRate: 16000,
		SpeechRate: -100,
		PitchRate: 0,
		TtsNus: 1,
	}

	return tts.sendPost(text)
}

func (t *TtsRequest) sendPost(text string) error {
	if IsEmpty(t.EncodeType) {
		t.EncodeType = "wav"
	}
	fmt.Println(text)
	url := t.concatURL()
	fmt.Printf("url:%s\n", url)
	var (
		textMB      = md5base64(text)
		method      = "POST"
		contentType = "text/plain"
		accept      = "audio/" + t.EncodeType + ",application/json"
		temp        = []rune(text)
		length      = len(temp)
		zone, _     = time.LoadLocation("GMT")
		date        = time.Now().In(zone).Format("Mon, 02 Jan 2006 15:04:05 MST")
	)
	stringToSign := method + "\n" + accept + "\n" + textMB + "\n" + contentType + "\n" + date
	signature := hMACSha1(stringToSign, akSecret)
	authHeader := "Dataplus " + akId + ":" + signature

	body := bytes.NewReader([]byte(text))
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Printf("http.NewRequest,[err=%s][url=%s]\n", err, url)
		return err
	}
	request.Header.Set("accept", accept)
	request.Header.Set("content-type", contentType)
	request.Header.Set("date", date)
	request.Header.Set("Authorization", authHeader)
	request.Header.Set("Content-Length", ParseString(length))
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("http.Do failed,[err=%s][url=%s]\n", err, url)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("http.Do failed,[err=%s][url=%s]\n", err, url)
	}
	fmt.Printf("status:%s  code:%d", resp.Status, resp.StatusCode)
	ioutil.WriteFile("./"+GenerateUuid()+"."+t.EncodeType, b, 0666)
	return err
}

func md5base64(s string) string {
	m := md5.New()
	io.WriteString(m, s)
	textMd5 := m.Sum(nil)
	textBase64 := base64.StdEncoding.EncodeToString(textMd5)
	return textBase64
}

func hMACSha1(s string, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	io.WriteString(mac, s)
	textMac := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(textMac)
}

func (t *TtsRequest) concatURL() string {
	var url string
	tMap, e := ParseMap(t)
	if e != nil {
		return EMPTY_STRING
	}
	fmt.Println(tMap)
	for k, v := range tMap {
		if !IsEmpty(v) {
			url += "&" + k + "=" + ParseString(v)
		}
	}
	return ttsURL + SubString(url, 1, -1, EMPTY_STRING)
}
