package interfaces

type APIRequest struct {
	Type string `json:"type"`
}

type LogonRequest struct {
	Type string `json:"type"`
	Data struct {
		IMEI       string `json:"imei"`
		AccountKey string `json:"accountKey"`
	} `json:"data"`
}

type MessageRequest struct {
	Type    string `json:"type"`
	Message string `json:"data"`
}

type PTTRequest struct {
	Type string `json:"type"`
	Data struct {
		Image  string `json:"image"`
		Active bool   `json:"active"`
	} `json:"data"`
}

type AudioRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LogonResponse struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type MessageResponse struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type AudioMessageResponse struct {
	Type string `json:"type"`
	Data struct {
		Text  string `json:"text"`
		Audio string `json:"audio"`
	} `json:"data"`
}