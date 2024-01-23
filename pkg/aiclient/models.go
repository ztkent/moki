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

func IsAnyscaleModel(name string) (AnyscaleModel, bool) {
	switch name {
	case Mistral7BInstruct.String(), "m7b":
		return Mistral7BInstruct, true
	case Llama27bChat.String(), "l7b":
		return Llama27bChat, true
	case Llama213bChat.String(), "l13b":
		return Llama213bChat, true
	case Mixtral8x7BInstruct.String(), "m8x7b":
		return Mixtral8x7BInstruct, true
	default:
		return "", false
	}
}

func IsOpenAIModel(name string) (OpenAIModel, bool) {
	switch name {
	case GPT35Turbo.String(), "turbo":
		return GPT35Turbo, true
	default:
		return "", false
	}
}
