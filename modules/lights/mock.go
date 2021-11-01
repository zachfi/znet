package lights

type MockLight struct {
	AlertCalls       map[string]int
	DimCalls         map[string]int
	OffCalls         map[string]int
	OnCalls          map[string]int
	RandomColorCalls map[string]int
	SetColorCalls    map[string]int
	ToggleCalls      map[string]int
}

func (m *MockLight) Alert(groupName string) error {
	if len(m.AlertCalls) == 0 {
		m.AlertCalls = make(map[string]int)
	}
	m.AlertCalls[groupName]++
	return nil
}

func (m *MockLight) Dim(groupName string, brightness int32) error {
	if len(m.DimCalls) == 0 {
		m.DimCalls = make(map[string]int)
	}
	m.DimCalls[groupName]++
	return nil
}

func (m *MockLight) Off(groupName string) error {
	if len(m.OffCalls) == 0 {
		m.OffCalls = make(map[string]int)
	}
	m.OffCalls[groupName]++
	return nil
}

func (m *MockLight) On(groupName string) error {
	if len(m.OnCalls) == 0 {
		m.OnCalls = make(map[string]int)
	}
	m.OnCalls[groupName]++
	return nil
}

func (m *MockLight) RandomColor(groupName string, colors []string) error {
	if len(m.RandomColorCalls) == 0 {
		m.RandomColorCalls = make(map[string]int)
	}
	m.RandomColorCalls[groupName]++
	return nil
}

func (m *MockLight) SetColor(groupName string, hex string) error {
	if len(m.SetColorCalls) == 0 {
		m.SetColorCalls = make(map[string]int)
	}
	m.SetColorCalls[groupName]++
	return nil
}

func (m *MockLight) Toggle(groupName string) error {
	if len(m.ToggleCalls) == 0 {
		m.ToggleCalls = make(map[string]int)
	}
	m.ToggleCalls[groupName]++
	return nil
}
