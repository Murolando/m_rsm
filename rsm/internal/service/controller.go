package service

import "github.com/Murolando/m_rsm/internal/entitites"

// сontroller interface
type Cotroller interface {
	Acceptor
	Proposer
}

type Acceptor interface {
	Promise()
	Accepted()
}
type Proposer interface {
	Prepare(value string)
	Propose()
	GenerateNewBallotNumber() *entitites.BallotNumber
}
