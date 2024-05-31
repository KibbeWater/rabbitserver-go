package interfaces

type AssistantResponse struct {
	Kernel struct {
		AssistantResponse string `json:"assistantResponse"`
	} `json:"kernel"`
}

type AssistantDeviceResponse struct {
	Kernel struct {
		AssistantResponseDevice struct {
			Text  string `json:"text"`
			Audio string `json:"audio"`
		}
	} `json:"kernel"`
}

type RabbitRegisterResponse struct {
	ActualUserID string `json:"actualUserId"`
	UserID       string `json:"userId"`
	AccountKey   string `json:"accountKey"`
	UserName     string `json:"userName"`
}
