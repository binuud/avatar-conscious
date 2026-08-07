package consciousServer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	aiModuleProto "github.com/dronasys-com/brain-nodes/gen/go/v1/aiModule"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
}

func MakeLLMChatRequest(prompt string, messages []openai.ChatCompletionMessageParamUnion, model string, temperature float64) (*openai.ChatCompletion, error) {

	llmURL := os.Getenv("LLM_ENDPOINT") + "/engines/llama.cpp/v1/"

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	ctx := context.Background()

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Opt(temperature),
	}

	completion, err := client.Chat.Completions.New(ctx, param)

	if err != nil {
		return nil, err
	}

	return completion, nil

}

func MakeLLMSingleRequest(prompt string, systemPrompt string, model string, temperature float64) (string, error) {

	llmURL := os.Getenv("LLM_ENDPOINT") + "/engines/llama.cpp/v1/"

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	ctx := context.Background()

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(prompt),
	}

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Opt(temperature),
	}

	completion, err := client.Chat.Completions.New(ctx, param)

	if err != nil {
		return "", err
	}

	var response = completion.Choices[0].Message.Content

	fmt.Println(response)
	return response, nil

}

func MQTTPubTts(text string) error {

	topic := "aiModule/tts"
	data := aiModuleProto.GenerateRequest{
		Name:        "tts",
		Description: "Tts Publish",
		Prompt:      text,
	}

	// Now you can use `data` as your typed object
	// logrus.Infof("Parsed data: %+v\n", data.Uuid)

	// Marshal converts the struct into a byte slice
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatal(err)
	}

	// Convert byte slice directly to a string
	jsonString := string(jsonData)
	logrus.Infof("Out string %s", jsonString)

	token := mqttClient.Publish(topic, 1, false, jsonData)
	token.Wait() // Block until the broker acknowledges receipt

	if token.Error() != nil {
		log.Printf("Failed to publish message: %v", token.Error())
	} else {
		fmt.Printf("Successfully published payload: %s\n", jsonData)
	}
	return token.Error()
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("\n Failed to upgrade websocket: %v", err)
		return
	}
	defer ws.Close()

	// logrus.Infof("Handling request %v", r)

	// Subscribe to a topic
	mqttClient.Subscribe("aiModule/response", 0, func(client mqtt.Client, msg mqtt.Message) {
		// Forward MQTT messages to WebSocket client
		logrus.Printf("\n Received aiModule/response on mqtt, send to websocket: %s", string(msg.Payload()))
		err := ws.WriteMessage(websocket.TextMessage, msg.Payload())
		if err != nil {
			logrus.Errorf("Write to websocket error: %v", err)
		}
	})

	for {
		// Read message from WebSocket client
		_, msg, err := ws.ReadMessage()
		if err != nil {
			logrus.Errorf("Read from websocket error: %v", err)
			break
		}

		logrus.Infof("\nReceived from WS: %s\n", string(msg))

		// var data aiModuleProto.GenerateRequest
		// err = protojson.Unmarshal(msg, &data)
		// // err = json.Unmarshal(msg, &data)
		// if err != nil {
		// 	logrus.Errorf("Error unmarshaling JSON: %v", err)
		// 	continue
		// }

		// // Now you can use `data` as your typed object
		// // logrus.Infof("Parsed data: %+v\n", data.Uuid)
		// var token mqtt.Token
		// switch data.Name {
		// case "Tts":
		// 	token = mqttClient.Publish("aiModule/tts", 0, false, msg)
		// case "Image":
		// 	token = mqttClient.Publish("aiModule/image", 0, false, msg)
		// case "Vision":
		// 	token = mqttClient.Publish("aiModule/vision", 0, false, msg)
		// }

		// token.Wait()
	}
}
