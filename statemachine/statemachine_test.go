package statemachine

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateMachineError_Error(t *testing.T) {

	expectedMessage := "something went wrong"

	err := StateMachineError{message: expectedMessage}

	if err.Error() != expectedMessage {
		t.Errorf("Error() = %q; want %q", err.Error(), expectedMessage)
	}
}
func TestState_GetStartCommand(t *testing.T) {

	state := State{
		name:              "Init",
		startCommand:      "start",
		availableCommands: []string{"run", "stop"},
		availableStates:   []State{},
	}

	expectedCommand := "start"

	if state.GetStartCommand() != expectedCommand {
		t.Errorf("GetStartCommand(): %q; want %q", state.GetStartCommand(), expectedCommand)
	}
}
func TestState_GetName(t *testing.T) {
	state := State{name: "Jambo"}

	if state.GetName() != "Jambo" {
		t.Errorf("GetName(): %q; want %q", state.GetName(), "Jambo")
	}
}

func TestState_GetAvailableCommands(t *testing.T) {
	state := State{
		availableCommands: []string{"run", "stop"},
	}

	expectedCommands := []string{"run", "stops"}

	if len(expectedCommands) != len(state.availableCommands) {
		t.Errorf("GetAvailableCommands(): len state.AvailableCommands = %d; expected %d",
			len(state.availableCommands), len(expectedCommands))
		return
	}
	for index := range expectedCommands {
		if expectedCommands[index] != state.availableCommands[index] {
			t.Errorf("GetAvailableCommands(): %q; expected %q",
				state.availableCommands[index], expectedCommands[index])
		}
	}
}

func TestState_GetAvailableStates(t *testing.T) {
	state1 := State{name: "Jambo💀"}
	state2 := State{name: "Igor🌈"}
	mainState := State{availableStates: []State{state1, state2}}
	expectedState := State{availableStates: []State{state1, state2}}

	if !reflect.DeepEqual(mainState.availableStates, expectedState.availableStates) {
		t.Errorf("GetAvailableStates(): states not equal,\n %v, expected %v",
			mainState.availableStates, expectedState.availableStates)
	}

}

func TestState_SetAvailableStates(t *testing.T) {
	state1 := State{name: "Stater1"}
	state2 := State{name: "Stater2"}

	mainState := &State{}
	mainState.SetAvailableStates(state1, state2)
	assert.Len(t, mainState.availableStates, 2)
	assert.Equal(t, mainState.availableStates[0].name, "Stater1")
	assert.Equal(t, mainState.availableStates[1].name, "Stater2")

	// assert.Panics(t, func() {
	// 	mainState.SetAvailableStates("InvalidType") //Не получилось передать херь, чтоб вызвать панику
	// }, "Statemachine|SetAvailableStates|\nerror: Object is not type State")
}

func TestState_NewState(t *testing.T) {
	receivedName := "Name"
	receivedStartCommand := "Run"
	receivedAvailableCommands := []string{"Run", "Stop"}

	expectedState := &State{
		name:              "Name",
		startCommand:      "Run",
		availableCommands: []string{"Run", "Stop"},
	}

	testState := NewState(receivedName, receivedStartCommand, receivedAvailableCommands...)

	if !(assert.Equal(t, expectedState, testState)) {
		t.Errorf("NewState(): states not equal,\n %v, expected %v",
			testState, expectedState)
	}
}

func TestStateMachine_CompareStates(t *testing.T) {

}

func TestStateMachine_AddStates(t *testing.T) {
	state1, state2, state3 := State{}, State{}, State{}
	expectedStateMachine := &StateMachine{
		states: []State{state1, state2, state3},
	}
	testStateMachine := &StateMachine{}
	testStateMachine.AddStates(state1, state2, state3)

	if !reflect.DeepEqual(testStateMachine.states, expectedStateMachine.states) {
		t.Errorf("AddState(): states not equal,\n %v, expected %v",
			testStateMachine.states, expectedStateMachine.states)
	}
}

func TestStateMachine_SetState(t *testing.T) {

	state1 := State{name: "Start"}
	state2 := State{name: "Middle"}
	state3 := State{name: "End"}

	sm := &StateMachine{
		states: []State{state1, state2, state3},
	}

	// Test 1: Setting status with empty name (expect error)
	err := sm.SetState(State{name: ""})
	if err == nil || err.Error() != "State has empty name" {
		t.Errorf("expected error 'State has empty name', got %v", err)
	}

	// Test 2: Transition to state that is not in list (expect error)
	err = sm.SetState(State{name: "NonExistent"})
	if err == nil || err.Error() != "This state is not unavailable" {
		t.Errorf("expected error 'This state is not unavailable', got %v", err)
	}

	// Set the initial state and check the result
	err = sm.SetState(state1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(sm.activeState, state1) {
		t.Errorf("expected activeState to be 'Start', got %v", sm.activeState)
	}

	// Re-setting active state (non expect error)
	err = sm.SetState(state1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test 3: Transition from Start to Middle state (non expect error)
	err = sm.SetState(state2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(sm.activeState, state2) {
		t.Errorf("expected activeState to be 'Middle', got %v", sm.activeState)
	}
	if !reflect.DeepEqual(sm.previousState, state1) {
		t.Errorf("expected previousState to be 'Start', got %v", sm.previousState)
	}

	// Test 4: Transition to End state if transition is allowed
	state2.availableStates = []State{state3}
	sm.activeState = state2
	err = sm.SetState(state3)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(sm.activeState, state3) {
		t.Errorf("expected activeState to be 'End', got %v", sm.activeState)
	}
	if !reflect.DeepEqual(sm.previousState, state2) {
		t.Errorf("expected previousState to be 'Middle', got %v", sm.previousState)
	}

	// Test 5: Transition to inaccessible state (expect error)
	err = sm.SetState(state1)
	if err == nil || err.Error() != "Can not move to this state" {
		t.Errorf("expected error 'Can not move to this state', got %v", err)
	}
}

func TestStateMachine_SetStateByName(t *testing.T) {
	state1 := State{name: "State1"}
	state2 := State{
		name:            "State2",
		availableStates: []State{state1},
	}
	state3 := State{
		name:            "State3",
		availableStates: []State{state1, state2},
	}

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
