package testbed

import (
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/quantumUI"
	"google.golang.org/protobuf/encoding/protojson"
)

type Pane struct {
	*pb.QuantumUI
}

func NewPane() *Pane {
	return &Pane{
		&pb.QuantumUI{
			Name:        "First Pane",
			Description: "This is a test pane",
			Title:       "Avatar Pane",
			Components: []*pb.UiComponent{
				{
					Id:          "firstName",
					Label:       "First Name",
					PlaceHolder: "Enter your first name",
					Value:       "",
					Width:       50,
					Type:        pb.UiComponent_INPUT,
				},
				{
					Id:          "lastName",
					Label:       "Last Name",
					PlaceHolder: "Enter your last name",
					Value:       "",
					Width:       50,
					Type:        pb.UiComponent_INPUT,
				},
				{
					Value: "lorem ipsum dolor sit amet consectetur adipiscing elit optio ducimus temporibus ullamco ipsum dolor omnis aut et fugiat fugiat et voluptas et vero est autem minim possimus et libero facere imperdiet nihil consequatur eiusmod nostrud voluptas vero aliquip possimus duis qui dolorum et cupiditate elit minus quas ea voluptas consequatur",
					Width: 100,
					Type:  pb.UiComponent_PLAINTEXT,
				},
			},
		},
	}
}

func (p *Pane) ShowPane() (string, error) {
	jsonBytes, err := protojson.Marshal(p)

	var prettyJSON map[string]interface{}
	err = json.Unmarshal(jsonBytes, &prettyJSON)
	if err != nil {
		log.Fatalf("Failed to unmarshal into map: %v", err)
	}

	prettyJsonBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		log.Fatalf("Failed to indent: %v", err)
	}
	fmt.Println(string(prettyJsonBytes))
	return string(prettyJsonBytes), nil
}
