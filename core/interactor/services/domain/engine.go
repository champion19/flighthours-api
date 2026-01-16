package domain

// Engine represents the engine type domain model
// Contains the engine type information for aircraft models
type Engine struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToLogger returns a slice of strings for logging engine information
func (e *Engine) ToLogger() []string {
	return []string{
		"id:" + e.ID,
		"name:" + e.Name,
	}
}
