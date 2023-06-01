package main

import (
	"fmt"
	"sync"
)

type TransactionManager struct {
	participants []*Participant
}

type Participant struct {
	ID int
	state string
	committed bool
}

func NewTransactionManager(participants int) *TransactionManager {
	tm := &TransactionManager{}
	for i := 0; i <  participants; i++ {
		tm.participants = append(tm.participants, &Participant{ID: i})
	}
	return tm
}

func (tm *TransactionManager) Phase1() bool {
	fmt.Println("Phase 1: Prepare")
	agreed := true
	for _, p :=range tm.participants {
		if p.committed == false {
			p.state = "agreed"
		} else {
			agreed = false
		}
	}
	return agreed
}

func (tm *TransactionManager) Phase2() {
	fmt.Println("Phase 2: Commit")
	var wg sync.WaitGroup
	wg.Add(len(tm.participants))
	for _, p := range tm.participants {
		go func(p *Participant) {
			defer wg.Done()
			if p.state == "agreed" {
				p.committed = true
				fmt.Printf("Participant %d committed\n", p.ID)
			}
		}(p)
	}
	wg.Wait()
}

func (tm *TransactionManager) IsCommited() bool {
	for _, p := range tm.participants {
		if !p.committed {
			return false
		}
	}
	return true
}

func main() {
	tm := NewTransactionManager(3)
	if tm.Phase1() {
		tm.Phase2()
	}
	if tm.IsCommited() {
		fmt.Println("All participants committed.")
	} else {
		fmt.Println("Not all participants could commit.")
	}
}