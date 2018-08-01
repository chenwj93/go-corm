package logs

import (
	"log"
	"fmt"
)

func init(){
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Debug(format ...interface{}) {
	if DebugLevel <= DEBUG {
		log.Output(CallDepth, fmt.Sprint("[corm-DEBUG]", format))
	}
}

func Info(format ...interface{}) {
	if DebugLevel <= INFO {
		log.Output(CallDepth, fmt.Sprint("[corm-INFO]", format))
	}
}

func Warn(format ...interface{}) {
	if DebugLevel <= WARN {
		log.Output(CallDepth, fmt.Sprint("[corm-WARN]", format))
	}
}

func Error(format ...interface{}) {
	if DebugLevel <= ERROR {
		log.Output(CallDepth, fmt.Sprint("[corm-ERROR]", format))
	}
}

func Fatal(format ...interface{}) {
	if DebugLevel <= FATAL {
		log.Output(CallDepth, fmt.Sprint("[corm-FATAL]", format))
	}
}
