package common

import "encoding/json"

type DNAMonitorRequest struct {
	Qid    string                 `json:"qid"`
	Method string                 `json:"method"`
	Params map[string]interface{} `josn:"params"`
}

func NewDNAMonitorRequest(data []byte) (*DNAMonitorRequest, error) {
	req := &DNAMonitorRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

type DNAMonitorResponse struct {
	Qid         string      `json:"qid"`
	Method      string      `json:"method"`
	ErrorCode   int         `json:"errorcode"`
	ErrorString string      `json:"errorstring"`
	Result      interface{} `json:"result"`
}

func NewDNAMonitorResponse(qid, method string, result interface{}, errorCode int) *DNAMonitorResponse {
	return &DNAMonitorResponse{
		Qid:       qid,
		Method:    method,
		Result:    result,
		ErrorCode: errorCode,
	}
}

func (this *DNAMonitorResponse) Marshal() ([]byte, error) {
	this.ErrorString = GetErrDesc(this.ErrorCode)
	return json.Marshal(this)
}
