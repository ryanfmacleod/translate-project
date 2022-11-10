package cli

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateURL = "https://translate.google.com/translate_a/single?"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", translateURL, nil)
	query := req.URL.Query()

	query.Add("client", "gtx")

	query.Add("s1", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	if err != nil {
		log.Fatal("1: There was a problem!")
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("2: There was a problem!")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Something")
		}
	}(res.Body)

	if res.StatusCode == http.StatusTooManyRequests {
		str <- "Too many request"
		wg.Done()
		return
	}

	parsedJson, err := gabs.ParseJSONBuffer(res.Body)

	if err != nil {
		log.Fatalf("3: There was a problem! %s", err)
	}

	firstEl, err := parsedJson.ArrayElement(0)

	if err != nil {
		log.Fatal("4: There was a problem")
	}

	secEl, err := firstEl.ArrayElement(0)

	if err != nil {
		log.Fatal("5: There was a problem")
	}

	translatedString, err := secEl.ArrayElement(0)

	if err != nil {
		log.Fatal("6: There was a problem")
	}

	str <- translatedString.Data().(string)
	wg.Done()
}
