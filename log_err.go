package rlog

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/rohanthewiz/serr"
)

// Logging for structured errors (SErr)
// Take an err (preferrably a SErr), and optional message argument followed by attribute - value pairs
// The first extra is always treated as a message. Additional attribute - values should be in pairs
// Therefore typically you will have an even number of arguments to this function
// unless just a single error is supplied
// Example:
//	er := serr.New("Just testing an error", "attribute1", "value1", "attribute2", "value2")
//	logger.LogErr(er, "Testing out LogErr()", 'attribute3", "value3", "attribute4", "value4")
//  logger.LogErr(er)
//  logger.LogErr(er, "Just an err and a message")
func LogErr(err error, extras ...string) {
	if err == nil {
		Log(LogLevel.Info, "In LogErr Not logging a nil err", "called from", serr.FunctionLoc(serr.CallerIndirection.Caller))
		return
	}

	var errs []string // accumulate "error" fields

	// Add error string from original error
	if er := err.Error(); er != "" {
		errs = []string{er}
	}

	lnExtras := len(extras)
	msgs := []string{} // accumulate "msg" fields

	flds := logrus.Fields{}

	// The first extra is always a message
	if lnExtras > 0 {
		msgs = []string{extras[0]}
	}

	var pairs []string
	if lnExtras > 1 {
		pairs = extras[1:]
	}

	// If error is structured error, get key vals
	if ser, ok := err.(serr.SErr); ok {
		// if lnExtras > 1 {
		// 	ser.(pairs...)
		// }

		for key, val := range ser.FieldsMap() {
			if key != "" {
				switch strings.ToLower(key) {
				case "error":
					errs = append(errs, val)
				case "msg":
					msgs = append(msgs, val)
				case strings.ToLower(serr.UserMsgKey):
					continue // that one is for UI only
				case strings.ToLower(serr.UserMsgSeverityKey):
					continue // that one is for UI only
				default:
					flds[key] = val
				}
			}
		}

	} else {
		key := ""
		for i, str := range pairs {
			if i%2 == 0 { // even position is a key
				key = str
			} else {
				flds[key] = str
			}
		}

		// Fixup / Validate
		if lnExtras > 1 && len(pairs)%2 != 0 {
			logrus.Warn("Other than a single error object, an even number of arguments os required to the LogErr function. Odd argument may be dropped")
		}
	}

	// message is required by logrus so use the original error string if msgs empty
	if len(msgs) == 0 {
		msgs = []string{err.Error()}
	}
	// Populate the "error" field
	if len(errs) > 0 {
		flds["error"] = strings.Join(errs, " - ")
	}

	msg := strings.Join(msgs, " - ")
	msg = logPrefix + msg
	logrus.WithFields(flds).Error(msg)
}
