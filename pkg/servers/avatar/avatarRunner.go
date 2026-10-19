package avatar

import (
	pb "github.com/binuud/conscious/gen/go/v1/avatar"
	"github.com/sirupsen/logrus"
)

type Avatar struct {
	*pb.Avatar
}

func NewAvatar() *Avatar {
	avatar := &pb.Avatar{
		Uuid:        "2b957a49-d9cf-449e-b435-3c1097ca318a",
		Name:        "Sample Election Commision SIR",
		Description: "Demo to showcase avatar agent, that asks users for election ID, and compares same with election commission DB and voter list.",
		Loops: []*pb.AvatarLoop{
			{
				Name:        "Check Voter ID",
				Description: "Check if the voter id is valid against mock server",
				RunActions: []*pb.AvatarAction{
					{
						Name:        "Get Voter ID Details",
						Description: "Ask user to show their voter id on web camera, and get the voter ID, Name details from the webcamera",
						Disabled:    false,
						Type:        pb.AvatarActionType_AiImageDescribe,
						Prompt:      "Look at the image, and get the details of Name, Voter ID. Keep your response short, and report the identifiers from the image as a json object.",
					},
					{
						Name:        "Confirm Valid Voter ID",
						Description: "Confirm that the Voter ID is valid, by checking with the election commission server",
						Disabled:    false,
						Type:        pb.AvatarActionType_WebApi,
						Prompt:      "Look at the image, and get the details of Name, Voter ID. Keep your response short, and report the identifiers from the image as a json object.",
					},
				},
			},
			{
				Name:        "Check Aadhar ID",
				Description: "Check if the aadhar id is valid against mock server",
			},
		},
	}
	// logrus.Infof("Avatar %v \n", avatar)
	return &Avatar{avatar}
}

func (a *Avatar) Run() {
	logrus.Infof("Running Avatar %s/%s \n", a.Uuid, a.Name)
	logrus.Infof("Description %s \n", a.Description)
	for _, item := range a.Loops {
		logrus.Infof("Running Loop %s \n", item.Name)
		logrus.Infof("Description %s \n", item.Description)
	}
}
