package main

// Outbound messages (companion → HA)

type stripState struct {
	Index   int     `json:"index"`
	Label   string  `json:"label"`
	Mute    bool    `json:"mute"`
	Gain    float64 `json:"gain"`
	Virtual bool    `json:"virtual"`
	A1      bool    `json:"a1"`
	A2      bool    `json:"a2"`
	A3      bool    `json:"a3"`
	A4      bool    `json:"a4"`
	A5      bool    `json:"a5"`
	B1      bool    `json:"b1"`
	B2      bool    `json:"b2"`
	B3      bool    `json:"b3"`
}

type busState struct {
	Index int     `json:"index"`
	Label string  `json:"label"`
	Mute  bool    `json:"mute"`
	Gain  float64 `json:"gain"`
}

type stateMsg struct {
	Type     string       `json:"type"`
	Kind     string       `json:"kind"`
	Strips   []stripState `json:"strips"`
	Buses    []busState   `json:"buses"`
	Protocol string       `json:"protocol"`
}

type updateMsg struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Index  int    `json:"index"`
	Param  string `json:"param"`
	Value  any    `json:"value"`
}

// Inbound messages (HA → companion)

type setMsg struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Index  int    `json:"index"`
	Param  string `json:"param"`
	Value  any    `json:"value"`
}
