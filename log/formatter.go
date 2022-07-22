package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/kr/logfmt"
)

type PrettyFormatter struct {
	Out           io.Writer
	DisableColors bool
}

func (f PrettyFormatter) Write(p []byte) (int, error) {
	if string(p) == "\n" {
		return 0, nil
	}
	var (
		fields  []string
		level   string
		message string
	)

	handler := func(key, val []byte) error {
		// v := strings.TrimSuffix(strings.TrimPrefix(string(val), "'"), "'")
		v := string(val)
		switch k := string(key); k {
		case "lvl":
			level = v
		case "msg":
			message = v
		case "time":
			return nil
		default:
			fields = append(fields, fmt.Sprintf(`%s="%s"`, k, v))
		}
		return nil
	}
	if err := logfmt.Unmarshal(p, logfmt.HandlerFunc(handler)); err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	switch level {
	case "debug":
		cW(f.Out, !f.DisableColors, nWhite, "[%s] ", "DEBU")
	case "info":
		cW(f.Out, !f.DisableColors, nGreen, "[%s] ", "INFO")
	case "warn":
		cW(f.Out, !f.DisableColors, nYellow, "[%s] ", "WARN")
	case "error":
		cW(f.Out, !f.DisableColors, nRed, "[%s] ", "ERR ")
	}

	fmt.Fprintf(f.Out, "%-44s %s\n", message, strings.Join(fields, " "))

	return 0, nil
}

type JSONFormatter struct {
	Out io.Writer
}

func (f JSONFormatter) Write(p []byte) (int, error) {
	if string(p) == "\n" {
		return 0, nil
	}

	data := map[string]string{}

	handler := func(key, val []byte) error {
		// v := strings.TrimSuffix(strings.TrimPrefix(string(val), "'"), "'")
		data[string(key)] = string(val)
		return nil
	}
	if err := logfmt.Unmarshal(p, logfmt.HandlerFunc(handler)); err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	fmt.Fprintln(f.Out, string(raw))

	return 0, nil
}
