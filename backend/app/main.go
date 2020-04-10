package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/umputun/remark/backend/app/cmd"
)

// Opts with all cli commands and flags
type Opts struct {
	ServerCmd  cmd.ServerCommand  `command:"server"`
	ImportCmd  cmd.ImportCommand  `command:"import"`
	BackupCmd  cmd.BackupCommand  `command:"backup"`
	RestoreCmd cmd.RestoreCommand `command:"restore"`
	AvatarCmd  cmd.AvatarCommand  `command:"avatar"`
	CleanupCmd cmd.CleanupCommand `command:"cleanup"`
	RemapCmd   cmd.RemapCommand   `command:"remap"`

	RemarkURL    string `long:"url" env:"REMARK_URL" required:"true" description:"url to remark"`
	SharedSecret string `long:"secret" env:"SECRET" required:"true" description:"shared secret key"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"debug mode"`
}

var revision = "dev"

var logger *zap.Logger

func main() {
	fmt.Printf("remark42 %s\n", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		initLogger(opts.Dbg)
		// commands implements CommonOptionsCommander to allow passing set of extra options defined for all commands
		c := command.(cmd.CommonOptionsCommander)
		c.SetCommon(cmd.CommonOpts{
			RemarkURL:    opts.RemarkURL,
			SharedSecret: opts.SharedSecret,
			Revision:     revision,
		})
		for _, entry := range c.HandleDeprecatedFlags() {
			log.Printf("[WARN] --%s is deprecated and will be removed in v%s, please use --%s instead",
				entry.Old, entry.RemoveVersion, entry.New)
		}
		err := c.Execute(args)
		if err != nil {
			log.Printf("[ERROR] failed with %+v", err)
		}
		return err
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

// levels from github.com/go-pkgz/lgr/logger.go
var gopkgzlgrLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "PANIC", "FATAL"}

func gopkgzlgrLevelToZapFunc(logger *zap.SugaredLogger, lvl string) func(template string, args ...interface{}) {
	switch lvl {
	case "TRACE":
		fallthrough
	case "DEBUG":
		return logger.Debugf
	case "INFO":
		return logger.Infof
	case "WARN":
		return logger.Warnf
	case "ERROR":
		return logger.Errorf
	case "PANIC":
		return logger.Panicf
	case "FATAL":
		return logger.Fatalf
	}
	return logger.Infof
}

// extractLevel parses messages with optional level prefix and returns level and the message with stripped level
func (l *loggerWriter) extractLevel(line string) (level, msg string) {
	for _, lv := range gopkgzlgrLevels {
		if strings.HasPrefix(line, lv) {
			return lv, strings.TrimSpace(line[len(lv):])
		}
		if strings.HasPrefix(line, "["+lv+"]") {
			return lv, strings.TrimSpace(line[len("["+lv+"]"):])
		}
	}
	return "INFO", line
}

// remark42's original log setup func
func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup()
	log.Setup(log.Msec, log.LevelBraces)
}

func redirectGopkgzlgrLogAt(l *zap.Logger, dbg bool) func() {
	// caller skip 3 level for github.com/go-pkgz/lgr v0.6.3
	logger := l.WithOptions(zap.AddCallerSkip(3)).Sugar()
	logFunc := logger.Infof
	// log.CallerPkg log.CallerFile, log.CallerFunc, log.Msec are disabled by default
	log.Setup()
	// we need enable log.Debug Option to get the debug level output
	log.Setup(log.Debug,
		log.LevelBraces,
		log.Format(`{{.Level}} {{.Message}}`), // only ouput level and message
		log.Out(&loggerWriter{logger, logFunc}), // hacking
		log.Err(&loggerWriter{logger, logFunc})) // hacking
	return func() {
		setupLog(dbg)
	}
}

type loggerWriter struct {
	logger *zap.SugaredLogger
	logFunc func(template string, args ...interface{})
}

func (l *loggerWriter) Write(p []byte) (int, error) {
	pLen := len(p)
	p = bytes.TrimSpace(p)
	lvl, msg := l.extractLevel(string(p))
	//l.logger.Warnf("lvl: %s, msg: %s", lvl, msg)
	gopkgzlgrLevelToZapFunc(l.logger, lvl)(msg)
	return pLen, nil
}

func initLogger(dbg bool) func() {
	// prod stackLevel := ErrorLevel
	// dev stackLevel = WarnLevel
	zapCfg := zap.NewDevelopmentConfig()
	//zapCfg := zap.NewProductionConfig()
	//zapCfg.Encoding = "console"
	if !dbg {
		zapCfg.DisableCaller = true
	}
	// if Development , stackLevel = WarnLevel
	zapCfg.Development = false
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	tmpLogger, _ := zapCfg.Build()
	logger = tmpLogger.Named("[remark42]")
	//The default global logger used by zap.L() and zap.S() is a no-op logger.
	//To configure the global loggers, you must use ReplaceGlobals.
	zap.ReplaceGlobals(logger)
	undo := redirectGopkgzlgrLogAt(logger, dbg)
	if dbg {
		zap.S().Infof("debug enabled")
	}
	return func() {
		logger.Sync() // flushes buffer, if any
		undo()
	}
}

// getDump reads runtime stack and returns as a string
func getDump() string {
	maxSize := 5 * 1024 * 1024
	stacktrace := make([]byte, maxSize)
	length := runtime.Stack(stacktrace, true)
	if length > maxSize {
		length = maxSize
	}
	return string(stacktrace[:length])
}

func init() {
	// catch SIGQUIT and print stack traces
	sigChan := make(chan os.Signal)
	go func() {
		for range sigChan {
			log.Printf("[INFO] SIGQUIT detected, dump:\n%s", getDump())
		}
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)
}
