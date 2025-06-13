package service

import (
	"log"

	"github.com/Murolando/m_rsm/internal/entitites"
)

type ServiceProposer struct {
	proposers []*entitites.Proposer
}

func (s *ServiceProposer) ProposeValue(proposerID int, value string) {
	for i := range 5 {
		log.Println(i, proposerID, value)
		lg, err := s.proposers[proposerID].ProposeValue(value)
		if err != nil {
			log.Println(err, value)
			continue
		}
		log.Println(lg)
		return
	}
}
