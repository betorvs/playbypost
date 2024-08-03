package types

type Msg struct {
	Msg string `json:"msg"`
}

type Composed struct {
	Msg     string            `json:"msg"`
	Opt     []GenericIDName   `json:"options"`
	Details map[string]string `json:"details"`
}
