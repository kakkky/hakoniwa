package usecase

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type SendMessageFromBuildingManagerToResident struct {
	agentCommandPublisher domain.AgentCommandPublisher
}

func NewSendMessageFromBuildingManagerToResident(agentCommandPublisher domain.AgentCommandPublisher) *SendMessageFromBuildingManagerToResident {
	return &SendMessageFromBuildingManagerToResident{
		agentCommandPublisher: agentCommandPublisher,
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
	if err := sm.agentCommandPublisher.PublishCommand(ctx, cmd); err != nil {
		return err
	}
	return nil
}
