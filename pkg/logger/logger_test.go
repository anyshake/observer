package logger

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestLoggerLevelsBufferAndFile(t *testing.T) {
	Init()
	if err := SetLevel(INFO); err != nil {
		t.Fatal(err)
	}
	if err := SetLevel(WARN); err != nil {
		t.Fatal(err)
	}
	if err := SetLevel(ERROR); err != nil {
		t.Fatal(err)
	}
	if err := SetLevel(FATAL); err != nil {
		t.Fatal(err)
	}
	if err := SetLevel(LogLevel(99)); err == nil {
		t.Fatal("unknown level accepted")
	}
	if err := SetLevel(INFO); err != nil {
		t.Fatal(err)
	}

	named := GetLogger("Observer")
	named.Infof("info %s", "message")
	named.Infoln("info", "line")
	named.Warnf("warn %d", 1)
	named.Warnln("warn")
	named.Errorf("error %t", true)
	named.Errorln("error")

	if got := GetLogger(TestLoggerLevelsBufferAndFile); got == nil {
		t.Fatal("function logger is nil")
	}
	if got := GetLogger(42); got == nil {
		t.Fatal("unknown logger is nil")
	}

	buffer := RegisterBufferLogger(2)
	GetLogger("buffer").Infoln("buffered")
	if buffer.Len() == 0 {
		t.Fatal("buffer logger stored nothing")
	}

	path := filepath.Join(t.TempDir(), "observer.log")
	RegisterFileLogger(path, 1, 1, 1)
	GetLogger("file").Infoln("filed")

	stdlogLine := []byte("stdlib line\n")
	if n, err := (standardLogWriter{}).Write(stdlogLine); err != nil || n != len(stdlogLine) {
		t.Fatalf("standardLogWriter.Write = %d, %v", n, err)
	}
	if n, err := (standardLogWriter{}).Write([]byte("   \n")); err != nil || n != 4 {
		t.Fatalf("blank standard log write = %d, %v", n, err)
	}
}

func TestFatalExits(t *testing.T) {
	Init()
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()
	previous := zerolog.FatalExitFunc
	t.Cleanup(func() { zerolog.FatalExitFunc = previous })

	for _, call := range []func(){
		func() { GetLogger("fatal").Fatalf("stop %d", 1) },
		func() { GetLogger("fatal").Fatalln("stop") },
	} {
		exited := false
		zerolog.FatalExitFunc = func() {
			exited = true
			panic("logger exit")
		}
		func() {
			defer func() { _ = recover() }()
			call()
		}()
		if !exited {
			t.Fatal("fatal logger returned")
		}
	}
}
