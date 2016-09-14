//package syslogger provides a NewSysLogger func that given a prefix string
// and a syslog prioryti construct without failing a SysLogger defined as a struct
// embeding a chan of messages a done channel and a syslogger.
// the other provide functions are Log() wich sends the message, and Close()
// wich close the logger is it important to close the logger otherwise we would
// leak from the gourutine who does the job in the background.

package syslogger

import (
	"log"
	"log/syslog"
	"time"
)

const (
	msgMaxLen = 79
)

type Logger struct {
	logger   *log.Logger
	Messages chan string
	Done     chan struct{}
}

func New(prefix string, p syslog.Priority) *Logger {
	stdlog, err := syslog.NewLogger(p|syslog.LOG_USER, log.Lshortfile)
	if err != nil {
		//sonic shouldn't start if the syslog doesn't work so panic
		panic(err)
	}

	stdlog.SetPrefix(prefix)

	log := &Logger{
		logger:   stdlog,
		Messages: make(chan string, 32),
		Done:     make(chan struct{}, 1),
	}

	go func(log *Logger) {
		for {
			select {
			case <-log.Done:
				close(log.Messages)
				for msg := range log.Messages {
					log.log(msg)
				}
				time.Sleep(time.Millisecond * 250)

				return

			case msg := <-log.Messages:
				log.log(msg)
			}
		}
	}(log)

	return log

}

// clients code should use uppercase Log to send messages to the listeining goroutine
func (s *Logger) toLog(msg string) {
	truncated := []byte(msg)

	if len(msg) > msgMaxLen {
		// msg should must not be over 79 characthers so panic
		truncated = truncated[:79]
	}
	s.Messages <- string(truncated)
}

// Log does what it says but it doesn't assure the message will be send.
func (s *Logger) Log(msg string) {
	s.toLog(msg)
}

// LogAndWait as the name suggest logs a message and waits almost assuring
// the the message will be sent, use case are when sending a log before
// the planed programs finish. Adds a faily expencive overhead, 5 millisecons
func (s *Logger) MustLog(msg string) {
	s.toLog(msg)
	time.Sleep(time.Millisecond * 5)
}

// Typically called in defer log.Close() fashion for short lived objects otherwise
// must be called when the logger isn't needed any longer, once called if any
// messages are still left on the Messages channel, they will be sent to syslog before
// actually exting, any message sent after the close method is called will a result
// in a sent operation to a close channel.

func (s *Logger) Close() {
	s.Done <- struct{}{}
}

// log not uppercase is used internaly to actually print to syslog
func (s *Logger) log(msg string) {
	s.logger.Output(2, msg)
}
