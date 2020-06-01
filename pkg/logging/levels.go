package logging

// LogLevel just wraps a string to be able to add LogContext specific to log levels
type LogLevel string

// Info is sets level=info in the log output
const Info = LogLevel("info")

// Error sets level=error in the log output
const Error = LogLevel("error")

// Fatal sets level=fatal in the log output
const Fatal = LogLevel("fatal")

// LogContext returns the map that states that key value of `level={{l}}`
func (l LogLevel) LogContext() map[string]string {
	return map[string]string{
		"level": string(l),
	}
}