package statemachine

import (
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	"github.com/H1ghN0on/go-tgbot-engine/logger"

	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type Command string

type StateMachine struct {
	activeState   State
	previousState State
	states        []State
}

type State struct {
	name              string
	startCommand      bottypes.Command
	availableCommands []bottypes.Command
	availableStates   []State
	canRestart        bool
}

type StateMachineError struct {
	message string
}

func (err StateMachineError) Error() string {
	return err.message
}

func (state State) GetStartCommand() bottypes.Command {
	return state.startCommand
}

func (state State) GetName() string {
	return state.name
}

func (state State) GetAvailableCommands() []bottypes.Command {
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

func (state State) CanRestart() bool {
	return state.canRestart
}

func NewState(name string, startCommand bottypes.Command, availableCommands ...bottypes.Command) *State {
	return &State{
		name:              name,
		startCommand:      startCommand,
		availableCommands: availableCommands,
		canRestart:        false,
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

func (sm *StateMachine) SetState(state handlers.Stater) error {
	if state.GetName() == "" {
		logger.StateMachine().Critical("state has empty name")
		return StateMachineError{message: "State has empty name"}
	}

	idx := slices.IndexFunc(sm.states, sm.CompareStates(state.(State)))
	if idx == -1 {
		logger.StateMachine().Critical("state", sm.activeState.name)
		return StateMachineError{message: "This state is not unavailable"}
	}

	if sm.activeState.GetName() == "" || slices.ContainsFunc(sm.activeState.availableStates, sm.CompareStates(sm.states[idx])) {
		sm.previousState = sm.activeState
		sm.activeState = sm.states[idx]

		logger.StateMachine().Info("state updated:", sm.activeState.name)

		return nil
	}
	return StateMachineError{message: "Can not move to this state"}
}

func (sm *StateMachine) SetStateByName(stateName string) error {
	err := sm.SetState(State{name: stateName})
	if err != nil {
		return StateMachineError{message: err.Error()}
	}
	return nil
}

func (sm *StateMachine) GetPreviousState() handlers.Stater {
	return sm.previousState
}

func (sm *StateMachine) GetActiveState() handlers.Stater {
	return sm.activeState
}

func NewStateMachine() *StateMachine {

	startState := NewState(
		"start-state",

		cmd.ShowCommandsCommand,

		cmd.LevelOneCommand,
		cmd.LevelTwoCommand,
		cmd.LevelThreeCommand,
		cmd.ShowCommandsCommand,
		cmd.KeyboardStartCommand,
		cmd.LevelFourStartCommand,
		cmd.BigMessagesCommand,
		cmd.SetInfoStartCommand,
		cmd.CheckboxStartCommand,
		cmd.DynamicKeyboardStartCommand,
	)

	levelFourState := NewState(
		"level-four-state",

		cmd.LevelFourStartCommand,

		cmd.LevelFourStartCommand,
		cmd.LevelFourOneCommand,
		cmd.LevelFourTwoCommand,
		cmd.LevelFourThreeCommand,
		cmd.LevelFourFourCommand,
		cmd.BackStateCommand,
	)

	keyboardState := NewState(
		"keyboard-state",

		cmd.KeyboardStartCommand,

		cmd.KeyboardStartCommand,
		cmd.KeyboardOneCommand,
		cmd.KeyboardTwoCommand,
		cmd.KeyboardThreeCommand,
		cmd.KeyboardFinishCommand,
		cmd.BackStateCommand,
		cmd.BackCommandCommand,
	)

	infoState := NewState(
		"info-state",

		cmd.SetInfoStartCommand,

		cmd.SetInfoStartCommand,
		cmd.SetNameCommand,
		cmd.SetSurnameCommand,
		cmd.SetAgeCommand,
		cmd.SetInfoEndCommand,
		cmd.BackStateCommand,
		cmd.BackCommandCommand,
	)

	checkboxState := NewState(
		"checkbox-state",

		cmd.CheckboxStartCommand,

		cmd.CheckboxStartCommand,
		cmd.CheckboxFirstCommand,
		cmd.CheckboxSecondCommand,
		cmd.CheckboxThirdCommand,
		cmd.CheckboxFourthCommand,
		cmd.CheckboxAcceptCommand,
		cmd.BackStateCommand,
		cmd.NothingnessCommand,
	)

	dynamicKeyboardState := NewState(
		"dynamic-keyboard-state",

		cmd.DynamicKeyboardStartCommand,

		cmd.DynamicKeyboardFirstStageCommand,
		cmd.DynamicKeyboardSecondStageCommand,
		cmd.DynamicKeyboardFinishCommand,
		cmd.BackStateCommand,
		cmd.BackCommandCommand,
	)

	startState.SetAvailableStates(*levelFourState, *keyboardState, *infoState, *checkboxState, *startState, *dynamicKeyboardState)
	levelFourState.SetAvailableStates(*startState)
	keyboardState.SetAvailableStates(*startState)
	infoState.SetAvailableStates((*startState))
	checkboxState.SetAvailableStates((*startState))
	dynamicKeyboardState.SetAvailableStates((*startState))

	sm := &StateMachine{}

	sm.AddStates(*startState, *levelFourState, *keyboardState, *infoState, *checkboxState, *dynamicKeyboardState)

	err := sm.SetStateByName("start-state")
	if err != nil {
		panic(err.Error())
	}

	return sm
}
