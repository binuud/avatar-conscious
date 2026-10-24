package main

import (
	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarServer"
	"github.com/sirupsen/logrus"
)

func main() {
	bootstrap.Bootstrap()
	logrus.Infof("Test_Avatar election-officer")
	a := avatarServer.NewAvatarFromStore("election-officer")
	a.Build()
	a.Save()
	a.Run()
}
