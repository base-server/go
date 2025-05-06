package main_sub_test

import (
	"flag"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/common/main_sub"
	"github.com/common-library/go/log/slog"
)

func TestRunWithKlog(t *testing.T) {
	{
		os.Args = []string{"test"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		job := func() error { return nil }
		if err := (&main_sub.Main{}).RunWithKlog(main_sub.GRPC, job, job); err.Error() != "invalid flag" {
			t.Fatal(err)
		}
	}

	{
		os.Args = []string{"test", "-config-file=invalid"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		job := func() error { return nil }
		if err := (&main_sub.Main{}).RunWithKlog(main_sub.GRPC, job, job); err.Error() != `Unsupported Config Type ""` {
			t.Fatal(err)
		}
	}

	{
		go func() {
			os.Args = []string{"test", "-config-file=../config/config.yaml"}
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			job := func() error { return nil }
			if err := (&main_sub.Main{}).RunWithKlog(main_sub.GRPC, job, job); err != nil {
				t.Fatal(err)
			}
		}()
		time.Sleep(200 * time.Millisecond)

		if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRunWithSlog(t *testing.T) {
	{
		os.Args = []string{"test"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		job := func(*slog.Log) error { return nil }
		if err := (&main_sub.Main{}).RunWithSlog(main_sub.GRPC, job, job); err.Error() != "invalid flag" {
			t.Fatal(err)
		}
	}

	{
		os.Args = []string{"test", "-config-file=invalid"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		job := func(*slog.Log) error { return nil }
		if err := (&main_sub.Main{}).RunWithSlog(main_sub.GRPC, job, job); err.Error() != `Unsupported Config Type ""` {
			t.Fatal(err)
		}
	}

	{
		go func() {
			os.Args = []string{"test", "-config-file=../config/config.yaml"}
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			job := func(*slog.Log) error { return nil }
			if err := (&main_sub.Main{}).RunWithSlog(main_sub.GRPC, job, job); err != nil {
				t.Fatal(err)
			}
		}()
		time.Sleep(200 * time.Millisecond)

		if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
			t.Fatal(err)
		}
	}
}
