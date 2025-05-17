package logger

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestFileLogger(t *testing.T) {
	t.Run("should be able to append message in the log file", func(t *testing.T) {
		file := createTempFile(t)
		defer file.Close()

		sut := FileLogger{
			file: file,
		}

		want := "test-message"
		sut.AppendMessage(mockLogMessage(want))

		got := removeTimestampFromRow(returnLastRow(t, file.Name()))

		toEqual(t, got, want)
	})

	t.Run("should append the log message in the end of the file", func(t *testing.T) {
		file := createTempFile(t)
		defer file.Close()

		sut := FileLogger{
			file: file,
		}

		messages := []mockLogMessage{
			mockLogMessage("first message"),
			mockLogMessage("second message"),
			mockLogMessage("third message"),
		}
		lastMessage := messages[len(messages)-1].String()

		for _, message := range messages {
			sut.AppendMessage(message)
		}

		got := removeTimestampFromRow(returnLastRow(t, file.Name()))

		toEqual(t, got, lastMessage)
	})

	t.Run("should log the message with a timestamp", func(t *testing.T) {
		file := createTempFile(t)
		defer file.Close()

		timestampRegex := regexp.MustCompile(`^\[[A-Z][a-z]{2}, \d{2} [A-Z][a-z]{2} \d{4} \d{2}:\d{2}:\d{2} [-+]\d{2}\]$`)

		sut := FileLogger{
			file: file,
		}

		message := mockLogMessage("teste message")

		sut.AppendMessage(message)

		got := strings.Split(returnLastRow(t, file.Name()), LOG_SEPARATOR_SYMBOL)
		gotedTimestamp, gotedMessage := got[0], got[1]

		if !timestampRegex.MatchString(gotedTimestamp) {
			t.Errorf("the log message must be saved with the current timestamp")
		}
		toEqual(t, gotedMessage, message.String())
	})
}

type mockLogMessage string

func (s mockLogMessage) Content() string {
	return string(s)
}

func (s mockLogMessage) String() string {
	return string(s)
}

func toEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func createTempFile(t *testing.T) *os.File {
	t.Helper()
	file, err := os.CreateTemp(os.TempDir(), "test.txt")
	if err != nil {
		t.Fatal("Internal test error")
	}

	return file
}

func returnLastRow(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(dir)
	if err != nil {
		t.Fatal("Internal test error")
	}

	rows := strings.Split(string(data), "\n")
	return strings.TrimSuffix(rows[len(rows)-2], "\n")
}

func removeTimestampFromRow(message string) string {
	return strings.Split(message, LOG_SEPARATOR_SYMBOL)[1]
}
