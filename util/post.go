package util

import  (
  "net/http"
  "encoding/json"
  "bytes"
  "fmt"
)

func Post(req *http.Request, resp *http.Response, p *InboundPayload) (*http.Response, error) {
  var refreshResp RefreshResponse
  err := json.NewDecoder(resp.Body).Decode(&refreshResp)
  if err != nil {
    return nil, fmt.Errorf("Error marshalling JSON: %v", err)
    }
  reqBody, err := json.Marshal(p.Data)
  if err != nil {
    return nil, fmt.Errorf("Error marshalling JSON: %v", err)
    }

  postReq, err := http.NewRequest(req.Method, p.Endpoint, bytes.NewReader(reqBody))
  if err != nil {
    return nil, fmt.Errorf("Error posting lead: %v", err)
  }

	// Copy all headers from r to req
	for key, values := range req.Header {
		for _, value := range values {
			postReq.Header.Add(key, value)
		}
	}
  postResp, err := http.DefaultClient.Do(postReq)
  if err != nil {
    return postResp, err
  }

  return postResp, nil
}



