package memory_test

import (
	"testing"

	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

func init() {
	bootstrap.Bootstrap()
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_Memory(t *testing.T) {
	// Step 1: save memory
	t.Run("Step=Save_Memory", func(t *testing.T) {
		_save_memory(t)
	})

	// Step 2: Depends on Step 1
	t.Run("Step=Load_Memory", func(t *testing.T) {
		_load_memory(t)
	})

}

func _save_memory(t *testing.T) {

	memoryName := "test-memory"

	m := memory.NewMemory(memoryName)
	m.AppendMessage(openai.UserMessage("This is prompt one"))
	err := m.Save()
	if err != nil {
		t.Errorf("client: cannot save memory as binary file: %v \n", err)
		return
	}
	m.DumpMemory(nil)
	t.Logf("able to save memory to binary file")
}

func _load_memory(t *testing.T) {

	memoryName := "test-memory"
	m, err := memory.LoadMemory(memoryName)
	if err != nil {
		// memory file not found
		t.Errorf("error cannot load memory")
		return
	} else {
		t.Log("able to load file")
	}

	if m.ToolRegistry != nil && len(m.ToolRegistry.Tools) > 0 {
		t.Logf("Tools have been saved %v", m.ToolRegistry.Tools)
	} else {
		t.Errorf("error cannot save tools")
	}

	lastMessage := m.Messages[len(m.Messages)-1]
	if lastMessage.OfUser != nil {
		if lastMessage.OfUser.Content.OfString.Value == "This is prompt one" {
			t.Log("got appended message from saved memory")
		} else {
			t.Errorf("cannot load memory from binary file, after save: %v \n", err)
		}
	} else {
		t.Errorf("cannot save messages: %v \n", err)
	}

	m.AppendMessage(openai.UserMessage("This is prompt two"))

	m.DumpMemory(nil)
	lastMessage = m.Messages[len(m.Messages)-1]
	if lastMessage.OfUser != nil {
		if lastMessage.OfUser.Content.OfString.Value == "This is prompt two" {
			t.Log("got appended message from saved memory")
		} else {
			t.Errorf("cannot save memory to binary file after append: %v \n", err)
		}
	} else {
		t.Errorf("cannot save memory to binary file after append: %v \n", err)
	}

}
