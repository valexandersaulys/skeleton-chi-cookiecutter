package config

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Flagged variables
var LogLevel = flag.String("log-level", "INFO",
	"Specify the log level at runtime, can be TRACE, DEBUG, INFO, WARN, or ERROR. Defaults to INFO.")
var LogPath = flag.String("log-path", "/tmp/{{cookiecutter.application_name}}.logs", "Specify the path for (most) logs")
var PrintRoutes = flag.Bool("routes", false, "Generate router documentation")
var RunMigrations = flag.Bool("migrate", true, "Whether to run model migrations")
var RuntimePort = flag.Int("port", 3000, "Specify the runtime port to run on")
var Timeout = flag.Int("timeout", -1, "Specify timeout before 504 Gateway Timeout error to client")
var MaxIdleDbConnections = flag.Int("max-idle-db-connections", 10, "Set the max number of idle database connections")
var MaxOpenDbConnections = flag.Int("max-open-db-connections", 100, "Set the max number of open database connections")

// var RawSqlQuery = flag.String("sql-query", "", "If specified, this will execute a sql query")

var RUNTIME_ENVIRONMENTS = []string{"TESTING", "LOCAL", "DEV", "PROD"}

// Environmental Variables
var Environment *string
var AddDummyModels *bool
var CsrfProtectionKey *[]byte
var CookieStoreSessionKeyAuthentication *[]byte
var CookieStoreSessionKeyEncryption *[]byte

func RunInit() {
	flag.Parse()

	log_path := fmt.Sprintf("%s.%s", *LogPath, time.Now().Format("20060102"))
	logFile, _ := os.OpenFile(
		log_path,
		os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	log.SetOutput(logFile)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   false,
		FullTimestamp: true,
	})
	switch *LogLevel {
	case "TRACE", "trace":
		log.SetLevel(log.TraceLevel)
	case "DEBUG", "debug":
		log.SetLevel(log.DebugLevel)
	case "INFO", "info":
		log.SetLevel(log.InfoLevel)
	case "WARN", "warn":
		log.SetLevel(log.WarnLevel)
	case "ERROR", "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
	log.Error(fmt.Sprintf("Setting log level to %d", log.GetLevel()))

	_environment := strings.ToUpper(os.Getenv("RUNTIME_ENV"))
	if !slices.Contains(RUNTIME_ENVIRONMENTS, _environment) {
		errorLine := fmt.Sprintf("Runtime environment `RUNTIME_ENV` is not one of the following: %v", RUNTIME_ENVIRONMENTS)
		panic(errorLine)
	}
	Environment = &_environment
	_addDummyModels, _ := strconv.ParseBool(os.Getenv("ADD_DUMMIES"))
	AddDummyModels = &_addDummyModels
	_csrf_protection_key := []byte(os.Getenv("CSRF_PROTECTION_KEY"))
	CsrfProtectionKey = &_csrf_protection_key

	var _cookie_store_session_key_auth []byte
	var _cookie_store_session_key_en []byte

	if os.Getenv("COOKIE_STORE_AUTH_KEY") != "" {
		_cookie_store_session_key_auth = []byte(os.Getenv("COOKIE_STORE_AUTH_KEY"))
	} else {
		_cookie_store_session_key_auth = []byte("COOKIE_STORE_AUTH_KEY")
	}
	CookieStoreSessionKeyAuthentication = &_cookie_store_session_key_auth

	if os.Getenv("COOKIE_STORE_ENCRYPT_KEY") != "" {
		_cookie_store_session_key_en = []byte(os.Getenv("COOKIE_STORE_ENCRYPT_KEY"))
	} else {
		_cookie_store_session_key_en = []byte("COOKIE_STORE_ENCRYPT_KEY")
	}
	CookieStoreSessionKeyEncryption = &_cookie_store_session_key_en
}
