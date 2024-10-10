package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	infoLogger.Println(message)
	fmt.Printf("%s [INFO] %s \n", time.Now().Format("2006-01-02 15:04:05"), message)

}

func Error(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	errorLogger.Println(message)
	fmt.Printf("%s [ERROR] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}
