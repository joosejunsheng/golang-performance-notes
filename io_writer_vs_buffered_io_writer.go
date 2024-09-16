package main

import (
	"fmt"
	"os"
	"time"
)

const (
	SOME_TEXT = "SOME TEXTTT"
)

type BufferedWriter struct {
	file           *os.File
	buffer         []byte
	bufferEndIndex int
}

// Option function to modify BufferedWriter
type Option func(*BufferedWriter)

func NewBufferedWriter(options ...Option) *BufferedWriter {

	bw := &BufferedWriter{
		buffer:         make([]byte, 4096),
		bufferEndIndex: 0,
	}

	for _, opt := range options {
		opt(bw)
	}

	return bw
}

// WithFile, option to set file for BufferedWriter
func WithFile(file *os.File) Option {
	return func(bw *BufferedWriter) {
		bw.file = file
	}
}

// WithBufferSize, option to set buffer size for BufferedWriter
func WithBufferSize(size int) Option {
	return func(bw *BufferedWriter) {
		bw.buffer = make([]byte, size)
	}
}

func (w *BufferedWriter) Write(content []byte) {
	if len(content) >= len(w.buffer) {
		// Write directly into file since buffer size is not enough for content
		// Write current buffered content into file
		w.Flush()
		// Write the entire content into file first
		w.file.Write(content)
	} else {
		if w.bufferEndIndex+len(content) > len(w.buffer) {
			// Write current buffered content into file since buffer size is not enough for current buffered + incoming content
			w.Flush()
		}

		// Copy incoming content into buffer
		copy(w.buffer[w.bufferEndIndex:], content)
		// Update buffer end index
		w.bufferEndIndex += len(content)
	}
}

func (w *BufferedWriter) Flush() {
	// Write into file and update end index to 0
	w.file.Write(w.buffer[0:w.bufferEndIndex])
	w.bufferEndIndex = 0
}

func (w *BufferedWriter) WriteString(content string) {
	w.Write([]byte(content))
}

func WriteNormally(filePath string) {
	outputFile, err := os.OpenFile(filePath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	for i := 0; i < 999999; i++ {
		outputFile.WriteString(SOME_TEXT)
	}
}

func WriteUsingBuffer(filePath string) {
	outputFile, err := os.OpenFile(filePath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	writer := NewBufferedWriter(
		WithFile(outputFile),
		WithBufferSize(8192),
	)
	defer writer.Flush()

	for i := 0; i < 999999; i++ {
		writer.WriteString(SOME_TEXT)
	}
}

func main() {
	normalFilePath := "normal_write.txt"
	bufferedFilePath := "buffered_write.txt"

	start := time.Now()
	WriteNormally(normalFilePath)
	elapsed := time.Since(start)
	fmt.Printf("WriteNormally took %s\n", elapsed)

	os.Remove(normalFilePath)

	start = time.Now()
	WriteUsingBuffer(bufferedFilePath)
	elapsed = time.Since(start)
	fmt.Printf("WriteUsingBuffer took %s\n", elapsed)

	os.Remove(bufferedFilePath)

	fmt.Println("Using Buffer Reduced Unnecessary I/O Time!!!")
}
