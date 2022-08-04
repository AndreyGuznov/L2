package main

import "fmt"

/*
State поведенческий паттерн
позволяет управлять поведением объекта в зависимости от состояния.
*/

type Activity interface {
	someActivity()
}

type Catching struct{}

func (c *Catching) someActivity() {
	fmt.Println("Catching ball...")
}

type Walking struct{}

func (w *Walking) someActivity() {
	fmt.Println("Walking...")
}

type Working struct{}

func (w *Working) someActivity() {
	fmt.Println("Working...")
}

type SomeMan struct {
	activity Activity
}

func (s *SomeMan) setActivity(activity Activity) {
	s.activity = activity
}

func (s *SomeMan) changeActivity() {
	switch s.activity.(type) {
	case *Working:
		s.setActivity(&Walking{})
	case *Walking:
		s.setActivity(&Catching{})
	case *Catching:
		s.setActivity(&Working{})
	}
}

func (s *SomeMan) someActivity() {
	s.activity.someActivity()
}

func main() {
	activity := Working{}
	SomeMan := SomeMan{}
	SomeMan.setActivity(&activity)
	for i := 0; i < 5; i++ {
		SomeMan.someActivity()
		SomeMan.changeActivity()
	}
}
