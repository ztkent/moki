package aiclient

type AnyscaleModel string
type OpenAIModel string

const (
	Mistral7BInstruct   AnyscaleModel = "mistralai/Mistral-7B-Instruct-v0.1"
	Llama27bChat        AnyscaleModel = "meta-llama/Llama-2-7b-chat-hf"
	Llama213bChat       AnyscaleModel = "meta-llama/Llama-2-13b-chat-hf"
	Llama270bChat       AnyscaleModel = "meta-llama/Llama-2-70b-chat-hf"
	Mixtral8x7BInstruct AnyscaleModel = "mistralai/Mixtral-8x7B-Instruct-v0.1"
	CodeLlama34b        AnyscaleModel = "codellama/CodeLlama-34b-Instruct-hf"
	CodeLlama70b        AnyscaleModel = "codellama/CodeLlama-70b-Instruct-hf"
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
	case Llama270bChat.String(), "l70b":
		return Llama270bChat, true
	case CodeLlama34b.String(), "cl34b":
		return CodeLlama34b, true
	case CodeLlama70b.String(), "cl70b":
		return CodeLlama70b, true
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
