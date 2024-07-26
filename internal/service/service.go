package service

import (
	"PandoraPlusHelper/internal/repository"
	"PandoraPlusHelper/pkg/jwt"
	"PandoraPlusHelper/pkg/log"
	"PandoraPlusHelper/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(tm repository.Transaction, logger *log.Logger, sid *sid.Sid, jwt *jwt.JWT) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}
