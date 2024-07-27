package service

import (
	"PandoraFuclaudePlusHelper/internal/repository"
)

type Coordinator struct {
	OpenaiTokenSvc   OpenaiTokenService
	OpenaiAccountSvc OpenaiAccountService
	ClaudeTokenSvc   ClaudeTokenService
	ClaudeAccountSvc ClaudeAccountService
}

func NewServiceCoordinator(service *Service,
	openaiTokenRepository repository.OpenaiTokenRepository, openaiAccountRepository repository.OpenaiAccountRepository,
	claudeTokenRepository repository.ClaudeTokenRepository, claudeAccountRepository repository.ClaudeAccountRepository) *Coordinator {
	coordinator := &Coordinator{}

	openaiTokenSvc := NewOpenaiTokenService(service, openaiTokenRepository, openaiAccountRepository, coordinator)
	openaiAccountSvc := NewOpenaiAccountService(service, openaiTokenRepository, openaiAccountRepository, coordinator)

	claudeTokenSvc := NewClaudeTokenService(service, claudeTokenRepository, claudeAccountRepository, coordinator)
	claudeAccountSvc := NewClaudeAccountService(service, claudeTokenRepository, claudeAccountRepository, coordinator)

	coordinator.OpenaiTokenSvc = openaiTokenSvc
	coordinator.OpenaiAccountSvc = openaiAccountSvc
	coordinator.ClaudeTokenSvc = claudeTokenSvc
	coordinator.ClaudeAccountSvc = claudeAccountSvc

	return coordinator
}
