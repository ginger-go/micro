package logger

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginger-go/micro"
	"github.com/ginger-go/micro/plugins/apicall"
	"github.com/ginger-go/micro/plugins/auth"
)

type sendLogRequest struct {
	Filename string `uri:"filename" binding:"required"`
	Content  string `uri:"content" binding:"required"`
}

func sendLog() {
	now := time.Now()
	if now.Hour() == 0 && now.Minute() >= 0 && now.Minute() <= 3 {
		yes := now.AddDate(0, 0, -1)
		filename := yes.Format("2006-01-02") + ".log"
		content := getLogFromAndToAndLevel(uint(yes.Unix()), uint(now.Unix()), "ALL")
		sendLogToLogService(filename, content)
	}
	filename := now.Format("2006-01-02") + ".log"
	content := getLogFromAndToAndLevel(uint(now.Unix()), uint(now.Unix()), "ALL")
	sendLogToLogService(filename, content)
}

func sendLogToLogService(filename, content string) {
	_, err := apicall.POST[struct{}](LOG_SERVICE_IP+"/micro/log", &sendLogRequest{
		Filename: filename,
		Content:  content,
	}, map[string]string{
		"Authorization": "Bearer " + auth.SYSTEM_TOKEN,
	}, "", nil)
	if err != nil {
		panic("failed to send log to log service")
	}
}

type getLogger struct {
	From  uint   `uri:"from" binding:"required"`
	To    uint   `uri:"to" binding:"required"`
	Level string `uri:"level" binding:"required"` // all, info, error, warn, debug
}

func getLog(ctx *gin.Context) {
	req := micro.GinRequest[getLogger](ctx)
	ctx.String(200, getLogFromAndToAndLevel(req.From, req.To, strings.ToUpper(req.Level)))
}

func getLogFromAndToAndLevel(from uint, to uint, level string) string {
	files, err := os.ReadDir(LOG_FOLDER)
	if err != nil {
		return ""
	}
	fromTime := time.Unix(int64(from), 0)
	toTime := time.Unix(int64(to), 0)

	var selectedFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileTime, err := time.Parse("2006-01-02", strings.ReplaceAll(file.Name(), ".log", ""))
		if err != nil {
			continue
		}
		if fileTime.After(fromTime) && fileTime.Before(toTime) {
			selectedFiles = append(selectedFiles, file.Name())
		}
	}

	var logs []string
	for _, file := range selectedFiles {
		logs = append(logs, getLogFromFile(file, level))
	}

	return strings.Join(logs, "\n")
}

func getLogFromFile(filename, level string) string {
	file, err := os.ReadFile(LOG_FOLDER + "/" + filename)
	if err != nil {
		return ""
	}
	var logs []string
	for _, line := range strings.Split(string(file), "\n") {
		if strings.Contains(line, level) || level == "ALL" {
			logs = append(logs, line)
		}
	}
	return strings.Join(logs, "\n")
}
