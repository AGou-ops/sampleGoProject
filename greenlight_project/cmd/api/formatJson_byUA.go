package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

// > User-Agent: curl/7.80.0

func writeJsonByUA(req *http.Request, data interface{}) (string, error) {
	client_ua := req.UserAgent()
	result, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", errors.New("indent failed")
	}
	isCurlUA, _ := regexp.MatchString("curl.*", client_ua)

	if isCurlUA {
		return string(result), nil
	}

	return "", errors.New("not curl UA")
}
