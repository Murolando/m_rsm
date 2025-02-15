package entitites

import (
	"fmt"
	"sync"
)

const (
	DefaultValue = "DEFAULT_VALUE"
	FalsePromise = "FALSE_PROMISE"
)

type BallotNumber struct {
	Number     int64
	ProposerID int
}

// PROPOSER
type Proposer struct {
	ProposerID int
	// Последни принятый номер
	LastBallotNumber *BallotNumber
	// Список асепторов
	acceptors []*Acceptor
}

func NewProposer(proposerID int, acceptors []*Acceptor) *Proposer {
	return &Proposer{
		ProposerID:       proposerID,
		LastBallotNumber: &BallotNumber{},
		acceptors:        acceptors,
	}
}

type PrepareMessage struct {
	// n
	ProposerBallotNumber *BallotNumber
}

type ProposeMessage struct {
	// (n_accepted,v_accepted)
	NewAcceptedBallotNumer *BallotNumber
	NewAcceptedValue       string
}

func (p *Proposer) GenerateNewBallotNumber() {
	p.LastBallotNumber.Number += 1
}

func (p *Proposer) ProposeValue(value string) bool {
	fmt.Println("Phase 1")
	// Phase 1
	// Prepare
	// generate ballot number
	p.GenerateNewBallotNumber()
	// prepare message
	prepareMessage := &PrepareMessage{
		ProposerBallotNumber: p.LastBallotNumber,
	}
	// prepare
	fmt.Println("Phase 1, Prepare")
	promises := p.Prepare(prepareMessage)
	// check quorum
	if promises == nil || len(promises) < len(p.acceptors)/2+1 {
		return false
	}
	fmt.Println("Phase 1, Choose", fmt.Sprintln(promises))
	// choose value
	proposeMessage := &ProposeMessage{
		NewAcceptedBallotNumer: p.LastBallotNumber,
		NewAcceptedValue:       value,
	}

	var number int64
	for _, msg := range promises {
		if msg.LastAcceptedBallotNumer.Number > number && msg.LastAcceptedValue != DefaultValue{
			number = msg.LastAcceptedBallotNumer.Number
			proposeMessage.NewAcceptedValue = msg.LastAcceptedValue
		}
	}
	fmt.Println("Phase 2")
	// Phase 2

	// Propose
	fmt.Println("Phase 2, Propose")
	acceptedMessages := p.Propose(proposeMessage)
	if acceptedMessages == nil || len(acceptedMessages) < len(p.acceptors)/2+1 {
		fmt.Println("CONSESUS NOT FOUND", p.ProposerID, proposeMessage.NewAcceptedValue)
		return false
	}
	fmt.Println("CONSENSUS!", p.ProposerID, proposeMessage.NewAcceptedValue)
	return true
}

func (p *Proposer) Prepare(prepareMessage *PrepareMessage) []*PromiseMessage {
	quorumMu := sync.Mutex{}
	promises := make([]*PromiseMessage, 0, len(p.acceptors))

	for i, acceptor := range p.acceptors {
		fmt.Println("Phase 1, Prepare, Promise", i, acceptor)
		go func() {
			msg := acceptor.Promise(prepareMessage)
			if msg.LastAcceptedValue != FalsePromise {
				fmt.Println("Phase 1", i, msg)
				quorumMu.Lock()
				promises = append(promises, msg)
				quorumMu.Unlock()
			}
		}()
	}
	Wait()
	Wait()
	fmt.Println("Phase 1, Prepare, Check", promises)
	if len(promises) >= len(p.acceptors)/2+1 {
		return promises
	}
	return nil
}
func (p *Proposer) Propose(proposeMessage *ProposeMessage) []*AcceptedMessage {
	quorumMu := sync.Mutex{}
	acceptedMessages := make([]*AcceptedMessage, 0, len(p.acceptors))

	for i, acceptor := range p.acceptors {
		go func() {
			fmt.Println("Phase 2, Propose, Accept", i, acceptor, acceptor.MaxPromissedBallotNumer, proposeMessage.NewAcceptedValue)
			msg := acceptor.Accept(proposeMessage)
			if msg.NewAcceptedBallotNumer != nil {
				fmt.Println("Phase 2", i, msg)
				quorumMu.Lock()
				acceptedMessages = append(acceptedMessages, msg)
				quorumMu.Unlock()
			}
		}()
	}
	Wait()
	Wait()
	fmt.Println("Phase 2, Propose, check", acceptedMessages)
	if len(acceptedMessages) >= len(p.acceptors)/2+1 {
		return acceptedMessages
	}
	return nil
}
