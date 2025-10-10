package file

import (
	"errors"
	"fmt"
	"os"
)

/* TODO
getCurrentTimeRoundedTo(
			config.fileLogsDirectoryGranularity,
		).UTC().Format(
			config.fileLogsDirectoryFormatLayout,
		)


	switch entryType {
	case config.EntryTypeOperation:
		filePath += config.operationsFileName
	case config.EntryTypeAudit:
		filePath += "audit." + config.logsFileName
	case config.EntryTypeEvent:
		filePath += "event." + config.logsFileName
	case config.EntryTypeLog:
		filePath += config.logsFileName
	default:
		return nil, errors.New("invalid entry type")
	}
*/

// GetLogFileDescriptor returns a file descriptor for the given file name in the log base directory
func GetLogFileDescriptor(
	baseFilePath string,
	dirname string,
	fileName string,
) (*os.File, error) {
	var logsBasePath = baseFilePath + dirname

	var err = os.MkdirAll(logsBasePath, os.ModePerm)

	if err != nil {
		fmt.Println("could not make base directory:", err.Error())
		return nil, err
	}

	var filePath = logsBasePath + "/" + fileName

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil || file == nil {
		if err == nil {
			err = errors.New(filePath)
		}

		fmt.Println("could not open file:", err.Error())

		return nil, err
	}

	return file, nil
}
