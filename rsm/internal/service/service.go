package service

import "github.com/Murolando/m_rsm/internal/entitites"

type Service struct {

	ServiceProposer
}

func New(proposers []*entitites.Proposer) *Service {
	return &Service{
		ServiceProposer: ServiceProposer{
			proposers: proposers,
		},
	}
}
