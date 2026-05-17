package agents

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type commandSubscriber struct {
	runtime  *Runtime
	cmdInbox domain.AgentCommandInbox
}

func newCommandSubscriber(runtime *Runtime) *commandSubscriber {
	return &commandSubscriber{
		runtime:  runtime,
		cmdInbox: make(domain.AgentCommandInbox, 32),
	}
}

func (ar *commandSubscriber) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case cmd := <-ar.cmdInbox:
			switch cmdV := cmd.(type) {
			case domain.AddResidentAgentCommand:
				ar.runtime.addResidentAgent(&cmdV.Resident)
				// 後続コマンド (e.g. PublishEventCommand) が届くまでに
				// route 登録を完了させるため、ここで reconcile を同期実行する
				ar.runtime.reconcileResidentAgents(ctx)
			case domain.PublishEventCommand:
				ar.runtime.eventBroker.inbox <- cmdV.Event
			}
		}
	}
}
