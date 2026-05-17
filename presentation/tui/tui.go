package tui

import (
	"github.com/kakkky/hakoniwa/usecase"
)

type TUI struct {
	registerResident *usecase.RegisterResident
	sendMessage      *usecase.SendMessageFromBuildingManagerToResident
}

func NewTUI(
	registerResident *usecase.RegisterResident,
	sendMessage *usecase.SendMessageFromBuildingManagerToResident,
) *TUI {
	return &TUI{
		registerResident: registerResident,
		sendMessage:      sendMessage,
	}
}
