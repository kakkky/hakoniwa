package presentation

import (
	"github.com/kakkky/hakoniwa/usecase"
)

type TextUserInterface struct {
	registerResident *usecase.RegisterResident
	sendMessage      *usecase.SendMessageFromBuildingManagerToResident
}

func NewTextUserInterface(
	registerResident *usecase.RegisterResident,
	sendMessage *usecase.SendMessageFromBuildingManagerToResident,
) *TextUserInterface {
	return &TextUserInterface{
		registerResident: registerResident,
		sendMessage:      sendMessage,
	}
}
