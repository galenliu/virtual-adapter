package gateway_addon_golang

type ActionNotification func(action *Action)

type Action struct {
	Name        string `json:"name"`
	AtType      string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ActionFunc      ActionNotification
}

func NewAction() *Action {
	action := &Action{}
	return action
}


