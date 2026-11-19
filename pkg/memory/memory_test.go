package memory_test

import (
	"testing"

	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/sirupsen/logrus"
)

var (
	memoryName = "test-memory"
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
	t.Run("Step=Clear_Memory", func(t *testing.T) {
		_clear_memory(t)
	})
	t.Run("Step=Delete_Memory", func(t *testing.T) {
		_delete_memory(t)
	})

}

func _save_memory(t *testing.T) {

	m := memory.NewEmtpyMemory()
	m.Meta.Name = memoryName
	m.AppendUserMessage("This is prompt one")
	m.AppendAssistantMessage("Ok prompt one")
	m.DumpTools()
	err := m.Save() // this is reduntant, but to check if save returns error
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

	// we are using global tools
	// to prevent tool loading at every http request
	// if m.ToolRegistry != nil && len(m.ToolRegistry.Tools) > 0 {
	// 	t.Logf("Tools have been saved %v", m.ToolRegistry.Tools)
	// } else {
	// 	t.Errorf("error cannot save tools")
	// }

	lastMessage := m.Messages[len(m.Messages)-1]
	if lastMessage.OfAssistant != nil {
		if lastMessage.OfAssistant.Content.OfString.Value == "Ok prompt one" {
			t.Log("got appended message from saved memory")
		} else {
			t.Errorf("cannot load memory from binary file, after save: ")
		}
	} else {
		t.Errorf("cannot save messages")
	}
	m.AppendUserMessage("This is prompt two")
	m.AppendAssistantMessage("Ok prompt two")
	m.DumpMemory(nil)
	lastMessage = m.Messages[len(m.Messages)-2]
	if lastMessage.OfUser != nil {
		if lastMessage.OfUser.Content.OfString.Value == "This is prompt two" {
			t.Log("got appended message from saved memory")
		} else {
			t.Errorf("cannot save memory to binary file after append:")
		}
	} else {
		t.Errorf("cannot save memory to binary file after append:")
	}

	lastContext := m.Contexts[len(m.Contexts)-1]
	if lastContext.Out.Response == "Ok prompt two" {
		t.Log("got appended context message from saved memory")
	} else {
		t.Errorf("cannot save context memory to binary file after append:")
	}
}

func _clear_memory(t *testing.T) {

	memoryName := "test-memory"
	m := memory.NewMemory(memoryName)
	m.DumpMemory(nil)
	lastMessage := m.Messages[len(m.Messages)-1]
	if lastMessage.OfAssistant != nil {
		if lastMessage.OfAssistant.Content.OfString.Value == "Ok prompt two" {
			t.Log("got appended message from saved memory")
		} else {
			t.Errorf("cannot load memory from binary file")
		}
	} else {
		t.Errorf("cannot save messages")
	}

	m.ClearMemory()
	m.DumpMemory(nil)
	if len(m.Messages) > 1 { // first memory is always system prompt
		t.Errorf("cannot clear memory")
	}
	if memoryName != m.Meta.Name {
		t.Errorf("error when clearing memory meta got cleared")
	}

	if len(m.Contexts) > 0 {
		t.Errorf("cannot clear context memory")
	}

}

func _delete_memory(t *testing.T) {

	err := memory.DeleteMemory(memoryName)
	if err != nil {
		// memory file not found
		t.Errorf("error cannot delete memory")
		return
	} else {
		t.Log("able to delete file")
	}

	_, err = memory.LoadMemory(memoryName)
	if err == nil {
		// memory file not found
		t.Errorf("error able to load deleted memory")
		return
	} else {
		t.Log("not able to load deleted file")
	}
}
