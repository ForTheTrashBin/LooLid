package osspecific

import "time"

type TimeSpec struct {
	TimeModify time.Time
	TimeAccess time.Time
	TimeCreate time.Time
}
