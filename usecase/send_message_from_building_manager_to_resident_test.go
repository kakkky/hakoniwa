package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/domain/mock"
	"github.com/kakkky/hakoniwa/usecase"
	"go.uber.org/mock/gomock"
)

func TestSendMessageFromBuildingManagerToResident_Exec(t *testing.T) {
	tests := []struct {
		name      string
		to        domain.Resident
		msg       string
		setupMock func(cmder *mock.MockAgentCommander, to domain.Resident, msg string)
	}{
		{
			name: "PublishCommand に PublishEventCommand を渡し成功すれば nil",
			to: domain.Resident{
				ID:   "resident-1",
				Name: "山田",
			},
			msg: "おはようございます",
			setupMock: func(cmder *mock.MockAgentCommander, to domain.Resident, msg string) {
				cmder.EXPECT().
					PublishCommand(gomock.Any(), gomock.Cond(func(cmd domain.AgentCommand) bool {
						pec, ok := cmd.(domain.PublishEventCommand)
						if !ok {
							return false
						}
						me, ok := pec.Event.(domain.MessageEvent)
						if !ok {
							return false
						}
						return me.To().ID == to.ID &&
							me.To().Name == to.Name &&
							me.From().ID == domain.BuildManagerActorID &&
							me.From().Name == domain.BuoldManagerActorName &&
							me.Payload() == msg
					})).
					Return(nil)
			},
		},
		{
			name: "空メッセージでも PublishCommand に到達する",
			to: domain.Resident{
				ID:   "resident-2",
				Name: "花子",
			},
			msg: "",
			setupMock: func(cmder *mock.MockAgentCommander, to domain.Resident, msg string) {
				cmder.EXPECT().
					PublishCommand(gomock.Any(), gomock.Cond(func(cmd domain.AgentCommand) bool {
						pec, ok := cmd.(domain.PublishEventCommand)
						if !ok {
							return false
						}
						me, ok := pec.Event.(domain.MessageEvent)
						return ok && me.To().ID == to.ID && me.Payload() == msg
					})).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cmder := mock.NewMockAgentCommander(ctrl)
			tt.setupMock(cmder, tt.to, tt.msg)

			uc := usecase.NewSendMessageFromBuildingManagerToResident(cmder)
			if err := uc.Exec(context.Background(), tt.to, tt.msg); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSendMessageFromBuildingManagerToResident_Exec_Error(t *testing.T) {
	publishErr := errors.New("publish failed")

	tests := []struct {
		name      string
		setupMock func(cmder *mock.MockAgentCommander)
		wantErr   error
	}{
		{
			name: "PublishCommand がエラーを返すとそのまま返される",
			setupMock: func(cmder *mock.MockAgentCommander) {
				cmder.EXPECT().PublishCommand(gomock.Any(), gomock.Any()).Return(publishErr)
			},
			wantErr: publishErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cmder := mock.NewMockAgentCommander(ctrl)
			tt.setupMock(cmder)

			uc := usecase.NewSendMessageFromBuildingManagerToResident(cmder)
			err := uc.Exec(context.Background(), domain.Resident{ID: "x", Name: "y"}, "msg")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error: got=%v want=%v", err, tt.wantErr)
			}
		})
	}
}
