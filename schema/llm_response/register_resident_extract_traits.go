package llmresponse

import (
	_ "embed"
	"encoding/json"
)

//go:embed register_resident_extract_traits.json
var registerResidentExtractTraits []byte

var RegisterResidentExtractTraits = json.RawMessage(registerResidentExtractTraits)

//go:embed resident_agent_process.json
var residentAgentProcess []byte

var ResidentAgentProcess = json.RawMessage(residentAgentProcess)
