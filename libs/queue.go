package libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// LogEntry representa una entrada de log que se enviará a la API.
type LogEntry struct {
	UserID   any    `json:"user_id"` // Puede ser nil
	URL      string `json:"url"`
	Payload  any    `json:"payload"`  // JSON indistinto
	Response any    `json:"response"` // JSON indistinto
	Action   string `json:"action"`
}

// LogSenderConfiguración para el LogSender.
type LogSenderConfig struct {
	APIURL        string
	QueueCapacity int
	RetryDelay    time.Duration
	MaxRetries    int
	BatchSize     int
	BatchInterval time.Duration
}

// LogSender es la estructura principal para enviar logs a una API.
type LogSender struct {
	config LogSenderConfig
	queue  chan LogEntry
	wg     sync.WaitGroup
	stop   chan struct{}
}

var (
	logSenderInstance *LogSender
	onceLogSender     sync.Once
)

// GetLogSender devuelve la única instancia de LogSender.
func GetLogSender(config LogSenderConfig) *LogSender {
	onceLogSender.Do(func() {
		logSenderInstance = NewLogSender(config)
	})
	return logSenderInstance
}

// NewLogSender crea una nueva instancia de LogSender.
func NewLogSender(config LogSenderConfig) *LogSender {
	ls := &LogSender{
		config: config,
		queue:  make(chan LogEntry, config.QueueCapacity),
		stop:   make(chan struct{}),
	}
	ls.startWorker()
	return ls
}

// Log lock a new log message to be sent.
func (ls *LogSender) Log(userID any, url string, payload any, response any, action string) {
	entry := LogEntry{
		UserID:   userID,
		URL:      url,
		Payload:  payload,
		Response: response,
		Action:   action,
	}
	select {
	case ls.queue <- entry:
		// LOG CITED correctly
	default:
		fmt.Println("Warning: Log queue is full, dropping log:", entry)
	}
}

// Startworker starts a goroutine that processes Logs's tail.
func (ls *LogSender) startWorker() {
	ls.wg.Add(1)
	go func() {
		defer ls.wg.Done()
		var batch []LogEntry
		ticker := time.NewTicker(ls.config.BatchInterval)
		defer ticker.Stop()

		for {
			select {
			case logEntry := <-ls.queue:
				batch = append(batch, logEntry)
				if len(batch) >= ls.config.BatchSize {
					ls.sendLogs(batch)
					batch = nil
				}
			case <-ticker.C:
				if len(batch) > 0 {
					ls.sendLogs(batch)
					batch = nil
				}
			case <-ls.stop:
				if len(batch) > 0 {
					ls.sendLogs(batch)
				}
				fmt.Println("Log sender worker stopped.")
				return
			}
		}
	}()
}

func (ls *LogSender) sendLogs(logs []LogEntry) {
    emptyPayload := map[string]any{}
    token, _, err := GenerateJWT(emptyPayload, 1*24*time.Hour)
    if err != nil {
        fmt.Println("Error generating JWT:", err)
        return
    }

    client := &http.Client{}
    for _, entry := range logs {
        jsonData, err := json.Marshal(entry)
        if err != nil {
            fmt.Println("Error marshaling log entry:", err)
            continue // jump this entrance and continue with the following
        }

        var retries int
        for retries < ls.config.MaxRetries {
            req, err := http.NewRequest("POST", ls.config.APIURL, bytes.NewBuffer(jsonData))
            if err != nil {
                // ... error handling ...
                continue
            }

            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

            resp, err := client.Do(req)
            if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
                defer resp.Body.Close()
                break
            }

            fmt.Printf("Error sending log entry (attempt %d): Status Code: %d. Retrying in %s...\n",
                retries+1, resp.StatusCode, ls.config.RetryDelay)
            if resp != nil {
                defer resp.Body.Close()
            }
            time.Sleep(ls.config.RetryDelay)
            retries++
        }
        if retries >= ls.config.MaxRetries {
            fmt.Println("Failed to send log entry after multiple retries:", entry)
        }
    }
    fmt.Println("Finished attempting to send all logs.")
}

// Stop detiene el worker de forma segura.
func (ls *LogSender) Stop() {
	close(ls.stop)
	ls.wg.Wait()
	fmt.Println("Log sender stopped.")
}
