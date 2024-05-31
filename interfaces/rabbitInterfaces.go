package interfaces

type AssistantResponse struct {
	Kernel struct {
		AssistantResponse string `json:"assistantResponse"`
	} `json:"kernel"`
}
