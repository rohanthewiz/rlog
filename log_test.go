package rlog

import (
	"errors"
	"testing"

	"github.com/rohanthewiz/serr"
)

func TestLog(t *testing.T) {
	formatter := "json"
	InitLog(LogConfig{
		Formatter: formatter,
		LogLevel:  "debug",
		SlackrusCfg: SlackrusCfg{
			Enabled:  false,
			// Endpoint: sc.Endpoint,
			// LogLevel: sc.LogLevel,
		},
	})
	defer CloseLog()

	Log("info", "Conveying some info", "attribute1", "value1", "attribute2", "value2")
	Log("error", "Some error occurred", "attribute1", "value1", "attribute2", "value2",
		"location", serr.FunctionLoc(serr.CallerIndirection.Caller))

	err := errors.New("This is the original error")

	// We can log a standard error, the message will be err.Error()
	LogErr(err, "message")
	// => ERRO[0000] message	error="This is the original error"

	// Multiple arguments after message are treated as a key, value list and will wrap the error
	// Be careful to use pairs of fields after message.
	LogErr(err, "message", "key1", "value1", "key2", "value2")
	// => ERRO[0000] message error="This is the original error" key1=value1 key2=value2


}