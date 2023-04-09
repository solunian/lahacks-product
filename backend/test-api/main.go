package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
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
	ID   string `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Body string `json:"body"`
}

type ResponseFinalMessage struct {
	Type      int     `json:"t"`
	FinalText string  `json:"final_text"`
	Links     []*Link `json:"links"`
}

type ConversationRequest struct {
	Question string `json:"question"`
	History  []struct {
		Assistant string `json:"assistant"`
		Human     string `json:"human"`
	}
}

func stupidAuthMiddleware(ctx iris.Context) {
	if ctx.GetHeader("Authorization") != AUTH_STRING {
		ctx.StatusCode(401)
		return
	}
	ctx.Next()
}

func main() {
	app := iris.New()

	type messageNumber struct {
		Number int `json:"number"`
	}

	app.Get("/json", func(ctx iris.Context) {
		ctx.Header("Transfer-Encoding", "chunked")
		i := 0
		ints := []int{1, 2, 3, 5, 7, 9, 11, 13, 15, 17, 23, 29}
		// Send the response in chunks and wait for half a second between each chunk,
		// until connection close.
		notifyClose := ctx.Request().Context().Done()
		for {
			select {
			case <-notifyClose:
				// err := ctx.Request().Context().Err()
				ctx.Application().Logger().Infof("Connection closed, loop end.")
				return
			default:
				ctx.JSON(messageNumber{Number: ints[i]})
				ctx.WriteString("\n")
				time.Sleep(500 * time.Millisecond) // simulate delay.
				if i == len(ints)-1 {
					ctx.Application().Logger().Infof("Loop end.")
					return
				}
				i++
				ctx.ResponseWriter().Flush()
			}
		}
	})

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

		text := req.Question
		words := strings.Split(text, " ")

		ctx.Header("Content-Type", "text/event-stream; charset=utf-8")
		ctx.Header("Transfer-Encoding", "chunked")
		ctx.Header("Access-Control-Allow-Origin", "*")

		notifyClose := ctx.Request().Context().Done()
		for i, word := range words {
			select {
			case <-notifyClose:
				ctx.Application().Logger().Infof("Connection closed, loop end.")
				return
			default:
				if i == len(words)-1 {
					finalMessage, err := json.Marshal(&ResponseFinalMessage{
						Type:      RESPONSE_TYPE_FINAL,
						FinalText: text,
						Links: []*Link{
							{
								ID:   "0",
								URL:  "https://google.com",
								Name: "Google",
								Body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Laoreet suspendisse interdum consectetur libero id faucibus nisl tincidunt eget. Volutpat odio facilisis mauris sit.",
							},
						},
					})
					if err != nil {
						ctx.StopWithError(500, err)
					}
					ctx.Write(finalMessage)
					ctx.WriteString("\n")

					ctx.Application().Logger().Infof("Loop end.")
					ctx.ResponseWriter().Flush()
					return
				}

				wordMessage, err := json.Marshal(&ResponseNewTextMessage{
					Type:    RESPONSE_TYPE_NEW_TEXT,
					NewText: " " + word,
				})
				if err != nil {
					ctx.StopWithError(500, err)
				}
				ctx.Write(wordMessage)
				ctx.WriteString("\n")

				time.Sleep(250 * time.Millisecond)

				ctx.ResponseWriter().Flush()
			}
		}
	})

	app.Listen(":3000")
}
