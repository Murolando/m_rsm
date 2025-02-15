package service

import (
	"fmt"

	"github.com/Murolando/m_rsm/internal/entitites"
)

type ServiceProposer struct {
	proposers []*entitites.Proposer
}

func (s *ServiceProposer) ProposeValue(value string, proposerID int) bool {
	fmt.Println("Service")
	for i := range 5 {
		fmt.Println("Iteration:", i+1)
		if res := s.proposers[proposerID].ProposeValue(value); res {
			return res
		}
		s.proposers[proposerID].GenerateNewBallotNumber()
	}
	return false
}
