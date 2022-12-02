package eml

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"sort"
	"time"
)

var maxInt64 = big.NewInt(math.MaxInt64)

// GenMessageID creates and returns a Message-ID, without surrounding angle brackets.
func GenMessageID() (string, error) {
	return generateID("")
}

// GenContentID creates and returns a Content-ID, without surrounding angle brackets.
func GenContentID(filename string) (string, error) {
	return generateID(filename)
}

// generateID creates a globally unique identifier in the Message-ID format (subset of email address),
// optionally having an additional string appended to the local part.
// Example: 11223344556677889900.11.1234567890@localhost
func generateID(appendWith string) (string, error) {
	random, err := rand.Int(rand.Reader, maxInt64)
	if err != nil {
		return "", nil
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	pid := os.Getpid()
	nanoTime := time.Now().UTC().UnixNano()
	if len(appendWith) == 0 {
		return fmt.Sprintf("%d.%d.%d@%s", nanoTime, pid, random, hostname), nil
	}
	return fmt.Sprintf("%d.%d.%d.%s@%s", nanoTime, pid, random, appendWith, hostname), nil
}

// max ...
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// sortedHeaderFields ...
func sortedHeaderFields(stringMap map[string][]string) []string {
	keyCount := 0
	sortedKeys := make([]string, len(stringMap))
	for k := range stringMap {
		sortedKeys[keyCount] = k
		keyCount++
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// bufioReader ...
func bufioReader(r io.Reader) *bufio.Reader {
	if bufferedReader, ok := r.(*bufio.Reader); ok {
		return bufferedReader
	}
	return bufio.NewReader(r)
}

// headerWriter ...
type headerWriter struct {
	w          io.Writer
	curLineLen int
	maxLineLen int
}

// Write ...
func (w *headerWriter) Write(p []byte) (int, error) {
	// TODO: logic for wrapping headers is actually pretty complex for some header types, like received headers
	var total int
	for len(p)+w.curLineLen > w.maxLineLen {
		toWrite := w.maxLineLen - w.curLineLen
		// Wrap at last space, if any
		lastSpace := bytes.LastIndexByte(p[:toWrite], byte(' '))
		if lastSpace > 0 {
			toWrite = lastSpace
		}
		written, err := w.w.Write(p[:toWrite])
		total += written
		if err != nil {
			return total, err
		}
		written, err = w.w.Write([]byte("\r\n "))
		total += written
		if err != nil {
			return total, err
		}
		p = p[toWrite:]
		w.curLineLen = 1 // Continuation lines are indented
	}
	written, err := w.w.Write(p)
	total += written
	w.curLineLen += written
	return total, err
}

// base64Writer ...
type base64Writer struct {
	w          io.Writer
	curLineLen int
	maxLineLen int
}

// Write ...
func (w *base64Writer) Write(p []byte) (int, error) {
	var total int
	for len(p)+w.curLineLen > w.maxLineLen {
		toWrite := w.maxLineLen - w.curLineLen
		written, err := w.w.Write(p[:toWrite])
		total += written
		if err != nil {
			return total, err
		}
		written, err = w.w.Write([]byte("\r\n"))
		total += written
		if err != nil {
			return total, err
		}
		p = p[toWrite:]
		w.curLineLen = 0
	}
	written, err := w.w.Write(p)
	total += written
	w.curLineLen += written
	return total, err
}

// leftTrimReader ...
type leftTrimReader struct {
	r    *bufio.Reader
	done bool
}

// Read ...
func (r *leftTrimReader) Read(p []byte) (n3 int, err3 error) {
	if r.done {
		// Delegate
		return r.r.Read(p)
	}
	// Peek and discard any whitespace, until we hit the first non-whitespace byte, then delegate
	r.r.Peek(1) // force a buffer load if empty
	maxBuffered := r.r.Buffered()
	if maxBuffered == 0 {
		r.done = true
		return r.r.Read(p)
	}
	peek, _ := r.r.Peek(maxBuffered)
	maxBuffered = len(peek)
	whiteSpaceCount := 0
	for whiteSpaceCount < maxBuffered && isASCIISpace(peek[whiteSpaceCount]) {
		whiteSpaceCount++
	}
	if whiteSpaceCount > 0 {
		discarded, err := r.r.Discard(whiteSpaceCount)
		if err == nil && discarded == whiteSpaceCount && whiteSpaceCount == maxBuffered {
			return r.Read(p)
		}
	}
	r.done = true
	return r.r.Read(p)
}

// isASCIISpace ...
func isASCIISpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
