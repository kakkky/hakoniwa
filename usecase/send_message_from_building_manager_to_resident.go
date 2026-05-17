package usecase

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type SendMessageFromBuildingManagerToResident struct {
	agentCommander domain.AgentCommander
}

func NewSendMessageFromBuildingManagerToResident(agentCommander domain.AgentCommander) *SendMessageFromBuildingManagerToResident {
	return &SendMessageFromBuildingManagerToResident{
		agentCommander: agentCommander,
	}
}

func (sm *SendMessageFromBuildingManagerToResident) Exec(ctx context.Context, to domain.Resident, msg string) error {
	cmd := domain.PublishEventCommand{
		Event: domain.MessageEvent{
			EventBase: domain.EventBase{
				EventTo: domain.EventTo{
					ID:   to.ID,
					Name: to.Name,
				},
				EventFrom: domain.EventFrom{
					ID:   domain.BuildManagerActorID,
					Name: domain.BuoldManagerActorName,
				},
			},
			Message: msg,
		},
	}
	if err := sm.agentCommander.PublishCommand(ctx, cmd); err != nil {
		return err
	}
	return nil
}
