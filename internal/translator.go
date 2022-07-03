package internal

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type translator struct {
	url        string
	sourceLang string
	targetLang string
}

type apiRequest struct {
	Q       string `json:"q"`
	Source  string `json:"source"`
	Target  string `json:"target"`
	Format  string `json:"format"`
	Api_key string `json:"api_key"`
}

type apiResponse struct {
	TranslatedText string `json:"translatedText"`
}

func newTranslator() translator {
	return translator{url: "https://libretranslate.de/translate", sourceLang: "de", targetLang: "en"}
}

func (t *translator) translate(text string) (string, error) {
	apiRequest := apiRequest{Q: text, Source: t.sourceLang, Target: t.targetLang, Format: "text", Api_key: ""}
	requestBytes, err := json.Marshal(apiRequest)
	if err != nil {
		return "", err
	}
	response, err := http.Post(t.url, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	apiResponse := apiResponse{}
	json.Unmarshal(responseBytes, &apiResponse)
	return apiResponse.TranslatedText, nil
}
