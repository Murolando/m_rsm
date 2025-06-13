package entitites

import (
	"sync"
)

// ACCEPTOR
type Acceptor struct {
	// id
	AcceptorID int

	// mutex
	Mu sync.Mutex
	// n_promissed
	MaxPromissedBallotNumber *BallotNumber
	// (n_accepted)
	LastAcceptedBallotNumber *BallotNumber
	// logs
	Logs []Log
}

func NewAcceptor(acceptorID int) *Acceptor {
	return &Acceptor{
		AcceptorID:               acceptorID,
		Mu:                       sync.Mutex{},
		MaxPromissedBallotNumber: &BallotNumber{},
		LastAcceptedBallotNumber: &BallotNumber{},
		Logs:                     make([]Log, 0),
	}
}

type PromiseMessage struct {
	// n_promissed
	MaxPromissedBallotNumber *BallotNumber
	// logs of this acceptor
	Logs []Log
	// promise
	Promise bool
}

type AcceptedMessage struct {
	// (n_accepted)
	NewAcceptedBallotNumber *BallotNumber
}

func (a *Acceptor) Promise(prepareMessage *PrepareMessage) *PromiseMessage {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	result := &PromiseMessage{
		MaxPromissedBallotNumber: &BallotNumber{},
	}

	if prepareMessage.ProposerBallotNumber.Number > a.MaxPromissedBallotNumber.Number {
		result.Promise = true

		a.MaxPromissedBallotNumber = prepareMessage.ProposerBallotNumber
		result.Logs = append(result.Logs, a.Logs...)
	}
	result.MaxPromissedBallotNumber = a.MaxPromissedBallotNumber

	return result
}

func (a *Acceptor) Accept(acceptMessage *ProposeMessage) *AcceptedMessage {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	result := &AcceptedMessage{
		NewAcceptedBallotNumber: &BallotNumber{},
	}
	if acceptMessage.NewAcceptedBallotNumber.Number >= a.MaxPromissedBallotNumber.Number {

		a.LastAcceptedBallotNumber = acceptMessage.NewAcceptedBallotNumber
		a.Logs = acceptMessage.NewAcceptedLog
	}
	result.NewAcceptedBallotNumber = a.LastAcceptedBallotNumber
	return result
}
