package pb

type SocketMessage struct {
	User    string `json:"user"`
	Group   string `json:"group"`
	Time    string `json:"time"`
	Content string `json:"content"`
	Token   string `json:"token"`
	Action  string `json:"action"`
}

type ChatMessage struct {
	User      string    `json:"user"`
	Time      string    `json:"time"`
	Content   string    `json:"content"`
}
