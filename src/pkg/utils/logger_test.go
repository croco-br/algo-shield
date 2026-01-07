package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Logger_Info_WhenCalled_ThenLogsMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST")
	logger.logger.SetOutput(&buf)

	logger.Info("test message: %s", "value")

	output := buf.String()
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "test message: value")
}

func Test_Logger_Error_WhenCalled_ThenLogsMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST")
	logger.logger.SetOutput(&buf)

	logger.Error("error occurred: %s", "details")

	output := buf.String()
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "error occurred: details")
}

func Test_Logger_Debug_WhenCalled_ThenLogsMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST")
	logger.logger.SetOutput(&buf)

	logger.Debug("debug info: %d", 42)

	output := buf.String()
	assert.Contains(t, output, "[DEBUG]")
	assert.Contains(t, output, "debug info: 42")
}

func Test_Logger_Warn_WhenCalled_ThenLogsMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST")
	logger.logger.SetOutput(&buf)

	logger.Warn("warning: %s", "something")

	output := buf.String()
	assert.Contains(t, output, "[WARN]")
	assert.Contains(t, output, "warning: something")
}

func Test_NewLogger_WhenCreated_ThenHasPrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("MYAPP")
	logger.logger.SetOutput(&buf)

	logger.Info("test")

	output := buf.String()
	assert.Contains(t, output, "MYAPP")
}

func Test_Logger_Info_WhenNoArgs_ThenLogsSimpleMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST")
	logger.logger.SetOutput(&buf)

	logger.Info("simple message")

	output := buf.String()
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "simple message")
}

func Test_Logger_Error_WhenMultipleArgs_ThenFormatsCorrectly(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST")
	logger.logger.SetOutput(&buf)

	logger.Error("error: %s code: %d", "failure", 500)

	output := buf.String()
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "error: failure code: 500")
}

func Test_Logger_WhenEmptyPrefix_ThenStillWorks(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("")
	logger.logger.SetOutput(&buf)

	logger.Info("test")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.NotEmpty(t, lines)
}
