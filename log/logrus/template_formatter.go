package logrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"text/template"
	"unicode/utf8"
)

type Purifier interface {
	Purify(original, derivative string) string
}

// TemplateFormatter formats logs into text
type TemplateFormatter struct {
	Purifier Purifier

	// system that already adds timestamps.
	DisableTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// PadLevelText Adds padding the level text so that all the levels output at the same length
	// PadLevelText is a superset of the DisableLevelTruncation option
	PadLevelText bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &TemplateFormatter{
	//     FieldMap: FieldMap{
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"}}
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// The max length of the level text, generated dynamically on init
	levelTextMaxLength int

	// Template for formatting
	Template *template.Template

	// SystemFields are the fields that are added by the system and should be
	SystemFields map[string]struct{}
}

func (f *TemplateFormatter) init(entry *Entry) {
	// Get the max length of the level text
	for _, level := range AllLevels {
		levelTextLength := utf8.RuneCount([]byte(level.String()))
		if levelTextLength > f.levelTextMaxLength {
			f.levelTextMaxLength = levelTextLength
		}
	}
}

// Format renders a single log entry
func (f *TemplateFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	var funcVal, fileVal string

	fixedKeys := make([]string, 0, 4+len(data))
	if !f.DisableTimestamp {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyTime))
	}
	fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyLevel))
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyMsg))
	}
	if entry.err != "" {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyLogrusError))
	}
	if entry.HasCaller() {
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		} else {
			funcVal = entry.Caller.Function
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}

		if funcVal != "" {
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyFunc))
		}
		if fileVal != "" {
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyFile))
		}
	}

	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
			fixedKeys = append(fixedKeys, keys...)
		} else {
			fixedKeys = append(fixedKeys, keys...)
			f.SortingFunc(fixedKeys)

		}
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	systemFields := f.SystemFields
	if systemFields == nil {
		systemFields = defaultSystemFields
	}

	params := make(map[string]interface{})
	rest := make(map[string]interface{})
	var entity, action, method string
	var subject, body string

	for _, key := range fixedKeys {
		var value interface{}
		switch {
		case key == f.FieldMap.resolve(FieldKeyTime):
			value = entry.Time.Format(timestampFormat)
		case key == f.FieldMap.resolve(FieldKeyLevel):
			value = entry.Level.String()
		case key == f.FieldMap.resolve(FieldKeyMsg):
			value = f.purify(entry.Message)
		case key == f.FieldMap.resolve(FieldKeyLogrusError):
			value = entry.err
		case key == f.FieldMap.resolve(FieldKeyFunc) && entry.HasCaller():
			value = funcVal
		case key == f.FieldMap.resolve(FieldKeyFile) && entry.HasCaller():
			value = fileVal
		case key == f.FieldMap.resolve(FieldKeyTraceID):
			value, _ = data[key]
		case key == f.FieldMap.resolve(FieldKeyEntity):
			value, _ = data[key]
			entity = f.value2string(value)
			continue
		case key == f.FieldMap.resolve(FieldKeyAction):
			value, _ = data[key]
			action = f.value2string(value)
			continue
		case key == f.FieldMap.resolve(FieldKeyMethod):
			value, _ = data[key]
			method = f.value2string(value)
			continue
		case key == f.FieldMap.resolve(FieldKeySubject):
			value, _ = data[key]
			subject = f.value2string(value)
			continue
		case key == f.FieldMap.resolve(FieldKeyData):
			value, _ = data[key]
			body = f.value2string(value)
			continue
		default:
			value = data[key]
		}

		val := f.value2string(value)
		if _, ok := systemFields[key]; ok {
			params[key] = val
		} else {
			rest[key] = val
		}
	}

	params["entity"] = f.formatEntity(entity, action)
	params["event"] = f.formatEvent(method, subject, body)

	if _, ok := params[FieldKeyTraceID]; !ok {
		params[FieldKeyTraceID] = ""
	}

	if len(rest) > 0 {
		var details []byte
		details, _ = json.Marshal(rest)
		params["details"] = string(details)
	}

	tpl := f.Template
	if tpl == nil {
		tpl = defaultTpl
	}
	_ = tpl.Execute(b, params)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TemplateFormatter) value2string(value interface{}) string {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	return stringVal
}

func (f *TemplateFormatter) formatEntity(entity, action string) string {
	if entity == "" {
		return ""
	}

	var result bytes.Buffer
	result.WriteByte(' ')
	result.WriteString(entity)
	if action == "" {
		result.WriteByte(':')
	} else {
		result.WriteByte(' ')
		result.WriteString(action)
	}

	return result.String()
}

func (f *TemplateFormatter) formatEvent(method, subject, body string) string {
	if body == "" && subject == "" {
		return ""
	}

	var wantSpace bool
	var result bytes.Buffer
	result.WriteByte(' ')

	if method != "" {
		result.WriteString(method)
		wantSpace = true
	}

	if subject != "" {
		if wantSpace {
			result.WriteByte(' ')
		}
		result.WriteString(subject)
		wantSpace = true
	}

	if body != "" {
		if wantSpace {
			result.WriteByte(' ')
		}
		result.WriteString(f.purify(body))
	}

	return result.String()
}

func (f *TemplateFormatter) purify(s string) string {
	if f.Purifier == nil {
		return s
	}

	return f.Purifier.Purify(s, s)
}

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
}

var defaultTemplate = `{{.time}} {{.level | ToUpper}} #{{.trace_id}}:{{.entity}} {{.msg}}{{.event}}{{if .details}} DETAILS {{.details}}{{end}}`

var defaultTpl = template.Must(template.New("log").Funcs(funcMap).Parse(defaultTemplate))

var defaultSystemFields = map[string]struct{}{
	FieldKeyTime:    {},
	FieldKeyLevel:   {},
	FieldKeyMsg:     {},
	FieldKeyFile:    {},
	FieldKeyFunc:    {},
	FieldKeyTraceID: {},
	FieldKeyEntity:  {},
	FieldKeyAction:  {},
	FieldKeyMethod:  {},
	FieldKeySubject: {},
	FieldKeyData:    {},
}
