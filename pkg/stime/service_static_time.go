package stime

import (
	"time"
)

type StaticTimeService interface {
	Change(update time.Time) time.Time
	TimeService
}

type staticTimeService struct {
	time time.Time
}

func NewStaticTimeService() StaticTimeService {
	return &staticTimeService{
		time: time.Now().In(time.UTC),
	}
}

// DeleteInvite - Delete an invite that was sent and invalidate the token.
func (s *staticTimeService) Now() time.Time {
	return s.time
}

func (s *staticTimeService) Change(update time.Time) time.Time {
	s.time = update
	return s.time
}
