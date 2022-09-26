package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/glide-im/glide/pkg/logger"
	"io"
	"net/http"
)

var baseUrl = "http://localhost:8081/api/"

func SetBaseUrl(url string) {
	baseUrl = url
}

func postJson(path string, data interface{}, resp interface{}) error {

	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}

	logger.D("<post> %s %s", path, string(bts))

	respData, err := http.Post(baseUrl+path, "application/json", bytes.NewBuffer(bts))

	if err != nil {
		return err
	}
	if respData.StatusCode != http.StatusOK {
		return errors.New(respData.Status)
	}

	bt, err := io.ReadAll(respData.Body)
	if err != nil {
		return err
	}

	logger.D("<resp> %s %s", path, string(bt))

	response := &CommonResponse{}
	if err = json.Unmarshal(bt, response); err != nil {
		return err
	}

	if response.Code != 100 {
		return &ApiError{
			Code:    response.Code,
			Message: response.Msg,
		}
	}

	if resp == nil {
		return nil
	}
	if err = json.Unmarshal(response.Data, resp); err != nil {
		return err
	}
	return nil
}

func get(path string) (response string, error error) {
	resp, err := http.Get(baseUrl + path)
	if err != nil {
		return "", err
	}
	bt, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bt), nil
}
