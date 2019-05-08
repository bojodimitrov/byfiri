package logger

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Log logs given message into file with today's date
func Log(message string) {
	file := getFile()
	if file == nil {
		return
	}
	defer func() {
		file.Close()
	}()
	w := bufio.NewWriter(file)
	w.WriteString(time.Now().Format("3:4:5-- ") + message)
	w.WriteString("\r\n")
	w.Flush()
}

func getFile() *os.File {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModeDir)
	}

	currentTime := time.Now()
	dateFileName := currentTime.Format("logs/02_01_2006") + ".txt"

	if _, err := os.Stat(dateFileName); err == nil {
		// path/to/whatever exists
		fi, err := os.OpenFile(dateFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("initiate log: error opening file")
			fmt.Println(err)
		}
		return fi
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		fi, err := os.Create(dateFileName)
		if err != nil {
			fmt.Println("initiate log: error creating file")
			fmt.Println(err)
		}
		return fi
	} else {
		//Schrödinger case: file could or could not exist
		fmt.Println("initiate log: mistery error")
		fmt.Println(err)
		return nil
	}
}
