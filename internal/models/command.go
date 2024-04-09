package models

// Command
type Command struct {
	ID      int    `json:"id" gorm:"unique"`
	Command string `json:"command"`
	Output  string `json:"output,omitempty"`
	Status  string `json:"status,omitempty"`
}
