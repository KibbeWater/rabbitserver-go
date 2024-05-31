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
