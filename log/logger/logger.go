package logger

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	backupTimeFormat = "2006-01-02T15-04-05.000" //формат времени для имени файла бекапа
	defaultMaxSize   = 100                       //максимальный размен лога по умолчанию
	megabyte         = 1024 * 1024               //количество байт в мегабайте. так код кажется читабельней
)

// Logger - система учета и записи логов с удаление старых логов
// реализована поддержка интерфейса io.WriteCloser для полной совместимости со стандартным пакетом Log
type Logger struct {
	// имя файла лога с полным путем. если путь не указан - в текущей папке создается папка "log" и файл как имя программі + ".log"
	FileName string
	// максимальный размер файла лога. при превышении - текущий бекапим и создаем новый
	MaxSize int `json:"maxsize"`
	// максимальное время жизни бекапа в днях
	MaxAge int `json:"maxage"`
	// максимальное число бекапов лога
	MaxBackups int `json:"maxbackups"`
	// использовать локальное время компа или UTC ддля формирования имени файла
	LocalTime bool `json:"localtime"`

	size int64
	file *os.File
	mu   sync.Mutex
}

// InitLog - initialization logger
func InitLog(fileName string, logger *Logger) {
	log.SetOutput(NewLogger(fileName, logger.MaxSize, logger.MaxAge, logger.MaxBackups))
}

// NewLogger - создание нового экземпляра Logger
func NewLogger(fileName string, maxSize, maxAge, maxBackups int) *Logger {

	return &Logger{
		FileName:   fileName,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  false,
	}
}

func (l *Logger) openExistingOrNew(writeLen int) error {
	filename := l.fileName()
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if info.Size()+int64(writeLen) >= l.maxSize() {
		return l.rotate()
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return l.openNew()
	}
	l.file = file
	l.size = info.Size()
	return nil
}

func (l *Logger) compress(name string) {
	inFile := name
	outFile := name + ".gz"
	zipHandle, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("[CompressFile ERROR] Opening file:", err)
	}
	defer zipHandle.Close()
	zipWriter, err := gzip.NewWriterLevel(zipHandle, 9)
	if err != nil {
		fmt.Println("[CompressFile ERROR] New gzip writer:", err)
	}
	var b []byte
	inReader, err := os.OpenFile(inFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("[CompressFile ERROR] Opening old log file:", err)
		return
	}
	defer inReader.Close()
	b, err = ioutil.ReadAll(inReader)
	if err != nil {
		fmt.Println("[CompressFile ERROR] Read old log:", err)
		return
	}
	_, err = zipWriter.Write(b)
	if err != nil {
		fmt.Println("[CompressFile ERROR] Writing:", err)
		return
	}
	err = zipWriter.Close()
	if err != nil {
		fmt.Println("[CompressFile ERROR] Closing zip writer:", err)
		return
	}
}
func (l *Logger) removeFile(name string) {
	err := os.Remove(name)
	if err != nil {
		fmt.Println("Delete old log fileError :", name, "  ", err)
		return
	}
}
func (l *Logger) openNew() error {
	err := os.MkdirAll(l.dir(), 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := l.fileName()
	mode := os.FileMode(0644)
	info, err := os.Stat(name)
	if err == nil {
		mode = info.Mode()
		newname := backupName(name, l.LocalTime)
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}
		l.compress(newname)
		l.removeFile(newname)
		// работает только в  linux
		if err := chown(name, info); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	l.file = f
	l.size = 0
	return nil
}

// получение времени вынесено в отдельную ф-цию исключительно для отладки
func currentTime() time.Time {
	return time.Now()
}

// backupName - возвращает имя файла бекапа
func backupName(name string, local bool) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	t := currentTime()
	if !local { //нормализация времени конечно изврат, но выгдить должно лучше
		t = t.UTC()
	}
	timestamp := t.Format(backupTimeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

// запись данных в файл
func (l *Logger) write(p []byte) (n int, err error) {
	writeLen := int64(len(p))
	if writeLen > l.maxSize() {
		return 0, fmt.Errorf(
			"Размер записи(%d) превышает максимальный размер лога(%d)", writeLen, l.maxSize(),
		)
	}

	if l.file == nil {
		if err = l.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}
	if l.size+writeLen > l.maxSize() {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = l.file.Write(p)
	l.size += int64(n)
	return n, err
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	n, err = l.write(p)
	return n, err
}

func (l *Logger) createDefaultFileName() string {
	dir, err := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "log"))
	if err != nil {
		panic(fmt.Errorf("Error get abs path:%v", err))
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("Error create log dir:%v", err))
	}
	file := filepath.Base(os.Args[0]) + ".log"
	return filepath.Join(dir, file)
}

// filename - возвращает имя текущего файла лога
func (l *Logger) fileName() string {
	if l.FileName != "" {
		return l.FileName
	}
	return l.createDefaultFileName()
}

// dir - возвращает каталог текущего файла лога
func (l *Logger) dir() string {
	return filepath.Dir(l.fileName())
}

// maxSize - вычисление максимального размера лога.
// вычислять только в момент создания логгера - плохая практика,
// а так можно менять максимальный размер в процессе работы логгера
func (l *Logger) maxSize() int64 {
	if l.MaxSize == 0 {
		return int64(defaultMaxSize * megabyte)
	}
	return int64(l.MaxSize) * int64(megabyte)
}

func (l *Logger) rotate() error {
	if err := l.close(); err != nil {
		return err
	}

	if err := l.openNew(); err != nil {
		return err
	}
	return l.cleanup()
}

// Rotate - бекапим текущий файл и создаем новый...и чистим старые файлы
func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotate()
}

// Close - закрытие файла-необходимо для поддержки интерфейса io.Closer
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

// закрываем текущий файл
func (l *Logger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

// чистка ненужных бекапов
func (l *Logger) cleanup() error {
	if l.MaxBackups == 0 && l.MaxAge == 0 {
		return nil
	}

	files, err := l.oldLogFiles()
	if err != nil {
		return err
	}

	var deletes []logInfo

	if l.MaxBackups > 0 && l.MaxBackups < len(files) {
		deletes = files[l.MaxBackups:]
		files = files[:l.MaxBackups]
	}
	if l.MaxAge > 0 {
		diff := time.Duration(int64(24*time.Hour) * int64(l.MaxAge))

		cutoff := currentTime().Add(-1 * diff)

		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				deletes = append(deletes, f)
			}
		}
	}

	if len(deletes) == 0 {
		return nil
	}
	// удаляем все файлы по созданному списку
	go deleteAll(l.dir(), deletes)

	return nil
}

func deleteAll(dir string, files []logInfo) {
	for _, f := range files {
		_ = os.Remove(filepath.Join(dir, f.Name()))
	}
}

func (l *Logger) prefixAndExt() (prefix, ext string) {
	filename := filepath.Base(l.fileName())
	ext = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(ext)] + "-"
	return prefix, ext
}

func (l *Logger) timeFromName(filename, prefix, ext string) string {
	if !strings.HasPrefix(filename, prefix) {
		return ""
	}
	filename = filename[len(prefix):]
	if strings.HasSuffix(filename, ext+".gz") {
		ext = ext + ".gz" //коректировка определения суфикса у упакованых файлов
	}
	if !strings.HasSuffix(filename, ext) {
		return ""
	}
	filename = filename[:len(filename)-len(ext)]
	return filename
}

// возвращает список бекапов
func (l *Logger) oldLogFiles() ([]logInfo, error) {
	files, err := ioutil.ReadDir(l.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	logFiles := []logInfo{}

	prefix, ext := l.prefixAndExt()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := l.timeFromName(f.Name(), prefix, ext)
		if name == "" {
			continue
		}
		t, err := time.Parse(backupTimeFormat, name)
		if err == nil {
			logFiles = append(logFiles, logInfo{t, f})
		}
	}

	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

// logInfo и byFormatTime для более наглядного преобразования имени файла обратно в дату и сортировки списка
type logInfo struct {
	timestamp time.Time
	os.FileInfo
}

type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
