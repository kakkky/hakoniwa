package agents

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kakkky/hakoniwa/domain"
)

const maxMeomryInLLMContext = 50

type residentAgent struct {
	*agentBase
	resident *domain.Resident
}

func newResidentAgent(base *agentBase, resident *domain.Resident) *residentAgent {
	base.id = id(resident.ID)
	base.name = name(resident.Name)
	//
	base.llmPrompt.AddSystemPrompt(``)
	return &residentAgent{
		agentBase: base,
		resident:  resident,
	}
}

func (ra *residentAgent) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-ra.inbox:
			if err := ra.processEvent(ctx, event); err != nil {
				return err
			}
		}
	}
}

func (ra *residentAgent) processEvent(ctx context.Context, event agentEvent) error {
	memories := ra.resident.Memories
	if len(memories) > maxMeomryInLLMContext {
		memories = memories[len(memories)-maxMeomryInLLMContext:]
	}
	var llmContext strings.Builder
	for _, memory := range memories {
		llmContext.WriteString(memory.String() + "\n")
	}

	switch eventV := event.(type) {
	case messageEvent:
		newMemoryContent := fmt.Sprintf(`%sから「%s」とメッセージを受け取った。`, "agent_name", eventV.message)
		newMemory := domain.Memory{
			Content:   newMemoryContent,
			OccuredAt: time.Now(),
		}
		ra.resident.AddMemory(newMemory)
		return ra.processMessageEvent(ctx, eventV, llmContext)
	case opportunityEvent:
		_ = eventV
	}

	return nil
}

func (ra *residentAgent) processMessageEvent(ctx context.Context, msgEvent messageEvent, llmContext strings.Builder) error {
	msg := msgEvent.payroad()

	systemPrompt := ``
	userPromptTemplate := ``

	userPrompt := fmt.Sprintf(userPromptTemplate, msg)

	ra.llmPrompt.AddSystemPrompt(systemPrompt)
	ra.llmPrompt.AddUserPrompt(userPrompt)

	res, err := domain.CallLLM(ctx, ra.llmProvider, ra.llmPrompt, residentLLMResponseSchema, parseResidentLLMResponse)
	if err != nil {
		return err
	}
	ra.handleResidentLLMResponse(res)

	return nil
}

func (ra *residentAgent) handleResidentLLMResponse(res residentLLMResponse) {
	if res.StateUpdate != nil {
		switch {
		case res.StateUpdate.Mood != "":
			ra.resident.UpdateMood(res.StateUpdate.Mood)
		}
	}
	action := res.Action
	switch action {
	case ResidentActionMessage:
		event := &messageEvent{
			eventBase: eventBase{
				toID:   id(res.To),
				fromID: id(ra.id),
			},
			message: res.Payload,
		}
		ra.sendEvent(event)
	case ResidentActionOpportunity:
		event := &opportunityEvent{
			eventBase: eventBase{
				toID:   id(res.To),
				fromID: id(ra.id),
			},
			opportunity: res.Payload,
		}
		ra.sendEvent(event)
	case ResidentActionStay:
		return
	}
}

//go:embed schema/resident_agent_llm_response_schema.json
var residentLLMResponseSchema string

type ResidentActionKind string

const (
	ResidentActionMessage    ResidentActionKind = "message"
	ResidentActionOpportunity ResidentActionKind = "opportunity"
	ResidentActionStay       ResidentActionKind = "stay"
)

type residentLLMResponse struct {
	Action      ResidentActionKind   `json:"action"`
	To          domain.ResidentID    `json:"to,omitempty"`
	Payload     string               `json:"payload"`
	StateUpdate *residentStateUpdate `json:"state_update,omitempty"`
}

type residentStateUpdate struct {
	Mood domain.Mood `json:"mood,omitempty"`
}

func parseResidentLLMResponse(raw domain.LLMResponse) (residentLLMResponse, error) {
	var res residentLLMResponse
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return residentLLMResponse{}, err
	}
	return res, nil
}
