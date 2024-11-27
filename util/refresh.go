package util

import (
  "net/http"
  "errors"
)

type InboundPayload struct {
  Data map[string]interface{} `json:"data"`
  Endpoint string `json:"endpoint"`
  AuthServer string `json:"authserver"`
  RefreshToken string `json:"refreshtoken"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	APIDomain   string `json:"api_domain"`
	TokenType   string `json:"token_type"`
}

func RefreshAuthToken(r *http.Request, p *InboundPayload) (*http.Response, error) {
  if p.Endpoint == "" || p.AuthServer == "" || p.RefreshToken == "" {
    return nil, errors.New("Missing endpoint, authserver, or refreshtoken")
  }

  refreshReq, err := http.NewRequest(r.Method, p.AuthServer, nil)  
  if err != nil {
    return nil, err
  }

  refreshResp, err := http.DefaultClient.Do(refreshReq)
  if err != nil {
    return refreshResp, err
  }

  return refreshResp, nil
}
