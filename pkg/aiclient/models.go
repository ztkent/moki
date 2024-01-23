package aiclient

type AnyscaleModel string
type OpenAIModel string

const (
	Mistral7BInstruct   AnyscaleModel = "mistralai/Mistral-7B-Instruct-v0.1"
	Llama27bChat        AnyscaleModel = "meta-llama/Llama-2-7b-chat-hf"
	Llama213bChat       AnyscaleModel = "meta-llama/Llama-2-13b-chat-hf"
	Mixtral8x7BInstruct AnyscaleModel = "mistralai/Mixtral-8x7B-Instruct-v0.1"
	GPT35Turbo          OpenAIModel   = "gpt-3.5-turbo"
)

func (a AnyscaleModel) String() string {
	return string(a)
}

func (o OpenAIModel) String() string {
	return string(o)
}
