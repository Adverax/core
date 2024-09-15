package log

import (
	"bytes"
	"compress/gzip"
	json2 "encoding/json"
	"fmt"
	"github.com/adverax/core/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type fileInfo struct {
	path    string
	modTime time.Time
}

type ChunkManagerOptions struct {
	MaxCount int
	MaxAge   int
	Folder   string
}

type ChunkManager struct {
	options   ChunkManagerOptions
	files     []*fileInfo
	mu        sync.Mutex
	cleanupCh chan struct{}
}

func NewChunkManager(options ChunkManagerOptions) *ChunkManager {
	_ = os.MkdirAll(options.Folder, 0755)

	cb := &ChunkManager{
		options:   options,
		cleanupCh: make(chan struct{}, 1),
	}

	go cb.cleanupWorker()

	return cb
}

func (that *ChunkManager) Save(data string) string {
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), that.getExtension(data))
	filePath := filepath.Join(that.options.Folder, filename)
	err := os.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return err.Error()
	}
	defer os.Remove(filePath)

	compressedFilePath, err := that.compressFile(filePath)
	if err != nil {
		return err.Error()
	}

	that.addFile(&fileInfo{path: compressedFilePath, modTime: time.Now()})

	select {
	case that.cleanupCh <- struct{}{}:
	default:
	}

	return fmt.Sprintf("CHUNK: %s", filepath.Base(compressedFilePath))
}

func (that *ChunkManager) compressFile(filePath string) (string, error) {
	compressedFilePath := filePath + ".gz"
	inFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer inFile.Close()

	outFile, err := os.Create(compressedFilePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, inFile)
	if err != nil {
		return "", err
	}

	return compressedFilePath, nil
}

func (that *ChunkManager) addFile(file *fileInfo) {
	that.mu.Lock()
	defer that.mu.Unlock()

	that.files = append(that.files, file)
}

func (that *ChunkManager) cleanupWorker() {
	for range that.cleanupCh {
		that.cleanupOldFiles()
	}
}

func (that *ChunkManager) cleanupOldFiles() {
	that.mu.Lock()
	defer that.mu.Unlock()

	// Sort files by modification time
	sort.Slice(that.files, func(i, j int) bool {
		return that.files[i].modTime.Before(that.files[j].modTime)
	})

	// Remove oldest files if count exceeds maxCount
	if len(that.files) > that.options.MaxCount {
		for _, file := range that.files[:len(that.files)-that.options.MaxCount] {
			os.Remove(file.path)
		}
		that.files = that.files[len(that.files)-that.options.MaxCount:]
	}

	// Remove files older than maxAge
	cutoff := time.Now().Add(-time.Duration(that.options.MaxAge) * time.Second)
	for i := 0; i < len(that.files); {
		if that.files[i].modTime.Before(cutoff) {
			os.Remove(that.files[i].path)
			that.files = append(that.files[:i], that.files[i+1:]...)
		} else {
			i++
		}
	}
}

func (that *ChunkManager) getExtension(s string) string {
	if isPNG(s) {
		return ".png"
	}

	if isJSON(s) {
		return ".json"
	}

	return ".txt"
}

type ChunkStorage interface {
	Save(data string) string
}

type PNGPurifier struct {
	storage ChunkStorage
	next    Purifier
}

func NewPNGPurifier(storage ChunkStorage, next Purifier) *PNGPurifier {
	return &PNGPurifier{
		storage: storage,
		next:    next,
	}
}

func (that *PNGPurifier) Purify(original, derivative string) string {
	if isPNG(derivative) {
		return that.storage.Save(derivative)
	}

	if that.next == nil {
		return derivative
	}

	return that.next.Purify(original, derivative)
}

type LenPurifier struct {
	maxLen  int
	storage ChunkStorage
	next    Purifier
}

func NewLenPurifier(storage ChunkStorage, maxLength int, next Purifier) *LenPurifier {
	return &LenPurifier{
		maxLen:  maxLength,
		storage: storage,
		next:    next,
	}
}

func (that *LenPurifier) Purify(original, derivative string) string {
	if len(derivative) <= that.maxLen {
		if that.next == nil {
			return derivative
		}

		return that.next.Purify(original, derivative)
	}

	return that.storage.Save(derivative)
}

type MultilinePurifier struct {
	next Purifier
}

func NewMultilinePurifier(next Purifier) *MultilinePurifier {
	return &MultilinePurifier{next: next}
}

func (that *MultilinePurifier) Purify(original, derivative string) string {
	derivative = that.purify(derivative)

	if that.next == nil {
		return derivative
	}

	return that.next.Purify(original, derivative)
}

func (that *MultilinePurifier) purify(s string) string {
	if s2, ok := that.purifyAsJson(s); ok {
		return s2
	}

	return that.purifyAsPlain(s)
}

func (that *MultilinePurifier) purifyAsJson(s string) (string, bool) {
	if !canBeJson.MatchString(s) {
		return s, false
	}

	var obj interface{}
	if err := json.Unmarshal([]byte(s), &obj); err == nil {
		b, _ := json.Marshal(obj)
		return string(b), true
	}

	return s, false
}

func (that *MultilinePurifier) purifyAsPlain(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Trim(line, "\n\r\t ")
	}

	return strings.Join(lines, " ")
}

func isPNG(s string) bool {
	return len(s) >= 4 && s[0] == 0x89 && s[1] == 0x50 && s[2] == 0x4E && s[3] == 0x47
}

var canBeJson = regexp.MustCompile(`^\s*[\[{"0-9]`)

func isJSON(s string) bool {
	if !canBeJson.MatchString(s) {
		return false
	}

	decoder := json2.NewDecoder(bytes.NewReader([]byte(s)))
	for {
		_, err := decoder.Token()
		if err != nil {
			return err == nil
		}
	}
}
