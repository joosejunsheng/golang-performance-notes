package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestBufferedWriterShortContent(t *testing.T) {
	bufferedFilePath := "buffered_write.txt"
	WriteUsingBuffer(bufferedFilePath, false)
	defer os.Remove(bufferedFilePath)

	content, _ := os.ReadFile(bufferedFilePath)
	if len(content) == 0 {
		t.Fatalf("Expected something, got '%s'", string(content))
	}
}

func TestBufferedWriterLongContent(t *testing.T) {
	bufferedFilePath := "buffered_write.txt"
	WriteUsingBuffer(bufferedFilePath, true)
	defer os.Remove(bufferedFilePath)

	content, _ := os.ReadFile(bufferedFilePath)
	if len(content) == 0 {
		t.Fatalf("Expected something, got '%s'", string(content))
	}
}

func TestBufferedWriterFasterForLongContent(t *testing.T) {
	normalFilePath := "normal_write.txt"
	bufferedFilePath := "buffered_write.txt"

	start := time.Now()
	WriteNormally(normalFilePath, true)
	normalElapsed := time.Since(start)
	fmt.Printf("WriteNormally took %s\n", normalElapsed)

	os.Remove(normalFilePath)

	newStart := time.Now()
	WriteUsingBuffer(bufferedFilePath, true)
	bufferedElapsed := time.Since(newStart)
	fmt.Printf("WriteUsingBuffer took %s\n", bufferedElapsed)

	os.Remove(bufferedFilePath)

	t.Logf("Normal write took: %v", normalElapsed)
	t.Logf("Buffered write took: %v", bufferedElapsed)

	if bufferedElapsed >= normalElapsed {
		t.Fatalf("BufferedWriter was not faster. Normal write: %v, Buffered write: %v", normalElapsed, bufferedElapsed)
	}
}
