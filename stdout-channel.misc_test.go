package stdio

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/dsnet/golib/memfile"
)

// https://stackoverflow.com/questions/17863821/how-to-read-last-lines-from-a-big-file-with-go-every-10-secs
func getLastLineWithSeek(fileHandle *memfile.File) string {
	line := ""

	b := fileHandle.Bytes()
	fileSize := len(b)
	if fileSize == 0 {
		return line
	}

	text := string(b)
	rawLines := strings.Split(text, "\n")

	if len(rawLines) == 0 {
		return line
	}

	lines := []string{}
	for _, l := range rawLines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}

		lines = append(lines, l)
	}

	if len(lines) == 0 {
		return line
	}

	line = lines[len(lines)-1]
	return line
}

func TestWriter_Write(t *testing.T) {
	// t.Error(">>> TestWriter_Write: ...")
	memOut := memfile.New([]byte(""))
	memErr := memfile.New([]byte(""))
	StdOut = io.Writer(memOut) // NewWriter(memFsOut)
	StdErr = io.Writer(memErr) // NewWriter(memFsErr)
	wrtOut := NewWriter(&StdOut)
	wrtErr := NewWriter(&StdErr)

	defer func() {
		<-WaitLoggerUntilEnd()
	}()
	InitWriter()

	// redundance call InitWriter should not ignore
	InitWriter()

	tests := []struct {
		description string
		rawPayload  []byte
		expectText  string
		writer      io.Writer
		memory      *memfile.File
	}{
		{
			description: "text 1",
			rawPayload:  []byte("text 1\n"),
			expectText:  "text 1",
			writer:      wrtOut,
			memory:      memOut,
		},
		{
			description: "text 2",
			rawPayload:  []byte("text 2\n"),
			expectText:  "text 2",
			writer:      wrtOut,
			memory:      memOut,
		},
		{
			description: "text 3",
			rawPayload:  []byte("text 3\n"),
			expectText:  "text 3",
			writer:      wrtErr,
			memory:      memErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt.writer.Write(tt.rawPayload)
			time.Sleep(10 * time.Millisecond)
			lastLine := getLastLineWithSeek(tt.memory)

			if lastLine != tt.expectText {
				t.Errorf(">>> TestWriter_Write, Expect: [%s], but got [%s]", tt.expectText, lastLine)
			}
		})
	}
}
