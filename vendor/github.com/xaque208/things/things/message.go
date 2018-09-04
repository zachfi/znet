package things

type Message struct {
	Device   string          `json:"id,omitempty"`
	Priority int             `json:"priority,omitempty"`
	Sensors  []SensorReading `json:"sensors,omitempty"`
	Events   []Event         `json:"events,omitempty"`
	Commands []Command       `json:"commands,omitempty"`
}

type SensorReading struct {
	Name  string  `json:"name,omitempty"`
	Value float32 `json:"value,omitempty"`
	Unit  string  `json:"unit,omitempty"`
}

type Event struct {
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}

type Command struct {
	Name      string
	Arguments CommandArguments
}

type CommandArguments map[string]interface{}
