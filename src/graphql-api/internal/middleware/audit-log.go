package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"graphql-api/config"
	"graphql-api/internal/auth"
	"graphql-api/pkg/data/models"
	"graphql-api/pkg/graphql"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	// "sync"
	"time"

	"graphql-api/pkg/graphql/utils"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/mssola/user_agent"
	"github.com/samborkent/uuidv7"
)

// AuditLog represents an entry in the audit log.
type AuditLog struct {
	Timestamp   time.Time // Timestamp of the log entry
	Method      string    // HTTP method (POST, GET, etc.)
	URL         string    // Request URL
	RequestBody string    // Request body
	UserName    string    // User ID (if authenticated)
	Duration    time.Duration
	// Add more fields as needed
}

// CustomResponseWriter wraps the standard http.ResponseWriter
type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Status     string        // Captures the HTTP status code
	Body       *bytes.Buffer // Captures the response body
}

// Write captures the body content while writing to the underlying ResponseWriter
func (crw *CustomResponseWriter) Write(data []byte) (int, error) {
	crw.Body.Write(data)                  // Capture the response body
	return crw.ResponseWriter.Write(data) // Pass through to the original ResponseWriter
}

// WriteHeader captures the status code while writing the HTTP response
func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	crw.StatusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode) // Pass through to the original ResponseWriter
}

// Middleware function to log GraphQL requests.
func AuditLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := &CustomResponseWriter{
			ResponseWriter: w,
			Body:           new(bytes.Buffer),
			StatusCode:     http.StatusOK, // Default to 200 OK
		}

		// Parse the GraphQL request from JSON
		var gqlRequest graphql.GraphQLRequest
		body, err := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		if err == nil {
			if err = json.Unmarshal(body, &gqlRequest); err != nil {
				http.Error(w, "Invalid GraphQL request", http.StatusBadRequest)
				return
			}
		}

		// Call the next handler
		next.ServeHTTP(crw, r)
		writeLog(r, crw, gqlRequest.Query, start)
	})
}

func writeLog(r *http.Request, crw *CustomResponseWriter, query string, start time.Time) {
	logEntry := prepareLog(r, crw, query, start)

	fmt.Printf("Audit Log: %+v\n", logEntry)

	// Start the log writing Go routine
	go writeLogToFile(*logEntry)
}

// Write log to a file
func writeLogToFile(logEntry models.LogModel) {

	var currentLogFile *os.File
	var currentLogFilePath string

	
		// Determine the log file name based on the current timestamp (every 5 minutes)
		logFileName := fmt.Sprintf("%04d-%02d-%02d-%02d_%02d.log",
			logEntry.Timestamp.Year(),
			logEntry.Timestamp.Month(),
			logEntry.Timestamp.Day(),
			logEntry.Timestamp.Hour(),
			(logEntry.Timestamp.Minute()/5)*5,
		)

		// Get the absolute path
		relativePath := "../../logs"
		logFilePath := filepath.Join(relativePath, logFileName)
		absolutePath, err := filepath.Abs(logFilePath)
		if err != nil {
			fmt.Println("log error:", err)
		}

		fmt.Println("logfile", absolutePath)

		// Check if we're still using the same log file, or need to switch
		if currentLogFilePath != logFilePath {
			if currentLogFile != nil {
				currentLogFile.Close() // Close the previous log file
			}
			currentLogFilePath = absolutePath

			// Ensure the "logs" directory exists
			if err := os.MkdirAll(relativePath, 0755); err != nil {
				log.Fatalf("Error creating logs directory: %v", err)
			}

			var err error
			currentLogFile, err = os.OpenFile(absolutePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Error opening log file: %v", err)
			}
		}

		// Write the log entry as JSON to the file
		jsonData, err := json.Marshal(logEntry)
		if err != nil {
			log.Fatalf("Error marshaling log data: %v", err)
		}

		_, err = currentLogFile.Write(append(jsonData, '\n')) // Add newline
		if err != nil {
			log.Fatalf("Error writing to log file: %v", err)
		}
	
}

func prepareLog(r *http.Request, crw *CustomResponseWriter, bodyString string, start time.Time) *models.LogModel {

	uaString := r.Header.Get("User-Agent")
	token := r.Header.Get("Authorization")
	ip := r.RemoteAddr
	actions := transformGraphResolves(bodyString)
	// Parse the User-Agent
	ua := user_agent.New(uaString)
	// Get the browser name and version
	browserName, browserVersion := ua.Browser()
	// Get the operating system name
	osInfo := ua.OS()
	device := ua.Model()
	userId := getUserIdFromJWT(token)

	// Parse the GraphQL response
	var gqlResponse graphql.GraphQLResponse
	if err := json.NewDecoder(crw.Body).Decode(&gqlResponse); err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}
	errors := utils.ErrorsToString(gqlResponse.Errors)

	logData := &models.LogModel{
		LogId:                uuidv7.New().String(),
		Timestamp:            time.Now(),
		Duration:             time.Since(start),
		Status:               http.StatusText(crw.StatusCode),
		ClientIp:             ip,
		ClientBrowser:        browserName,
		ClientBrowserVersion: browserVersion,
		ClientOs:             osInfo,
		ClientOsVersion:      ua.OSInfo().Version,
		ClientDevice:         device,
		UserId:               userId,
		Actions:              actions,
		Resource:             "GraphQLApi",
		Errors:               errors,
	}

	return logData
}

func getUserIdFromJWT(token string) int {
	config := config.NewConfig()
	user, err := auth.DecodeJWTToken(token, config.SecretKey)
	userId := -1
	if err == nil {
		userId = user.UserId
	}

	return userId
}

func transformGraphResolves(query string) string {
	// Parse the GraphQL query into an AST
	document, err := utils.ParseGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Failed to parse GraphQL query: %v", err)
	}

	// List to store the dot notation for "resolve" fields
	var resolveDotNotation []string

	// Find the first operation and generate dot notation for resolve fields
	for _, definition := range document.Definitions {
		if operation, ok := definition.(*ast.OperationDefinition); ok {
			utils.GenerateResolveDotNotation(operation.SelectionSet, "", &resolveDotNotation)
		}
	}

	methods := strings.Join(resolveDotNotation[:], ",")
	return methods
}