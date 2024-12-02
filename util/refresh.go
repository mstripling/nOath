package util

import (
    "net/http"
    "errors"
    "io"
    "encoding/json"
    "fmt"
    "bytes"
)

type InboundPayload struct {
    Data map[string]interface{} `json:"data"`
    Endpoint string `json:"endpoint"`
    AuthServer string `json:"auth_server"`
    RefreshToken string `json:"refresh_token"`
    Redirect bool `json:"redirect,omitempty"`
    RefreshHeader map[string]interface{} `json:"refresh_header,omitempty"`
    RefreshBody interface{} `json:"refresh_body,omitempty"`
    AccessTokenLocation interface{} `json:"access_token_location,omitempty"`
}

type RefreshResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    APIDomain   string `json:"api_domain"`
    TokenType   string `json:"token_type"`
}

func RefreshAuthToken(r *http.Request, p *InboundPayload) (*http.Response, error) {
    if p.Endpoint == "" || p.AuthServer == "" || p.RefreshToken == "" {
        return nil, errors.New("Missing endpoint, auth_server, or refresh_token")
    }

    var bodyReader io.Reader
    if p.RefreshBody != nil {
        // Marshal RefreshBody to JSON if it contains data
        bodyBytes, err := json.Marshal(p.RefreshBody)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal refresh body: %w", err)
        }
        bodyReader = bytes.NewReader(bodyBytes) // Convert JSON bytes to io.Reader
    }

    // Create the new HTTP request
    refreshReq, err := http.NewRequest(r.Method, p.AuthServer, bodyReader)
    if err != nil {
        return nil, fmt.Errorf("failed to create new request: %w", err)
    }

    // Use RefreshHeader from payload if it exists
    if len(p.RefreshHeader) > 0 {
        for key, value := range p.RefreshHeader {
            // Convert key and value to strings before adding to header
            keyStr := key // it's already a string
            valueStr := fmt.Sprintf("%v", value) // convert value to string using fmt.Sprintf
            refreshReq.Header.Set(keyStr, valueStr)
        }
    } else {
        // Fall back to headers from the original request if RefreshHeader is empty
        for key, values := range r.Header {
            for _, value := range values {
                refreshReq.Header.Add(key, value)
            }
        }
    }


    refreshResp, err := http.DefaultClient.Do(refreshReq)
        if err != nil {
        return refreshResp, err
    }

    return refreshResp, nil
}
