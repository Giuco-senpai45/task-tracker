package requests

import "time"

type TaskCreatePayload struct {
	Name     string    `json:"name"`
	Deadline time.Time `json:"deadline"`
}
