package statemachine

import (
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/errors"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
)

type Command string

type StateMachine struct {
	activeState State
	states      []State
}

type State struct {
	name              string
	availableCommands []string
	availableStates   []State
}

func (state State) GetName() string {
	return state.name
}

func (state State) GetAvailableCommands() []string {
	return state.availableCommands
}

func (state State) GetAvailableStates() []handlers.Stater {
	var convertedStates []handlers.Stater

	for _, state := range state.availableStates {
		convertedStates = append(convertedStates, state)
	}
	return convertedStates
}

func (state *State) SetAvailableStates(newStates ...handlers.Stater) {
	for _, newState := range newStates {
		state.availableStates = append(state.availableStates, newState.(State))
	}
}

func NewState(name string, availableCommands ...string) *State {
	return &State{
		name:              name,
		availableCommands: availableCommands,
	}
}

func (sm *StateMachine) CompareStates(a State) func(State) bool {
	return func(b State) bool {
		return a.GetName() == b.GetName()
	}
}

func (sm *StateMachine) AddStates(states ...handlers.Stater) {
	for _, state := range states {
		sm.states = append(sm.states, state.(State))
	}
}

func (sm *StateMachine) SetStateByName(stateName string) error {
	idx := slices.IndexFunc(sm.states, func(s State) bool {
		return s.GetName() == stateName
	})
	if idx == -1 {
		return errors.StateMachineError{Code: errors.WrongState, Message: "This state is not unavailable"}
	}

	if sm.activeState.GetName() == "" ||
		sm.activeState.GetName() == stateName ||
		slices.ContainsFunc(sm.activeState.availableStates, sm.CompareStates(sm.states[idx])) {
		sm.activeState = sm.states[idx]
		return nil
	}

	return errors.StateMachineError{Code: errors.WrongState, Message: "Can not move to this state"}
}

func (sm *StateMachine) SetState(state handlers.Stater) error {
	if !slices.ContainsFunc(sm.states, sm.CompareStates(state.(State))) {
		return errors.StateMachineError{Code: errors.WrongState, Message: "This state is not unavailable"}
	}

	if sm.activeState.GetName() == "" ||
		sm.activeState.GetName() == state.GetName() ||
		slices.ContainsFunc(sm.activeState.availableStates, sm.CompareStates(state.(State))) {
		sm.activeState = state.(State)
		return nil
	}
	return errors.StateMachineError{Code: errors.WrongState, Message: "Can not move to this state"}
}

func (sm *StateMachine) GetActiveState() handlers.Stater {
	return sm.activeState
}
