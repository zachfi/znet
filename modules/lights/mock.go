package lights

import "context"

type MockLight struct {
	AlertCalls         map[string]int
	OffCalls           map[string]int
	OnCalls            map[string]int
	RandomColorCalls   map[string]int
	SetBrightnessCalls map[string]int
	SetColorCalls      map[string]int
	SetColorTempCalls  map[string]int
	ToggleCalls        map[string]int
}

func (m *MockLight) Alert(ctx context.Context, groupName string) error {
	if len(m.AlertCalls) == 0 {
		m.AlertCalls = make(map[string]int)
	}
	m.AlertCalls[groupName]++
	return nil
}

func (m *MockLight) SetBrightness(ctx context.Context, groupName string, brightness int32) error {
	if len(m.SetBrightnessCalls) == 0 {
		m.SetBrightnessCalls = make(map[string]int)
	}
	m.SetBrightnessCalls[groupName]++
	return nil
}

func (m *MockLight) Off(ctx context.Context, groupName string) error {
	if len(m.OffCalls) == 0 {
		m.OffCalls = make(map[string]int)
	}
	m.OffCalls[groupName]++
	return nil
}

func (m *MockLight) On(ctx context.Context, groupName string) error {
	if len(m.OnCalls) == 0 {
		m.OnCalls = make(map[string]int)
	}
	m.OnCalls[groupName]++
	return nil
}

func (m *MockLight) RandomColor(ctx context.Context, groupName string, colors []string) error {
	if len(m.RandomColorCalls) == 0 {
		m.RandomColorCalls = make(map[string]int)
	}
	m.RandomColorCalls[groupName]++
	return nil
}

func (m *MockLight) SetColor(ctx context.Context, groupName string, hex string) error {
	if len(m.SetColorCalls) == 0 {
		m.SetColorCalls = make(map[string]int)
	}
	m.SetColorCalls[groupName]++
	return nil
}

func (m *MockLight) Toggle(ctx context.Context, groupName string) error {
	if len(m.ToggleCalls) == 0 {
		m.ToggleCalls = make(map[string]int)
	}
	m.ToggleCalls[groupName]++
	return nil
}

func (m *MockLight) SetColorTemp(ctx context.Context, groupName string, temp int32) error {
	if len(m.SetColorTempCalls) == 0 {
		m.SetColorTempCalls = make(map[string]int)
	}
	m.SetColorTempCalls[groupName]++
	return nil
}
