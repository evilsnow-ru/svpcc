package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"svpcc/config"
)

const (
	dataSuffix		= "/data"
	buffersSuffix	= "/buffers"
)

type BufferData struct {
	id string
	fileName string
	bufferSize string
}

func (buf *BufferData) Id() string {
	return buf.id
}

func (buf *BufferData) FileName() string {
	return buf.fileName
}

func (buf *BufferData) BufferSize() string {
	return buf.bufferSize
}

type dataEndpointResult struct {
	Data []dataEndpointEntry `json:"data"`
}

type dataEndpointEntry struct {
	Id string		`json:"id"`
	FileName string `json:"file"`
}

type buffersEndpointResult struct {
	Data []buffersEndpointEntry `json:"data"`
}

type buffersEndpointEntry struct {
	Id string			`json:"id"`
	BufferSize string	`json:"dataSize"`
}

func ReadBuffersData(config config.Config) ([]BufferData, error) {
	dataResult, err := readDataEndpoint(config)

	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]*BufferData)

	for i := range dataResult.Data {
		entry := dataResult.Data[i]
		dataMap[entry.Id] = &BufferData{
			id: entry.Id,
			fileName: entry.FileName,
			bufferSize: "",
		}
	}

	buffersResult, err := readBuffersEndpoint(config)

	if err != nil {
		return nil, err
	}

	for i := range buffersResult.Data {
		entry := buffersResult.Data[i]
		bufferData, present := dataMap[entry.Id]

		if present {
			bufferData.bufferSize = entry.BufferSize
		}
	}

	result := make([]BufferData, 0)

	for key := range dataMap {
		result = append(result, *dataMap[key])
	}

	sort.Slice(result, func (i, j int) bool {
		id1, _ := strconv.Atoi(result[i].Id())
		id2, _ := strconv.Atoi(result[j].Id())
		return id1 < id2
	})

	return result, nil
}

func readBuffersEndpoint(config config.Config) (*buffersEndpointResult, error) {
	fmt.Println("Accessing " + buffersSuffix + " endpoint...")
	endpointUrl := config.ServerAddress() + buffersSuffix
	response, err := http.Get(endpointUrl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error accessing endpoint at \"%s\": %s.", endpointUrl, err.Error()))
	}

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Http result code: %d, expected 200.", response.StatusCode))
	}

	dataBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading content from endpoint \"%s\": %s.", endpointUrl, err.Error()))
	}

	err = response.Body.Close()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error closing connection for endpoint \"%s\": %s", endpointUrl, err.Error()))
	}

	data := string(dataBytes)
	dataResult := buffersEndpointResult{
		Data: []buffersEndpointEntry{},
	}

	err = json.NewDecoder(strings.NewReader(data)).Decode(&dataResult)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error converting result from endpoint \"%s\": %s. Error: %s.", endpointUrl, data, err.Error()))
	}

	log.Printf("Received data items: %d", len(dataResult.Data))
	return &dataResult, nil
}

func readDataEndpoint(config config.Config) (*dataEndpointResult, error) {
	fmt.Println("Accessing " + dataSuffix + " endpoint...")
	endpointUrl := config.ServerAddress() + dataSuffix
	response, err := http.Get(endpointUrl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error accessing endpoint at \"%s\": %s.", endpointUrl, err.Error()))
	}

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Http result code: %d, expected 200.", response.StatusCode))
	}

	dataBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading content from endpoint \"%s\": %s.", endpointUrl, err.Error()))
	}

	err = response.Body.Close()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error closing connection for endpoint \"%s\": %s", endpointUrl, err.Error()))
	}

	data := string(dataBytes)
	dataResult := dataEndpointResult{
		Data: []dataEndpointEntry{},
	}

	err = json.NewDecoder(strings.NewReader(data)).Decode(&dataResult)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error converting result from endpoint \"%s\": %s. Error: %s.", endpointUrl, data, err.Error()))
	}

	log.Printf("Received data items: %d", len(dataResult.Data))
	return &dataResult, nil
}