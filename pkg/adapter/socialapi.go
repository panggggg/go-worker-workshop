package adapter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type ISocialAPIAdapter interface {
	GetThreads(keyword string) string
	GetAccountInfo(userid string) (string, error)
}

func GetThreads(keyword string) string {
	url := fmt.Sprintf("https://go-workshop-2zcpzmfnyq-de.a.run.app/thread/?hashtag=%s&page_size=5", keyword)
	response, err := http.Get(url)
	failOnError(err, "Failed to get timeline")

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func GetAccountInfo(userid string) (string, error) {
	url := fmt.Sprintf("https://go-workshop-2zcpzmfnyq-de.a.run.app/account/%s", userid)
	response, err := http.Get(url)
	failOnError(err, "Failed to get user info")

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body), nil
}
