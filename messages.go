package thegotribe

type Request struct {
	Category string      `json:"category"`
	Request  string      `json:"request"`
	Values   interface{} `json:"values"`
}

type Response struct {
	Request
	StatusCode int `json:"statuscode"`
}
