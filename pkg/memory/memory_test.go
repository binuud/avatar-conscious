package memory

import (
	"testing"

	"github.com/openai/openai-go/v3"
)

func Test_Memory(t *testing.T) {
	m := NewMemory()
	err := m.Save()
	if err != nil {
		t.Errorf("client: cannot save memory as binary file: %v \n", err)
		return
	}
	t.Logf("saved memory to binary file")
	// do not save tools because of issue with saving hash map
	newMemory := &Memory{Name: "test-memory"}
	err = newMemory.Load()
	if err != nil {
		t.Errorf("client: cannot load memory from binary file: %v \n", err)
		return
	}
	t.Logf("loaded memory from binary file")
	newMemory.AppendMemory(openai.UserMessage("This is prompt one"))

	m = &Memory{Name: "test-memory"}
	err = m.Load()
	if err != nil {
		t.Errorf("client: cannot load memory from binary file: %v \n", err)
		return
	}
	m.DumpMemory()
	lastMessage := m.Messages[len(m.Messages)-1]
	if lastMessage.OfUser != nil {
		if lastMessage.OfUser.Content.OfString.Value == "This is prompt one" {
			t.Log("got appended message from saved memory")
		} else {
			t.Errorf("cannot save memory to binary file after append: %v \n", err)
		}
	} else {
		t.Errorf("cannot save memory to binary file after append: %v \n", err)
	}

}
