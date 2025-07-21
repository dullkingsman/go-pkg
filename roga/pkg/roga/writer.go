package roga

import (
	"encoding/json"
	"errors"
	"github.com/dullkingsman/go-pkg/roga/internal/writer"
	"github.com/dullkingsman/go-pkg/roga/writable"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
)

type DefaultWriter struct{}

func (d DefaultWriter) WriteToStdout(items []writable.Writable, entryType EntryType, writer *writer.SyncedWriter, formatter writable.Formatter) (int, error) {
	for _, item := range items {
		if item == nil {
			continue
		}

		var bytes, _ = item.Json()

		//if len(bytes) > 0 {
		if _, err := writer.WriteString(string(bytes) + "\n"); err != nil {
			// TODO what to do
			continue
		}
		//}

		//switch entryType {
		//case EntryTypeOperation:
		//	operation := utils.SafeCastValue[Operation](item)
		//	entry := utils.CyanString("op") + "(" + utils.GreyString(operation.Id.String()) + ") " + operation.String(formatter)
		//	fmt.Printf(strings.TrimSpace(operation.Name + " " + entry))
		//
		//case EntryTypeAudit, EntryTypeEvent, EntryTypeLog:
		//	logItem := utils.SafeCastValue[Log](item)
		//	//var fmtFunc = utils.FormatInfoLog
		//
		//	//switch logItem.Level {
		//	//case LevelFatal, LevelError:
		//	//	fmtFunc = utils.FormatErrorLog
		//	//case LevelWarn:
		//	//	fmtFunc = utils.FormatWarnLog
		//	//case LevelInfo:
		//	//	fmtFunc = utils.FormatInfoLog
		//	//case LevelDebug:
		//	//	fmtFunc = utils.FormatDebugLog
		//	//default:
		//	//	fmtFunc = utils.FormatInfoLog
		//	//}
		//
		//	fmt.Printf(strings.TrimSpace(EntryTypeName[entryType] + " " + logItem.String(formatter)))
		//}

	}

	return writer.Writer.Buffered(), nil
}

func (d DefaultWriter) WriteToFile(items []writable.Writable, entryType EntryType, writer *writer.SyncedWriter) (int, error) {
	var jsonItems []byte

	for _, item := range items {
		if item == nil {
			continue
		}

		var itemId uuid.UUID

		if entryType == EntryTypeOperation {
			itemId = utils.SafeCastValue[Operation](item).Id
		} else {
			itemId = utils.SafeCastValue[Log](item).Id
		}

		var jsonBytes, err = json.Marshal(item)

		if err != nil {
			utils.LogError("roga:writer(file)", "could not marshal "+EntryTypeName[entryType]+" "+utils.GreyString(itemId.String()))
			// TODO what to do
			continue
		}

		jsonItems = append(jsonItems, '\n')
		jsonItems = append(jsonItems, jsonBytes...)
	}

	if len(jsonItems) == 0 {
		if len(jsonItems) != len(items) {
			return 0, errors.New("could not marshal all " + EntryTypeName[entryType] + "s to file: some of them could not be marshalled")
		}

		return 0, nil
	}

	var written, err = writer.Write(jsonItems)

	if err != nil {
		utils.LogError("roga:writer(file)", "could not write "+EntryTypeName[entryType]+"s logs to file: "+err.Error())
		return 0, err
	}

	return written, nil
}

type Writer interface {
	WriteToStdout(items []writable.Writable, entryType EntryType, writer *writer.SyncedWriter, formatter writable.Formatter) (int, error)
	WriteToFile(items []writable.Writable, entryType EntryType, writer *writer.SyncedWriter) (int, error)
}
