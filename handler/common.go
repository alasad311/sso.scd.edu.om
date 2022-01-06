package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"sso.scd.edu.om/entity"
)

func EmptyResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": 501, "msg": "No Authorization header provided"})
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func LoggerToFile(message string) gin.HandlerFunc {
	logFilePath := "./"
	logFileName := "scdsso.log"
	// log file
	fileName := path.Join(logFilePath, logFileName)
	// check file
	if _, err := os.Lstat(fileName); err == nil {
		os.Remove(fileName)
	}
	// Write file
	src, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err", err)
	}
	// instantiation
	logger := logrus.New()
	// Set output
	logger.Out = src
	// Set log level
	logger.SetLevel(logrus.DebugLevel)
	// Set rotatelogs
	logWriter, _ := rotatelogs.New(
		// Split file name
		fileName+".%d%m%Y.log",
		// Generate soft chain, point to the latest log file
		rotatelogs.WithLinkName(fileName),
		// Set maximum save time (7 days)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// Set log cutting interval (1 day)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// Add Hook
	logger.AddHook(lfHook)

	return func(c *gin.Context) {

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// start time
		startTime := time.Now()
		// Processing request
		c.Next()
		responseBody := bodyLogWriter.body.String()
		var responseMsg string
		var responseData interface{}
		if responseBody != "" {
			res := entity.Result{}
			err := json.Unmarshal([]byte(responseBody), &res)
			if err == nil {
				responseMsg = res.Message
				responseData = res.Data
			}
		}
		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}
		// End time
		endTime := time.Now()
		// latency
		latencyTime := endTime.Sub(startTime)
		// Request mode
		reqMethod := c.Request.Method
		// Request routing
		reqUri := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// Request Size
		size := c.Writer.Size()
		// Request Agent
		reqAgent := c.Request.UserAgent()
		// Request IP
		clientIP := c.ClientIP()
		// Log format
		logger.WithFields(logrus.Fields{
			"status_code":       statusCode,
			"req_method":        reqMethod,
			"req_uri":           reqUri,
			"client_ip":         clientIP,
			"latency_time":      latencyTime,
			"req_size":          size,
			"request_post_data": c.Request.PostForm.Encode(),
			"req_agent":         reqAgent,
			"response_msg":      responseMsg,
			"response_data":     responseData,
		}).Info(message)

	}

}
