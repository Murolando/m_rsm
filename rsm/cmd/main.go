package main

import (
	"fmt"
	"sync"

	"github.com/Murolando/m_rsm/internal/entities"
	"github.com/Murolando/m_rsm/internal/service"
)

func main() {
	acceptors := []*entitites.Acceptor{
		entitites.NewAcceptor(1),
		entitites.NewAcceptor(2),
		entitites.NewAcceptor(3),
		entitites.NewAcceptor(4),
		entitites.NewAcceptor(5),
	}
	fmt.Println(acceptors)
	proposers := []*entitites.Proposer{
		entitites.NewProposer(0, acceptors, 1),
		entitites.NewProposer(1, acceptors, 1),
	}

	service := service.New(proposers)

	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		service.ProposeValue(0, "murolando")
		wg.Done()
	}()
	go func() {
		service.ProposeValue(1, "ninja")
		wg.Done()
	}()
	wg.Wait()

	wg.Add(3)
	go func() {
		service.ProposeValue(0, "alexyB")
		wg.Done()
	}()
	go func() {
		service.ProposeValue(1, "oleg")
		wg.Done()
	}()
	go func() {
		service.ProposeValue(1, "joomba")
		wg.Done()
	}()
	wg.Wait()
}
