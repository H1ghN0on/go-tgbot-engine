package statemachine

import (
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/handlers"
)

type Command string

type StateMachine struct {
	activeState   State
	previousState State
	states        []State
}

type State struct {
	name              string
	startCommand      string
	availableCommands []string
	availableStates   []State
}

type StateMachineError struct {
	message string
}

func (err StateMachineError) Error() string {
	return err.message
}

func (state State) GetStartCommand() string {
	return state.startCommand
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

func NewState(name string, startCommand string, availableCommands ...string) *State {
	return &State{
		name:              name,
		startCommand:      startCommand,
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
	err := sm.SetState(State{name: stateName})
	if err != nil {
		return StateMachineError{message: err.Error()}
	}
	return nil
}

func (sm *StateMachine) SetState(state handlers.Stater) error {
	if state.GetName() == "" {
		return StateMachineError{message: "State has empty name"}
	}

	idx := slices.IndexFunc(sm.states, sm.CompareStates(state.(State)))
	if idx == -1 {
		return StateMachineError{message: "This state is not unavailable"}
	}

	if sm.activeState.GetName() == "" ||
		sm.activeState.GetName() == state.GetName() ||
		slices.ContainsFunc(sm.activeState.availableStates, sm.CompareStates(sm.states[idx])) {
		sm.previousState = sm.activeState
		sm.activeState = sm.states[idx]
		return nil
	}
	return StateMachineError{message: "Can not move to this state"}
}

func (sm *StateMachine) GetActiveState() handlers.Stater {
	return sm.activeState
}

func (sm *StateMachine) GetPreviousState() handlers.Stater {
	return sm.previousState
}
