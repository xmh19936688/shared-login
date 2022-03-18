package model

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Identity struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
