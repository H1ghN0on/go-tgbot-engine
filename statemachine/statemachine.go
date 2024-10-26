package statemachine

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	errs "github.com/H1ghN0on/go-tgbot-engine/errors"

	"slices"
)

type StateMachine struct {
	activeState bottypes.State
	states      []bottypes.State
}

func CompareStates(a bottypes.State) func(bottypes.State) bool {
	return func(b bottypes.State) bool {
		return a.Name == b.Name
	}
}

func (sm *StateMachine) AddStates(states ...bottypes.State) {
	sm.states = append(sm.states, states...)
}

func (sm *StateMachine) SetStateByName(stateName string) error {
	idx := slices.IndexFunc(sm.states, func(s bottypes.State) bool {
		return s.Name == stateName
	})
	if idx == -1 {
		return errs.StateMachineError{Code: errs.WrongState, Message: "This state is not unavailable"}
	}

	if sm.activeState.Name == "" ||
		sm.activeState.Name == stateName ||
		slices.ContainsFunc(sm.activeState.AvailableStates, CompareStates(sm.states[idx])) {
		sm.activeState = sm.states[idx]
		return nil
	}

	return errs.StateMachineError{Code: errs.WrongState, Message: "Can not move to this state"}
}

func (sm *StateMachine) SetState(state bottypes.State) error {
	if !slices.ContainsFunc(sm.states, CompareStates(state)) {
		return errs.StateMachineError{Code: errs.WrongState, Message: "This state is not unavailable"}
	}

	if sm.activeState.Name == "" ||
		sm.activeState.Name == state.Name ||
		slices.ContainsFunc(sm.activeState.AvailableStates, CompareStates(state)) {
		sm.activeState = state
		return nil
	}
	return errs.StateMachineError{Code: errs.WrongState, Message: "Can not move to this state"}
}

func (sm *StateMachine) GetActiveState() bottypes.State {
	return sm.activeState
}
