package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"svpcc/config"
)

const flushEndpointSuffix = "/flush"


func Flush(cfg config.Config, id int) (string, error) {

	var endpointUrl string

	if id >= 0 {
		endpointUrl = cfg.ServerAddress() + flushEndpointSuffix + "/" + strconv.Itoa(id)
	} else {
		endpointUrl = cfg.ServerAddress() + flushEndpointSuffix
	}

	fmt.Println(fmt.Sprintf("Accessing endpoint url: \"%s\".", endpointUrl))
	response, err := http.Get(endpointUrl)

	if err != nil {
		return "", err
	}

	if response == nil || response.StatusCode != 200 || response.Body == nil {
		errorMsg := fmt.Sprintf("Error getting response from \"%s\". Status code: %d.", endpointUrl, response.StatusCode)
		return "", errors.New(errorMsg)
	}

	dataBytes, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()

	if err != nil {
		return "", err
	}

	data := string(dataBytes)
	responseObject := statusResponse{}
	err = json.NewDecoder(strings.NewReader(data)).Decode(&responseObject)

	if  err != nil {
		return "", err
	}

	return responseObject.Status, nil
}