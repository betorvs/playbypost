package types

type Msg struct {
	Msg string `json:"msg"`
}

type Composed struct {
	Msg     string            `json:"msg"`
	Opt     []GenericIDName   `json:"options"`
	Opts    []Options         `json:"opts"`
	Details map[string]string `json:"details"`
}
