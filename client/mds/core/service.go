package core

import (
	"encoding/json"
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/utils"
	"net/http"
	"strings"
)

type Service struct {
	Endpoint string
	Api      *Root
}

type Root struct {
	HostUrl    *string
	OrgId      string
	AuthToUse  *model.ClientAuth
	HttpClient *http.Client
	Token      *string
}

func NewService(hostUrl *string, endPoint string, root *Root) *Service {
	return &Service{
		Endpoint: fmt.Sprintf("%s/api/%s", *hostUrl, endPoint),
		Api:      root,
	}
}

func (r *Root) Get(url *string, queryModel interface{}, dest interface{}) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		return nil, err
	}

	if queryModel != nil {
		pairs, err := utils.ToKeyValuePairs(queryModel)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = utils.ProcessAsQuery(req.URL.Query(), &pairs).Encode()
	}

	body, err := r.doRequest(req)
	if err != nil {
		return nil, err
	}

	if dest != nil {
		err = json.Unmarshal(body, &dest)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

func (r *Root) Post(url *string, reqBody interface{}, dest interface{}) ([]byte, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, *url, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := r.doRequest(req)
	if err != nil {
		return nil, err
	}

	if dest != nil && body != nil {
		if len(body) <= 0 {
			err = fmt.Errorf("error occured.")
			return nil, err
		}
		err = json.Unmarshal(body, &dest)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

func (r *Root) Delete(url *string, reqBody interface{}, dest interface{}) ([]byte, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, *url, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := r.doRequest(req)
	if err != nil {
		return nil, err
	}

	if dest != nil {
		err = json.Unmarshal(body, &dest)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

func (r *Root) Patch(url *string, reqBody interface{}, dest interface{}) ([]byte, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	fmt.Printf("BODY: %s", reqBody)
	req, err := http.NewRequest(http.MethodPatch, *url, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := r.doRequest(req)
	if err != nil {
		return nil, err
	}

	if dest != nil {
		err = json.Unmarshal(body, &dest)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

func (r *Root) Put(url *string, reqBody interface{}, dest interface{}) ([]byte, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	fmt.Printf("BODY: %s", rb)
	req, err := http.NewRequest(http.MethodPut, *url, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := r.doRequest(req)
	if err != nil {
		return nil, err
	}

	if dest != nil {
		err = json.Unmarshal(body, &dest)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}
