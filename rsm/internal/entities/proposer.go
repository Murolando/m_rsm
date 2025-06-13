package entitites

import (
	"errors"
	"log"
	"sync"
	"time"
)

const (
	FalsePromise               = "FALSE_PROMISE"
	DefaultValue               = "DEFAULT_VALUE"
	ProposerWaiterExpiredError = "proposer Waiter Exipered"
	FailedToProposeNewValue    = "failed to propose new value"
)

type BallotNumber struct {
	Number     int64
	ProposerID int
}

// PROPOSER
type Proposer struct {
	ProposerID int
	// mutex
	Mu sync.Mutex
	// Last ballot number
	LastBallotNumber *BallotNumber
	// Acceptors
	acceptors []*Acceptor

	// Quorum timer
	waitTimer time.Duration
}

func NewProposer(proposerID int, acceptors []*Acceptor, waitTimer time.Duration) *Proposer {
	return &Proposer{
		ProposerID:       proposerID,
		LastBallotNumber: &BallotNumber{},
		acceptors:        acceptors,
		waitTimer:        waitTimer,
	}
}

type PrepareMessage struct {
	// n
	ProposerBallotNumber *BallotNumber
}

type ProposeMessage struct {
	// (n_accepted,v_accepted)
	NewAcceptedBallotNumber *BallotNumber
	// log
	NewAcceptedLog []Log
}

func (p *Proposer) GenerateNewBallotNumber() {
	p.LastBallotNumber.Number += 1
}

func (p *Proposer) ProposeValue(message string) ([]Log, error) {
	// generate new ballot number
	p.GenerateNewBallotNumber()
	pNum := *p.LastBallotNumber

	prepareMessage := PrepareMessage{
		ProposerBallotNumber: &pNum,
	}

	// prepare
	promises, promise, err := p.Prepare(&prepareMessage)
	if err != nil {
		return nil, err
	}
	if !promise {
		p.LastBallotNumber.Number = promises[0].MaxPromissedBallotNumber.Number
		return nil, errors.New(FailedToProposeNewValue)
	}

	value := &Log{
		LastAcceptedBallotNumber: p.LastBallotNumber,
		LastAcceptedValue:        message,
	}
	newLog := mergeLists(promises, *p.LastBallotNumber, value)

	// propose
	proposeMessage := ProposeMessage{
		NewAcceptedBallotNumber: &pNum,
		NewAcceptedLog:          newLog,
	}
	_, accept, err := p.Propose(&proposeMessage)
	if err != nil {
		return nil, err
	}
	if !accept {
		p.LastBallotNumber.Number = promises[0].MaxPromissedBallotNumber.Number
		return nil, errors.New(FailedToProposeNewValue)
	}
	log.Println(p.acceptors[0].Logs, p.acceptors[1].Logs, p.acceptors[2].Logs)
	return p.acceptors[0].Logs, nil
}

func (p *Proposer) Prepare(prepareMessage *PrepareMessage) ([]*PromiseMessage, bool, error) {
	result := make([]*PromiseMessage, 0)
	var resultMu sync.Mutex

	// Get message from acceptors and parallelly send message to acceptors
	for _, a := range p.acceptors {
		go func() {
			r := a.Promise(prepareMessage)
			resultMu.Lock()
			result = append(result, r)
			resultMu.Unlock()
		}()
	}

	// Wait time and check quorum, if it is not there or it is false return error
	time.Sleep(p.waitTimer * time.Second)

	p.Mu.Lock()
	defer p.Mu.Unlock()

	if len(result) < len(p.acceptors)/2+1 {
		return result, false, errors.New(ProposerWaiterExpiredError)
	}
	for _, promise := range result {
		if !promise.Promise {
			result = []*PromiseMessage{promise}
			return result, false, nil
		}
	}

	return result, true, nil
}

func (p *Proposer) Propose(proposeMessage *ProposeMessage) ([]*AcceptedMessage, bool, error) {
	result := make([]*AcceptedMessage, 0)
	var resultMu sync.Mutex

	// Get message from acceptors and parallelly send message to acceptors
	for _, a := range p.acceptors {
		go func() {
			r := a.Accept(proposeMessage)
			resultMu.Lock()
			result = append(result, r)
			resultMu.Unlock()
		}()
	}
	// Wait time and check quorum, if it is not there or it is false return error
	time.Sleep(p.waitTimer * time.Second)
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if len(result) < len(p.acceptors)/2+1 {
		return result, false, errors.New(ProposerWaiterExpiredError)
	}
	for _, accept := range result {
		if *accept.NewAcceptedBallotNumber != *proposeMessage.NewAcceptedBallotNumber {
			result = []*AcceptedMessage{accept}
			return result, false, nil
		}
	}

	return result, true, nil
}

func mergeLists(promises []*PromiseMessage, proposerNumber BallotNumber, proposedValue *Log) []Log {
	maxLength := 0
	for _, promise := range promises {
		if len(promise.Logs) > maxLength {
			maxLength = len(promise.Logs)
		}
	}

	result := make([]Log, maxLength)
	added := false
	for i := range maxLength {
		var currentMax *Log
		for _, promise := range promises {
			if i < len(promise.Logs) {
				if currentMax == nil || promise.Logs[i].LastAcceptedBallotNumber.Number > currentMax.LastAcceptedBallotNumber.Number {
					currentMax = &promise.Logs[i]
				}
			}
		}
		if currentMax != nil {
			*currentMax.LastAcceptedBallotNumber = proposerNumber
		} else {
			currentMax = proposedValue
			added = true
		}
		result[i] = *currentMax
	}
	if !added {
		result = append(result, *proposedValue)
	}

	return result
}
