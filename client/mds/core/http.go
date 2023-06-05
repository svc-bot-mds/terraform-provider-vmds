package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	headerAuth        = "csp-auth-token"
	headerContentType = "content-type"
	headerTokenType   = "token-type"
	contentTypeJSON   = "application/json"
)

var headers = map[string]string{
	headerContentType: contentTypeJSON,
}

type HttpError struct {
	error
	StatusCode int
}

type ApiError struct {
	HttpError
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMsg"`
}

func (h HttpError) Error() string {
	return h.error.Error()
}

func (h HttpError) Unwrap() error {
	return h.error
}

func (r *Root) doRequest(req *http.Request) ([]byte, error) {
	r.addHeaders(req)

	res, err := r.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		errorWithMsg := fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
		var apiError ApiError
		if err = json.Unmarshal(body, &apiError); err != nil {
			return nil, errors.Join(errorWithMsg)
		}
		apiError.error = errorWithMsg
		apiError.StatusCode = res.StatusCode
		return nil, apiError
	}

	return body, nil
}

func (r *Root) addHeaders(req *http.Request) {
	for header, value := range headers {
		req.Header.Add(header, value)
	}
	if r.Token != nil {
		req.Header.Set(headerAuth, " "+*r.Token)
		//	TODO: add token-type
	}
}

func closeBody(body io.Closer) {
	if err := body.Close(); err != nil {
		fmt.Println(err)
	}
}
