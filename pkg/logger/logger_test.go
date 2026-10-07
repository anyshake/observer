package logger

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

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
	mode := os.Getenv("OBSERVER_LOGGER_FATAL")
	if mode == "f" || mode == "ln" {
		Init()
		if mode == "f" {
			GetLogger("fatal").Fatalf("stop %d", 1)
		} else {
			GetLogger("fatal").Fatalln("stop")
		}
		t.Fatal("logger fatal returned")
	}

	for _, mode := range []string{"f", "ln"} {
		t.Run(mode, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestFatalExits$", "-test.count=1")
			cmd.Env = append(os.Environ(), "OBSERVER_LOGGER_FATAL="+mode)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			if err := cmd.Run(); err == nil {
				t.Fatalf("fatal mode %s did not exit: %s", mode, stderr.String())
			}
		})
	}
	_ = log.Logger
}
