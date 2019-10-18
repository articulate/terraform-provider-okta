package sdk

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

type MappingProperties struct {
	Properties struct {
		Attribute struct {
			Expression string `json:"expression"`
			PushStatus string `json:"pushStatus"`
		} `json:"attribute"`
	} `json:"properties"`
}

type Mapping []struct {
	ID     string `json:"id"`
	Source struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Schema struct {
				Href string `json:"href"`
			} `json:"schema"`
		} `json:"_links"`
	} `json:"source"`
	Target struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Schema struct {
				Href string `json:"href"`
			} `json:"schema"`
		} `json:"_links"`
	} `json:"target"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

func (m *ApiSupplement) RemovePropertyMapping(mappingId, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s/", mappingId)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(req, nil)
}

func (m *ApiSupplement) ListProfileMappings(mappingId string) ([]*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings", nil)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var mapping []*Mapping
	resp, err := m.RequestExecutor.Do(req, &mapping)
	return mapping, resp, err
}

func (m *ApiSupplement) GetSingleProfileMapping(mappingId string) ([]*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s", mappingId)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var mapping []*Mapping
	resp, err := m.RequestExecutor.Do(req, &mapping)
	return mapping, resp, err
}

func (m *ApiSupplement) AddPropertyMapping(mappingId string, body Mapping, qp *query.Params) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s/", mappingId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	mapping := body
	resp, err := m.RequestExecutor.Do(req, &mapping)
	return &mapping, resp, err
}

func (m *ApiSupplement) UpdateMapping(mappingId string, body Mapping, qp *query.Params) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s/", mappingId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	mapping := body
	resp, err := m.RequestExecutor.Do(req, &mapping)
	if err != nil {
		return nil, resp, err
	}
	return &mapping, resp, nil
}
