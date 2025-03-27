package logger

// Logger interface describes project logger methods
type Logger interface {
	// Info logs the provided arguments at InfoLevel. Spaces are added between arguments when neither is a string.
	Info(args ...interface{})

	// Infof formats the message according to the format specifier and logs it at InfoLevel.
	Infof(template string, args ...interface{})

	// Warn logs the provided arguments at WarnLevel. Spaces are added between arguments when neither is a string.
	Warn(args ...interface{})

	// Warnf formats the message according to the format specifier and logs it at WarnLevel.
	Warnf(template string, args ...interface{})

	// Error logs the provided arguments at ErrorLevel. Spaces are added between arguments when neither is a string.
	Error(args ...interface{})

	// Errorf formats the message according to the format specifier and logs it at ErrorLevel.
	Errorf(template string, args ...interface{})

	// Debug logs the provided arguments at DebugLevel. Spaces are added between arguments when neither is a string.
	Debug(args ...interface{})

	// Debugf formats the message according to the format specifier and logs it at DebugLevel.
	Debugf(template string, args ...interface{})

	// Fatal constructs a message with the provided arguments and calls os.Exit. Spaces are added between arguments when neither is a string.
	Fatal(args ...interface{})

	// Fatalf formats the message according to the format specifier and calls os.Exit.
	Fatalf(template string, args ...interface{})

	// WithFields adds a variadic number of fields to the logging context. It accepts a
	// mix of strongly-typed Field objects and loosely-typed key-value pairs. When
	// processing pairs, the first element of the pair is used as the field key
	// and the second as the field value.
	WithFields(args ...interface{}) Logger
}
