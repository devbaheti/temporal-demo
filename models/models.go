package models

type Activities struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

type Workflow struct {
	ID         string                 `json:"id"`
	Parameters map[string]interface{} `json:"parameters"`
	Activities []Activities           `json:"activities,required"`
}
