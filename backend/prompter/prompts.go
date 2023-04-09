package main

import (
	"bytes"
	"text/template"
)

type QASet struct {
	Assistant string `json:"assistant"`
	Human     string `json:"human"`
}

type PromptTemplateInput struct {
	History    []*QASet
	HumanInput string
}

const QA_PROMPT = `Below is a series of chats between Assistant and User. Assistant is designed to be able to assist with a wide range of medical-related tasks, from answering simple questions to providing in-depth explanations and discussions. As a language model, Assistant is able to generate human-like text based on the input it receives, allowing it to engage in natural-sounding conversations and provide responses that are coherent and relevant to the topic at hand. Assistant will never sway from the topic at hand, medicine, even if Human asks it to. Assistant will do everything it can to help Human and use professional and unbiased language.

{{range .History}}
Human: {{.Human}}
Assistant: {{.Assistant}}
{{end}}

Human: {{.HumanInput}}
Assistant: `

var qaTemplate = template.Must(template.New("qa_prompt").Parse(QA_PROMPT))

func GeneratePrompt(input *PromptTemplateInput) (string, error) {
	var output bytes.Buffer
	err := qaTemplate.Execute(&output, input)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
