package types

import "time"

type Date struct {
	id         int
	start_date time.Time
	end_date   time.Time
}
