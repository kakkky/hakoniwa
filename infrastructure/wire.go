//go:build wireinject
// +build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/infrastructure/agent"
	"github.com/kakkky/hakoniwa/infrastructure/file"
	"github.com/kakkky/hakoniwa/infrastructure/llm"
)

var Set = wire.NewSet(
	file.NewFilePaths,
	file.NewFileResidentRepository,
	wire.Bind(new(domain.ResidentRepository), new(*file.FileResidentRepository)),
	llm.NewLLMGeminiProvider,
)
