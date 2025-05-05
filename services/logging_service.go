package services

import (
	"global-auth-server/libs"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// LoggingService is the structure that provides the logging service.
type LoggingService struct {
	logSender *libs.LogSender
}

var (
	loggingServiceInstance *LoggingService
	once                   sync.Once
)

// Newloggingservice creates and initializes the only instance of loggingservice.
// The logger configuration is defined here.
func NewLoggingService() *LoggingService {
	once.Do(func() {
		_ = godotenv.Load()

		LOG_URL_API := os.Getenv("LOG_URL_API")

		config := libs.LogSenderConfig{
			APIURL: LOG_URL_API + "/api/logs/create",
			QueueCapacity: 1000,                              
			RetryDelay:    5 * time.Second,                  
			MaxRetries:    3,                               
			BatchSize:     10,                              
			BatchInterval: 2 * time.Second,               
		}
		loggingServiceInstance = &LoggingService{
			logSender: libs.GetLogSender(config),
		}
	})
	return loggingServiceInstance
}

// Log Calls the logsender log method with the data provided.
func (ls *LoggingService) Log(userID any, url string, payload any, response any, action string) {
	ls.logSender.Log(userID, url, payload, response, action)
}

// Stop stops the logging service and waiting for pending logs to be sent.
func (ls *LoggingService) Stop() {
	if ls.logSender != nil {
		ls.logSender.Stop()
	}
}