package agents

import "github.com/kakkky/hakoniwa/domain"

type commandSubscriber struct {
	runtime  *Runtime
	cmdInbox cmdInbox
}

type cmdInbox chan domain.AgentCommand

func newCommandSubscriber() *commandSubscriber {
	return &commandSubscriber{}
}

func (ar *commandSubscriber) run() error {
	for {
		select {
		case cmd := <-ar.cmdInbox:
			switch cmdV := cmd.(type) {
			case domain.AddResidentAgentCommand:
				ar.runtime.addResidentAgent(&cmdV.Resident)
			}
		}
	}
}
