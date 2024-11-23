package statemachine

import (
	"testing"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/stretchr/testify/assert"
)

func TestStateMachineError_Error(t *testing.T) {

	expectedMessage := "something went wrong"

	err := StateMachineError{message: expectedMessage}

	assert.Equal(t, expectedMessage, err.Error(),
		"Error() = %q; want %q",
		err.Error(), expectedMessage)

}
func TestState_GetStartCommand(t *testing.T) {

	state := State{
		name:              "Init",
		startCommand:      bottypes.Command{Command: "start"},
		availableCommands: []bottypes.Command{{Command: "run"}, {Command: "stop"}},
		availableStates:   []State{},
	}

	expectedCommand := bottypes.Command{Command: "start"}

	assert.Equal(t, expectedCommand, state.GetStartCommand(),
		"GetStartCommand(): got %q, want %q",
		state.GetStartCommand(), expectedCommand)

}
func TestState_GetName(t *testing.T) {
	state := State{name: "Jambo"}

	assert.Equal(t, "Jambo", state.GetName(),
		"GetName(): got %q, want %q",
		state.GetName(), "Jambo")

}

func TestState_GetAvailableCommands(t *testing.T) {
	state := State{
		availableCommands: []bottypes.Command{{Command: "run"}, {Command: "stop"}},
	}

	expectedCommands := []bottypes.Command{{Command: "run"}, {Command: "stop"}}

	assert.Len(t, state.availableCommands, len(expectedCommands),
		"GetAvailableCommands(): len state.AvailableCommands = %d; expected %d",
		len(state.availableCommands), len(expectedCommands))

	for index := range expectedCommands {
		assert.Equal(t, expectedCommands[index], state.availableCommands[index],
			"GetAvailableCommands(): at index %d, got %q, expected %q",
			index, state.availableCommands[index], expectedCommands[index])
	}

}

func TestState_GetAvailableStates(t *testing.T) {
	state1 := State{name: "Jambo💀"}
	state2 := State{name: "Igor🌈"}
	mainState := State{availableStates: []State{state1, state2}}
	expectedState := State{availableStates: []State{state1, state2}}

	assert.Equal(t, expectedState.availableStates, mainState.availableStates,
		"GetAvailableStates(): states not equal,\n got %v, expected %v",
		mainState.availableStates, expectedState.availableStates)

}

func TestState_SetAvailableStates(t *testing.T) {
	state1 := State{name: "Stater1"}
	state2 := State{name: "Stater2"}

	mainState := &State{}
	mainState.SetAvailableStates(state1, state2)
	assert.Len(t, mainState.availableStates, 2)
	assert.Equal(t, mainState.availableStates[0].name, "Stater1")
	assert.Equal(t, mainState.availableStates[1].name, "Stater2")

}

func TestState_NewState(t *testing.T) {
	receivedName := "Name"
	receivedStartCommand := bottypes.Command{Command: "Run"}
	receivedAvailableCommands := []bottypes.Command{{Command: "Run"}, {Command: "Stop"}}

	expectedState := &State{
		name:              "Name",
		startCommand:      bottypes.Command{Command: "Run"},
		availableCommands: []bottypes.Command{{Command: "Run"}, {Command: "Stop"}},
	}

	testState := NewState(receivedName, receivedStartCommand, receivedAvailableCommands...)

	assert.Equal(t, expectedState, testState,
		"NewState(): states not equal,\n got %v, expected %v",
		testState, expectedState)

}

func TestStateMachine_AddStates(t *testing.T) {
	state1, state2, state3 := State{}, State{}, State{}
	expectedStateMachine := &StateMachine{
		states: []State{state1, state2, state3},
	}
	testStateMachine := &StateMachine{}
	testStateMachine.AddStates(state1, state2, state3)

	assert.Equal(t, expectedStateMachine.states, testStateMachine.states,
		"AddState(): states not equal,\n got %v, expected %v",
		testStateMachine.states, expectedStateMachine.states)

}

func TestStateMachine_SetState(t *testing.T) {

	state2 := State{name: "Middle"}
	state3 := State{name: "End"}

	state1 := State{
		name:            "Start",
		availableStates: []State{state2},
	}

	sm := &StateMachine{
		states: []State{state1, state2, state3},
	}

	// Test 1: Setting status with empty name
	err := sm.SetState(State{name: ""})
	assert.NotNil(t, err, "expected error 'State has empty name', but got nil")
	assert.Equal(t, "State has empty name", err.Error(), "unexpected error message: got %v", err.Error())

	// Test 2: Transition to state that is not in list
	err = sm.SetState(State{name: "NonExistent"})
	assert.NotNil(t, err, "expected error 'This state is not unavailable', but got nil")
	assert.Equal(t, "This state is not unavailable", err.Error(), "unexpected error message: got %v", err.Error())

	// Set the initial state and check the result
	err = sm.SetState(state1)
	assert.Nil(t, err, "unexpected error: %v", err)
	assert.Equal(t, state1, sm.activeState, "expected activeState to be 'Start', got %v", sm.activeState)

	// Re-setting active state
	assert.Nil(t, err, "unexpected error: %v", err)

	// Test 3: Transition from Start to Middle state
	err = sm.SetState(state2)
	assert.Nil(t, err, "unexpected error: %v", err)
	assert.Equal(t, state2, sm.activeState, "expected activeState to be 'Middle', got %v", sm.activeState)
	assert.Equal(t, state1, sm.previousState, "expected previousState to be 'Start', got %v", sm.previousState)

	// Test 4: Transition to End state if transition is allowed
	state2.availableStates = []State{state3}
	sm.activeState = state2
	err = sm.SetState(state3)
	assert.Nil(t, err, "unexpected error: %v", err)

	assert.Equal(t, state3, sm.activeState, "expected activeState to be 'End'")
	assert.Equal(t, state2, sm.previousState, "expected previousState to be 'Middle'")

	// Test 5: Transition to inaccessible state
	err = sm.SetState(state1)
	assert.NotNil(t, err)
	assert.Equal(t, "Can not move to this state", err.Error())
}

func TestStateMachine_SetStateByName(t *testing.T) {

	state3 := State{name: "State3"}

	state2 := State{
		name:            "State2",
		availableStates: []State{state3},
	}

	state1 := State{name: "State1", availableStates: []State{state2}}

	sm := &StateMachine{
		states:      []State{state1, state2, state3},
		activeState: State{name: ""},
	}

	err := sm.SetStateByName("State1")
	assert.Nil(t, err)
	assert.Equal(t, "State1", sm.activeState.GetName())

	err = sm.SetStateByName("State2")
	assert.Nil(t, err)
	assert.Equal(t, "State2", sm.activeState.GetName())

	err = sm.SetStateByName("NonExistentState")
	assert.NotNil(t, err)
	assert.Equal(t, "This state is not unavailable", err.Error())

	err = sm.SetStateByName("")
	assert.NotNil(t, err)
	assert.Equal(t, "State has empty name", err.Error())
}

func TestStateMachine_GetPreviousState(t *testing.T) {

	prevState := State{name: "PreviousState"}
	activeState := State{name: "ActiveState"}

	sm := &StateMachine{
		previousState: prevState,
		activeState:   activeState,
	}

	assert.Equal(t, prevState, sm.GetPreviousState(), "GetPreviousState() returned incorrect previous state")
}

func TestStateMachine_GetActiveState(t *testing.T) {

	activeState := State{name: "ActiveState"}

	sm := &StateMachine{
		activeState: activeState,
	}

	assert.Equal(t, activeState, sm.GetActiveState(), "GetActiveState() returned incorrect active state")
}
