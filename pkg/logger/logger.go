package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidLevel      = errors.New("invalid level")
	ErrInvalidOutputType = errors.New("invalid output type")
	ErrInvalidFormatType = errors.New("invalid format type")
)

type OutputType string

const (
	OutputStdout OutputType = "stdout"
	OutputStderr OutputType = "stderr"
	OutputFile   OutputType = "file"
)

type FormatType string

const (
	FormatJSON FormatType = "json"
	FormatText FormatType = "text"
)

type Config struct {
	Level  int
	Output OutputType
	Format FormatType
}

func InitWithDefaults() error {
	return Init(Config{
		Level:  int(zerolog.DebugLevel),
		Output: OutputStdout,
		Format: FormatJSON,
	})
}

func Init(cfg Config) error {
	level := zerolog.Level(cfg.Level)
	if level < zerolog.TraceLevel || level > zerolog.PanicLevel {
		return fmt.Errorf("%w: %d", ErrInvalidLevel, cfg.Level)
	}

	output, err := getOutput(cfg.Output)
	if err != nil {
		return err
	}

	writer := getWriter(output, cfg.Format)

	log.Logger = zerolog.New(writer).
		With().
		Timestamp().
		Logger().
		Level(level)

	zerolog.SetGlobalLevel(level)
	return nil
}

func ParseOutputType(s string) (OutputType, error) {
	switch strings.ToLower(s) {
	case "stdout":
		return OutputStdout, nil
	case "stderr":
		return OutputStderr, nil
	case "file":
		return OutputFile, nil
	default:
		return OutputStdout, fmt.Errorf("%w: %s", ErrInvalidOutputType, s)
	}
}

func ParseFormatType(s string) (FormatType, error) {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON, nil
	case "text":
		return FormatText, nil
	default:
		return FormatJSON, fmt.Errorf("%w: %s", ErrInvalidFormatType, s)
	}
}

func getWriter(output io.Writer, format FormatType) io.Writer {
	switch format {
	case FormatText:
		return zerolog.ConsoleWriter{Out: output}
	case FormatJSON:
		return output
	default:
		return output
	}
}

func getOutput(outputType OutputType) (io.Writer, error) {
	switch outputType {
	case OutputStdout:
		return os.Stdout, nil
	case OutputStderr:
		return os.Stderr, nil
	case OutputFile:
		file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return file, nil
	default:
		return os.Stdout, nil
	}
}
