package types

type Msg struct {
	Msg string `json:"msg"`
}

type Composed struct {
	Msg     string            `json:"msg"`
	Opts    []Options         `json:"opts"`
	Details map[string]string `json:"details"`
}
