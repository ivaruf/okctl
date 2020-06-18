package run

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type Runner interface {
	Run(progress io.Writer, args []string) ([]byte, error)
}

type Run struct {
	WorkingDirectory string
	BinaryPath       string
	Env              []string
	Logger           *logrus.Logger
}

func AnonymizeEnv(entries []string) []string {
	out := make([]string, len(entries))

	hide := []string{
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
	}

	for _, e := range entries {
		for _, h := range hide {
			if strings.Contains(e, h) {
				e = fmt.Sprintf("%s=XXXXXXX", h)
				break
			}
		}

		out = append(out, e)
	}

	return out
}

func (r *Run) Run(progress io.Writer, args []string) ([]byte, error) {
	var errOut, errErr error

	ctxLogger := logrus.WithFields(
		logrus.Fields{
			"component": "generic_runner",
			"binary":    r.BinaryPath,
			"args":      args,
			"env":       AnonymizeEnv(r.Env),
			"dir":       r.WorkingDirectory,
		},
	)

	cmd := &exec.Cmd{
		Path: r.BinaryPath,
		Args: append([]string{r.BinaryPath}, args...),
		Env:  r.Env,
		Dir:  r.WorkingDirectory,
	}

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	var errBuff, outBuff bytes.Buffer
	stdout := io.MultiWriter(progress, &outBuff)
	stderr := io.MultiWriter(progress, &errBuff)

	ctxLogger.Info("Starting execution of provided command")

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		_, errOut = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errErr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		return errBuff.Bytes(), err
	}

	if errOut != nil || errErr != nil {
		return errBuff.Bytes(), err
	}

	return outBuff.Bytes(), nil
}

func New(logger *logrus.Logger, workingDirectory, binaryPath string, env []string) *Run {
	return &Run{
		WorkingDirectory: workingDirectory,
		BinaryPath:       binaryPath,
		Env:              env,
		Logger:           logger,
	}
}