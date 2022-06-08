package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// Deprecated: move to package download
	Authlibversion = "1.1.41"
	// Deprecated: move to package download
	Authlibsha1 = "64e96fc1e29e312a4d1e9d530f492cfa4b089cf1"
)

// Deprecated: move to package download
var Authliburls = []string{
	"https://download.mcbbs.net/mirrors/authlib-injector/artifact/41/authlib-injector-1.1.41.jar",
	"https://authlib-injector.yushi.moe/artifact/41/authlib-injector-1.1.41.jar",
	"https://bmclapi2.bangbang93.com/mirrors/authlib-injector/artifact/41/authlib-injector-1.1.41.jar",
}

func Getauthlibapi(api string) (apiaddress string, err error) {
	u, err := url.Parse(api)
	if err != nil {
		return "", fmt.Errorf("Getauthlibapi: %w", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
		api = u.String()
	}
	c := &http.Client{
		Timeout:   5 * time.Second,
		Transport: Transport,
	}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return "", fmt.Errorf("Getauthlibapi: %w", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	reps, err := c.Do(req)
	if reps != nil {
		defer reps.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("Getauthlibapi: %w", err)
	}
	defer func() {
		err = checkapi(apiaddress)
	}()
	header := reps.Header.Get("X-Authlib-Injector-API-Location")
	if header == "" {
		return api, nil
	}
	hurl, err := url.Parse(header)
	if err != nil {
		return "", fmt.Errorf("Getauthlibapi: %w", err)
	}
	if hurl.Scheme == "" {
		api = "https://" + reps.Request.URL.Host + "/" + header
	} else {
		return header, nil
	}
	return api, nil
}

type yggdrasil struct {
	SignaturePublickey string `json:"signaturePublickey"`
}

func checkapi(url string) error {
	c := &http.Client{
		Timeout:   5 * time.Second,
		Transport: Transport,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("checkapi: %w", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	reps, err := c.Do(req)
	if reps != nil {
		defer reps.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("checkapi: %w", err)
	}
	b, err := ioutil.ReadAll(reps.Body)
	if err != nil {
		return fmt.Errorf("checkapi: %w", err)
	}
	var y yggdrasil
	err = json.Unmarshal(b, &y)
	if err != nil {
		return fmt.Errorf("checkapi: %w", err)
	}
	if y.SignaturePublickey == "" {
		return fmt.Errorf("checkapi: %w", JsonNotTrue)
	}
	return nil
}

var JsonNotTrue = errors.New("json not true")
