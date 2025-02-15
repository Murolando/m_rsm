package entitites

import (
	"fmt"
	"sync"
	"time"

	"math/rand"
)

// ACCEPTOR
type Acceptor struct {
	// id
	AcceptorID int

	// mutex
	Mu sync.Mutex
	// n_promissed
	MaxPromissedBallotNumer *BallotNumber
	// (n_accepted,v_accepted)
	LastAcceptedBallotNumer *BallotNumber
	LastAcceptedValue       string
	// logs
	Logs []string
}

func NewAcceptor(acceptorID int) *Acceptor {
	return &Acceptor{
		AcceptorID:              acceptorID,
		Mu:                      sync.Mutex{},
		MaxPromissedBallotNumer: &BallotNumber{},
		LastAcceptedBallotNumer: &BallotNumber{},
		LastAcceptedValue:       DefaultValue,
		Logs:                    make([]string, 0),
	}
}

type PromiseMessage struct {
	// n
	ProposerBallotNumber *BallotNumber
	// (n_accepted,v_accepted) || (null,-1)
	LastAcceptedBallotNumer *BallotNumber
	LastAcceptedValue       string
}

type AcceptedMessage struct {
	// (n_accepted)
	NewAcceptedBallotNumer *BallotNumber
}

func (a *Acceptor) Promise(prepareMessage *PrepareMessage) *PromiseMessage {
	Wait()
	fmt.Println("Phase 1, Promise, Acceptor", a.AcceptorID, prepareMessage.ProposerBallotNumber.Number, prepareMessage.ProposerBallotNumber.ProposerID)
	result := &PromiseMessage{
		ProposerBallotNumber:    a.MaxPromissedBallotNumer,
		LastAcceptedBallotNumer: nil,
		LastAcceptedValue:       FalsePromise,
	}
	a.Mu.Lock()
	defer a.Mu.Unlock()

	if prepareMessage.ProposerBallotNumber.Number > a.MaxPromissedBallotNumer.Number {
		// promise
		a.MaxPromissedBallotNumer.Number = prepareMessage.ProposerBallotNumber.Number

		// result
		result.ProposerBallotNumber = prepareMessage.ProposerBallotNumber
		result.LastAcceptedBallotNumer = a.LastAcceptedBallotNumer
		result.LastAcceptedValue = a.LastAcceptedValue
	}
	fmt.Println("Phase 1, Prepare, Promise, Unlock", a.AcceptorID, *result.ProposerBallotNumber)
	return result
}

func (a *Acceptor) Accept(acceptMessage *ProposeMessage) *AcceptedMessage {
	Wait()

	result := &AcceptedMessage{
		NewAcceptedBallotNumer: nil,
	}

	a.Mu.Lock()
	defer a.Mu.Unlock()
	fmt.Println("Phase 2, Accept, Acceptor", a.AcceptorID)
	if a.MaxPromissedBallotNumer.Number >= acceptMessage.NewAcceptedBallotNumer.Number {
		// accepted
		a.LastAcceptedBallotNumer.Number = acceptMessage.NewAcceptedBallotNumer.Number
		a.LastAcceptedBallotNumer.ProposerID = acceptMessage.NewAcceptedBallotNumer.ProposerID
		//value
		a.LastAcceptedValue = acceptMessage.NewAcceptedValue
		// result
		result.NewAcceptedBallotNumer = a.LastAcceptedBallotNumer
	}
	fmt.Println("Phase 2, Accept, Unlock", a.AcceptorID,result)
	return result
}

func Wait() {
	pause := rand.Int() % 10
	time.Sleep(time.Second * time.Duration(pause))
}
