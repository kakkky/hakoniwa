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

const maxMemoryInLLMContext = 50

type residentAgent struct {
	agentBase
	resident *domain.Resident
}

func newResidentAgent(base agentBase, resident *domain.Resident) *residentAgent {
	systemPromptTemplate := `レスポンスは以下のJSONスキーマの文字列で行うようにしてください。：%s`
	systemPrompt := fmt.Sprintf(systemPromptTemplate, residentLLMResponseSchema)
	base.llmPrompt.AddSystemPrompt(systemPrompt)
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

func (ra *residentAgent) processEvent(ctx context.Context, event domain.Event) error {
	memories := ra.resident.Memories
	if len(memories) > maxMemoryInLLMContext {
		memories = memories[len(memories)-maxMemoryInLLMContext:]
	}
	var llmContext strings.Builder
	for _, memory := range memories {
		llmContext.WriteString(memory.String() + "\n")
	}

	switch eventV := event.(type) {
	case domain.MessageEvent:
		newMemoryContent := fmt.Sprintf(`%sから「%s」とメッセージを受け取った。`, eventV.EventFrom.Name, eventV.Message)
		newMemory := domain.Memory{
			Content:   newMemoryContent,
			OccuredAt: time.Now(),
		}
		ra.resident.AddMemory(newMemory)
		return ra.processMessageEvent(ctx, eventV, llmContext)
	case domain.OpportunityEvent:
		_ = eventV
	}

	return nil
}

func (ra *residentAgent) processMessageEvent(ctx context.Context, msgEvent domain.MessageEvent, llmContext strings.Builder) error {
	msg := msgEvent.Payload()

	systemPromptTemplate := ``
	systemPrompt := fmt.Sprintf(systemPromptTemplate, nil)
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
		ra.sendEvent(domain.MessageEvent{
			EventBase: domain.EventBase{
				EventTo: domain.EventTo{ID: res.To},
				EventFrom: domain.EventFrom{
					ID:   domain.ActorID(ra.resident.ID),
					Name: domain.ActorName(ra.resident.Name),
				},
			},
			Message: res.Payload,
		})
	case ResidentActionOpportunity:
		ra.sendEvent(domain.OpportunityEvent{
			EventBase: domain.EventBase{
				EventTo: domain.EventTo{ID: res.To},
				EventFrom: domain.EventFrom{
					ID:   domain.ActorID(ra.resident.ID),
					Name: domain.ActorName(ra.resident.Name),
				},
			},
			Opportunity: res.Payload,
		})
	case ResidentActionStay:
		return
	}
}

//go:embed schema/resident_agent_llm_response_schema.json
var residentLLMResponseSchema string

type ResidentActionKind string

const (
	ResidentActionMessage     ResidentActionKind = "message"
	ResidentActionOpportunity ResidentActionKind = "opportunity"
	ResidentActionStay        ResidentActionKind = "stay"
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
