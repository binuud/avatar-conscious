package llm_test

import (
	"os"
	"strings"
	"testing"

	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	prompt       = "What is the value of 56 + 100"
	systemPrompt = "You are an friendly assistant, answer in brief. If you do not know the answer, do not try to guess."
)

func init() {
	bootstrap.Bootstrap()
	preCleanUp()
}

func preCleanUp() {
	audioFile := utils.GetGeneratedTTSFolder() + "test--tts-llm.mp3"
	err := os.Remove(audioFile)
	if err != nil {
		logrus.Errorf("Error deleting file:", err)
		return
	}
}

func Test_Get_Endpoint(t *testing.T) {

	got := llm.GetLLMEndpoint()
	if got == "" {
		t.Fatal("error llm endpoint not configured")
	}
	t.Logf("got llm endpoint %s", got)

}
func Test_Get_Models(t *testing.T) {

	got, err := llm.MakeLLMGetModels()
	if err != nil {
		t.Errorf("error Test_Get_Models: %v", err)
		return
	}
	if len(got) == 0 {
		t.Error("error Test_Get_Models: no models received")
		return
	}
	t.Logf("got %d models from llm", len(got))
}

func Test_LLM_Single_Request(t *testing.T) {

	model := utils.GetEnv("LLM_MODEL", "lmstudio-community/Qwen3.5-9B-MLX-4bit")
	got, err := llm.MakeLLMSingleRequest(prompt, systemPrompt, model, 0.8)
	if err != nil {
		t.Errorf("error Test_LLM_Chat: %v", err)
		return
	}

	if strings.Contains(got, "156") {
		t.Log("got expected value")
	}

}

func Test_LLM_Chat_Request(t *testing.T) {

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(prompt),
	}

	model := utils.GetEnv("LLM_MODEL", "mlx-community/Qwen3-TTS-12Hz-0.6B-Base-8bit")
	completion, err := llm.MakeLLMChatRequest(messages, nil, model, 0.8)
	if err != nil {
		t.Errorf("error Test_LLM_Chat_Request: %v", err)
		return
	}

	var got = completion.Choices[0].Message.Content
	if strings.Contains(got, "156") {
		t.Log("got expected value Test_LLM_Chat_Request")
	}

}

func Test_LLM_Audio_Request(t *testing.T) {

	tts_model := "mlx-community/Qwen3-TTS-12Hz-0.6B-Base-8bit"

	audioFile := utils.GetGeneratedTTSFolder() + "test--tts-llm.mp3"
	err := llm.MakeTTSRequest("Hi how are you", tts_model, audioFile)
	if err != nil {
		t.Errorf("error Test_LLM_Audio_Request: %v", err)
		return
	}
	if utils.FileExists(audioFile) {
		t.Logf("got expected audio: Test_LLM_Audio_Request: in file %s", audioFile)
	} else {
		t.Errorf("error log file not found %s", audioFile)
	}

}
