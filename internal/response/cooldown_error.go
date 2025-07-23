package response

import "time"

type CooldownError struct {
	Message    string
	RetryAfter time.Duration
}

func (e *CooldownError) Error() string {
	return e.Message
}

type LoginAttemptError struct {
	Message           string
	RemainingAttempts int
}

func (e *LoginAttemptError) Error() string {
	return e.Message
}
