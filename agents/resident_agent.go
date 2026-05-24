package agents

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kakkky/hakoniwa/domain"
	llmresponse "github.com/kakkky/hakoniwa/schema/llm_response"
)

type residentAgent struct {
	agentBase
	residentID domain.ResidentID
	repository domain.ResidentRepository
}

func newResidentAgent(base agentBase, repo domain.ResidentRepository, residentID domain.ResidentID) *residentAgent {
	return &residentAgent{
		residentID: residentID,
		agentBase:  base,
		repository: repo,
	}
}

type residentAgentProcessResponse struct {
	Action  residentActionKind `json:"action"`
	To      residentActionTo   `json:"to,omitempty"`
	Payload string             `json:"payload"`
	ToolUse *residentToolUse   `json:"tool_use,omitempty"`
}

type residentActionKind string

const (
	residentActionMessage     residentActionKind = "message"
	residentActionOpportunity residentActionKind = "opportunity"
	residentActionStay        residentActionKind = "stay"
)

type residentActionTo struct {
	ID   domain.ResidentID   `json:"id"`
	Name domain.ResidentName `json:"name"`
}

type residentToolUse struct {
	UpdateMood *updateMood `json:"update_mood,omitempty"`
}

type updateMood struct {
	Mood domain.Mood `json:"mood"`
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
	resident, err := ra.repository.FindByID(ra.residentID)
	if err != nil {
		return err
	}
	name := resident.Name
	age := resident.Age
	gender := resident.Gender
	traits := resident.Traits
	systemPromptTemplate :=
		`あなたはあるマンションの住民の一人です。
		プロフィール：
			- 名前：%s
			- 性別：%s
			- 年齢：%s
			- 特徴：%s
		あなたの記憶：
			- %s
		あなたの気分：
			- %s
		`
	memories := resident.Memories
	mood := resident.Mood
	systemPrompt := fmt.Sprintf(systemPromptTemplate, name, gender, age, traits, memories, mood)
	ra.agentBase.llmPrompt.AddSystemPrompt(systemPrompt)

	switch eventV := event.(type) {
	case domain.MessageEvent:
		newMemoryContent := fmt.Sprintf(`%sから「%s」とメッセージを受け取った。`, eventV.EventFrom.Name, eventV.Message)
		newMemory := domain.Memory{
			Content:   newMemoryContent,
			OccuredAt: time.Now(),
		}
		resident.AddMemory(newMemory)
		return ra.processMessageEvent(ctx, eventV)
	case domain.OpportunityEvent:
		newMemoryContent := fmt.Sprintf(`%sに%s。`, eventV.EventFrom.Name, eventV.Opportunity)
		newMemory := domain.Memory{
			Content:   newMemoryContent,
			OccuredAt: time.Now(),
		}
		resident.AddMemory(newMemory)
		return ra.processOppotunityeEvent(ctx, eventV)

	}

	defer ra.repository.Save(resident)

	return nil
}

func (ra *residentAgent) processMessageEvent(ctx context.Context, msgEvent domain.MessageEvent) error {
	msg := msgEvent.Payload()
	userPromptTemplate :=
		``
	userPrompt := fmt.Sprintf(userPromptTemplate, msg)
	ra.llmPrompt.AddUserPrompt(userPrompt)

	rawResp, err := ra.llmProvider.Generate(ctx, ra.llmPrompt, llmresponse.ResidentAgentProcess)
	if err != nil {
		return err
	}
	var residentAgentProcessResponse residentAgentProcessResponse
	if err := json.Unmarshal(rawResp, &residentAgentProcessResponse); err != nil {
		return nil
	}
	if err := ra.handleResidentLLMResponse(residentAgentProcessResponse); err != nil {
		return err
	}

	return nil
}

func (ra *residentAgent) processOppotunityeEvent(ctx context.Context, oppEvent domain.OpportunityEvent) error {
	msg := oppEvent.Payload()

	systemPromptTemplate :=
		``
	systemPrompt := fmt.Sprintf(systemPromptTemplate, nil)
	userPromptTemplate := ``

	userPrompt := fmt.Sprintf(userPromptTemplate, msg)

	ra.llmPrompt.AddSystemPrompt(systemPrompt)
	ra.llmPrompt.AddUserPrompt(userPrompt)

	rawResp, err := ra.llmProvider.Generate(ctx, ra.llmPrompt, llmresponse.ResidentAgentProcess)
	if err != nil {
		return err
	}
	var residentAgentProcessResponse residentAgentProcessResponse
	if err := json.Unmarshal(rawResp, &residentAgentProcessResponse); err != nil {
		return err
	}
	if err := ra.handleResidentLLMResponse(residentAgentProcessResponse); err != nil {
		return err
	}

	return nil
}

func (ra *residentAgent) handleResidentLLMResponse(res residentAgentProcessResponse) error {

	toolUse := res.ToolUse
	if toolUse != nil {
		switch {
		case toolUse.UpdateMood != nil:
			ra.resident.UpdateMood(toolUse.UpdateMood.Mood)
			if err := ra.repository.Save(ra.resident); err != nil {
				return err
			}
		}
	}

	action := res.Action
	switch action {
	case residentActionMessage:
		ra.sendEvent(domain.MessageEvent{
			EventBase: domain.EventBase{
				EventTo: domain.EventTo{
					ID:   res.To.ID,
					Name: res.To.Name,
				},
				EventFrom: domain.EventFrom{
					ID:   domain.ActorID(ra.resident.ID),
					Name: domain.ActorName(ra.resident.Name),
				},
			},
			Message: res.Payload,
		})
		newMemoryContent := fmt.Sprintf(`%sに「%s」とメッセージを受け取った。`, res.To.Name, res.Payload)
		newMemory := domain.Memory{
			Content:   newMemoryContent,
			OccuredAt: time.Now(),
		}
		ra.resident.AddMemory(newMemory)
	case residentActionOpportunity:
		ra.sendEvent(domain.OpportunityEvent{
			EventBase: domain.EventBase{
				EventTo: domain.EventTo{
					ID:   res.To.ID,
					Name: res.To.Name,
				},
				EventFrom: domain.EventFrom{
					ID:   domain.ActorID(ra.resident.ID),
					Name: domain.ActorName(ra.resident.Name),
				},
			},
			Opportunity: res.Payload,
		})
		newMemoryContent := fmt.Sprintf(`%sに%s`, res.To.Name, res.Payload)
		newMemory := domain.Memory{
			Content:   newMemoryContent,
			OccuredAt: time.Now(),
		}
		ra.resident.AddMemory(newMemory)
	case residentActionStay:
	}
	return nil
}
