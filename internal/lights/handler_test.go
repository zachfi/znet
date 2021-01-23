package lights

import "testing"

type mockLight struct {
	alertCalls       int
	dimCalls         int
	offCalls         int
	onCalls          int
	randomColorCalls int
	setColorCalls    int
	toggleCalls      int
}

func (m *mockLight) Alert(groupName string) error {
	m.alertCalls++
	return nil
}

func (m *mockLight) Dim(groupName string, brightness int32) error {
	m.dimCalls++
	return nil
}

func (m *mockLight) Off(groupName string) error {
	m.offCalls++
	return nil
}

func (m *mockLight) On(groupName string) error {
	m.onCalls++
	return nil
}

func (m *mockLight) RandomColor(groupName string, colors []string) error {
	m.randomColorCalls++
	return nil
}

func (m *mockLight) SetColor(groupName string, hex string) error {
	m.setColorCalls++
	return nil
}

func (m *mockLight) Toggle(groupName string) error {
	m.toggleCalls++
	return nil
}

func TestMockHandler_interface(t *testing.T) {
	var l Handler = &mockLight{}
	t.Log(l)
}
