package main

import (
	"fmt"
	"sync"

	"github.com/Murolando/m_rsm/internal/entitites"
	"github.com/Murolando/m_rsm/internal/service"
)

func main() {
	fmt.Println(1)
	acceptors := []*entitites.Acceptor{
		entitites.NewAcceptor(1),
		entitites.NewAcceptor(2),
		entitites.NewAcceptor(3),
		entitites.NewAcceptor(4),
		entitites.NewAcceptor(5),
	}
	fmt.Println(acceptors)
	proposers := []*entitites.Proposer{
		entitites.NewProposer(0, acceptors),
		entitites.NewProposer(1, acceptors),
	}

	service := service.New(proposers)

	// TestOneProposer(service, acceptors)

	TestTwoParallelProposers(service, acceptors)
}

func TestOneProposer(service *service.Service, acceptors []*entitites.Acceptor) {
	service.ProposeValue("svetlana", 0)

	for _, val := range acceptors {
		fmt.Println(val.LastAcceptedValue, val.LastAcceptedBallotNumer)
	}
}

func TestTwoParallelProposers(service *service.Service, acceptors []*entitites.Acceptor) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		service.ProposeValue("svetlana", 0)
	}()

	go func() {
		defer wg.Done()
		service.ProposeValue("oleg", 1)
	}()

	wg.Wait()
	for _, val := range acceptors {
		fmt.Println(val.AcceptorID, val.MaxPromissedBallotNumer, val.LastAcceptedValue, val.LastAcceptedBallotNumer)
	}
}
