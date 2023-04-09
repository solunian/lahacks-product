package main

import (
	"LAH-7/prompter/cohere"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

const AUTH_STRING = "LAH-SPRING-2023-9b2a97fa6ef69419"

const (
	RESPONSE_TYPE_NEW_TEXT = 0
	RESPONSE_TYPE_FINAL    = 1
)

type ResponseNewTextMessage struct {
	Type    int    `json:"t"`
	NewText string `json:"n"`
}

type Link struct {
	ID        string  `json:"id"`
	URL       string  `json:"url"`
	Name      string  `json:"name"`
	Body      string  `json:"body"`
	Certainty float32 `json:"certainty"`
}

type ResponseFinalMessage struct {
	Type      int     `json:"t"`
	FinalText string  `json:"final_text"`
	Links     []*Link `json:"links"`
}

type ConversationRequest struct {
	Question string `json:"question"`
	History  []*QASet
}

func stupidAuthMiddleware(ctx iris.Context) {
	if ctx.GetHeader("Authorization") != AUTH_STRING {
		ctx.StatusCode(401)
		return
	}
	ctx.Next()
}

type GraphQLSearchResponse struct {
	MedObject []struct {
		Additional struct {
			Certainty float32 `mapstructure:"certainty"`
			Distance  float32 `mapstructure:"distance"`
		} `mapstructure:"_additional"`
		Description string `mapstructure:"description"`
		DocID       string `mapstructure:"doc_id"`
		Name        string `mapstructure:"name"`
		URL         string `mapstructure:"url"`
	} `mapstructure:"MedObject"`
}

func GetLinks(cohereToken string, assistantResponse string, c *weaviate.Client) ([]*Link, error) {
	embeddings, err := cohere.GetEmbeddings(cohereToken, []string{assistantResponse})
	if err != nil {
		return nil, err
	}

	className := "MedObject"

	name := graphql.Field{
		Name: "name",
	}
	id := graphql.Field{
		Name: "doc_id",
	}
	url := graphql.Field{
		Name: "url",
	}
	description := graphql.Field{
		Name: "description",
	}
	_additional := graphql.Field{
		Name: "_additional", Fields: []graphql.Field{
			{Name: "certainty"}, // only supported if distance==cosine
			{Name: "distance"},  // always supported
		},
	}

	nearVector := c.GraphQL().NearVectorArgBuilder().
		WithVector(embeddings[0]).
		WithCertainty(0.75)

	ctx := context.Background()
	result, err := c.GraphQL().Get().
		WithClassName(className).
		WithFields(name, id, url, description, _additional).
		WithNearVector(nearVector).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	// r := result.Get
	var d GraphQLSearchResponse
	err = mapstructure.Decode(result.Data["Get"], &d)
	if err != nil {
		return nil, err
	}

	var links []*Link

	for _, v := range d.MedObject {
		links = append(links, &Link{
			ID:        v.DocID,
			URL:       v.URL,
			Name:      v.Name,
			Body:      v.Description,
			Certainty: v.Additional.Certainty,
		})
	}
	return links, nil
}

func main() {
	cohereToken := os.Getenv("COHERE_API_KEY")

	config := weaviate.Config{
		Scheme: "http",
		Host:   "localhost:8080",
	}
	weaveClient := weaviate.New(config)

	app := iris.New()

	app.Options("/conversation", func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.StatusCode(201)
		ctx.Done()
	})

	app.Use(stupidAuthMiddleware)

	// Fixed: updated JSON endpoint to return a JSON array instead of separate JSON objects
	app.Post("/conversation", func(ctx iris.Context) {
		var req ConversationRequest

		ctx.ReadJSON(&req)

		ctx.Header("Content-Type", "text/event-stream; charset=utf-8")
		ctx.Header("Transfer-Encoding", "chunked")

		prompt, err := GeneratePrompt(&PromptTemplateInput{
			History:    req.History,
			HumanInput: req.Question,
		})
		if err != nil {
			ctx.StopWithError(500, err)
		}
		log.Print(prompt)

		input := cohere.GenerateInput{
			Model:             "31989c7e-2182-4e52-b4c8-2c3d04d54626-ft",
			K:                 15,
			Prompt:            prompt,
			StopSequences:     []string{"\n"},
			ReturnLikelihoods: "NONE",
			MaxTokens:         250,
			Temperature:       0.3,
			Stream:            true,
		}

		results := make(chan cohere.Result)

		finalText := ""

		go func() {
			cohere.Generate(cohereToken, &input, results)
			if err != nil {
				log.Fatal(err)
			}
		}()

		notifyClose := ctx.Request().Context().Done()
		for word := range results {
			select {
			case <-notifyClose:
				ctx.Application().Logger().Infof("Connection closed, loop end.")
				return
			default:
				finalText += word.Text
				wordMessage, err := json.Marshal(&ResponseNewTextMessage{
					Type:    RESPONSE_TYPE_NEW_TEXT,
					NewText: word.Text,
				})
				if err != nil {
					ctx.StopWithError(500, err)
				}
				ctx.Write(wordMessage)
				ctx.WriteString("\n")

				ctx.ResponseWriter().Flush()
			}
		}

		links, err := GetLinks(cohereToken, finalText, weaveClient)

		finalMessage, err := json.Marshal(&ResponseFinalMessage{
			Type:      RESPONSE_TYPE_FINAL,
			FinalText: strings.TrimSpace(finalText),
			Links:     links,
		})
		if err != nil {
			ctx.StopWithError(500, err)
		}
		ctx.Write(finalMessage)
		ctx.WriteString("\n")

		ctx.Application().Logger().Infof("Loop end.")
		ctx.ResponseWriter().Flush()
	})

	app.Listen(":3000")
}
