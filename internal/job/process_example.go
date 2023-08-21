package job

import (
	"time"

	"github.com/google/uuid"
)

type ProcessExample struct {
	JobID uuid.UUID `json:"job_id"`
	Data  string    `json:"data"`
}

func (p *ProcessExample) Handle() error {
	// Process the job here.
	time.Sleep(1 * time.Second)
	return nil
}
