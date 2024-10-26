package statemachine

import (
	"slices"

	errs "github.com/H1ghN0on/go-tgbot-engine/errors"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
)

type Command string

type StateMachine struct {
	activeState State
	states      []State
}

type State struct {
	Name              string
	AvailableCommands []string
	AvailableStates   []State
}

func (state State) GetName() string {
	return state.Name
}

func (state State) GetAvailableCommands() []string {
	return state.AvailableCommands
}

func (state State) GetAvailableStates() []handlers.Stater {
	convertedStates := make([]handlers.Stater, len(state.AvailableStates))

	for _, state := range state.AvailableStates {
		// Assuming Stater has a method to convert to State
		convertedStates = append(convertedStates, state)
	}
	return convertedStates
}

func (sm *StateMachine) CompareStates(a State) func(State) bool {
	return func(b State) bool {
		return a.GetName() == b.GetName()
	}
}

func (sm *StateMachine) AddStates(states ...handlers.Stater) {
	for _, state := range states {
		// Assuming Stater has a method to convert to State
		sm.states = append(sm.states, state.(State))
	}
}

func (sm *StateMachine) SetStateByName(stateName string) error {
	idx := slices.IndexFunc(sm.states, func(s State) bool {
		return s.GetName() == stateName
	})
	if idx == -1 {
		return errs.StateMachineError{Code: errs.WrongState, Message: "This state is not unavailable"}
	}

	if sm.activeState.GetName() == "" ||
		sm.activeState.GetName() == stateName ||
		slices.ContainsFunc(sm.activeState.AvailableStates, sm.CompareStates(sm.states[idx])) {
		sm.activeState = sm.states[idx]
		return nil
	}

	return errs.StateMachineError{Code: errs.WrongState, Message: "Can not move to this state"}
}

func (sm *StateMachine) SetState(state handlers.Stater) error {
	if !slices.ContainsFunc(sm.states, sm.CompareStates(state.(State))) {
		return errs.StateMachineError{Code: errs.WrongState, Message: "This state is not unavailable"}
	}

	if sm.activeState.GetName() == "" ||
		sm.activeState.GetName() == state.GetName() ||
		slices.ContainsFunc(sm.activeState.AvailableStates, sm.CompareStates(state.(State))) {
		sm.activeState = state.(State)
		return nil
	}
	return errs.StateMachineError{Code: errs.WrongState, Message: "Can not move to this state"}
}

func (sm *StateMachine) GetActiveState() handlers.Stater {
	return sm.activeState
}
