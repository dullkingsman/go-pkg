package roga

import (
	"encoding/json"
	"fmt"
	"github.com/dullkingsman/go-pkg/roga/writable"
	"github.com/dullkingsman/go-pkg/utils"
	"strconv"
	"time"
)

type DefaultStdoutFormatter struct{}

func (d DefaultStdoutFormatter) Format(item writable.Writable) string {
	var disableColor = true

	if item == nil {
		return ""
	}

	var log = utils.SafeCastValue[Log](item)

	if log != nil {
		var logType = "log"

		if log.Type == TypeAudit {
			logType = "audit"
		} else if log.Type == TypeEvent {
			logType = "event"
		}

		return fmt.Sprintf(utils.CyanString(logType, disableColor) + "(" + utils.GreyString(log.Id.String(), disableColor) + ") " + log.Message)
	}

	var operation = utils.SafeCastValue[Operation](item)

	if operation == nil {
		return ""
	}

	var tmp = utils.CyanString("op", disableColor) + "(" + utils.GreyString(operation.Id.String(), disableColor) + ") "

	var additional = ""

	if operation.Measurements != nil {
		var prefix = "\n                    " + " " + utils.GetEquivalentWhiteSpace(operation.Name)

		var jsonBytes, err = json.MarshalIndent(operation.Measurements, prefix, "    ")

		if err != nil {
			utils.LogError("roga:operation-writer(stdout)", "could not marshal additional values")
		} else {
			additional += prefix + utils.BeautifyIndentedJson(jsonBytes, disableColor)
		}
	}

	if operation.EssentialMeasurements.EndTime.IsZero() {
		return tmp + "started at " + utils.GreyString(operation.EssentialMeasurements.StartTime.Format(time.RFC3339), disableColor) + additional
	}

	var took = operation.EssentialMeasurements.EndTime.Sub(operation.EssentialMeasurements.StartTime).Milliseconds()

	return tmp + "finished in " + utils.GreyString(strconv.FormatInt(took, 10)+"ms"+additional, disableColor)
}

type DefaultFileFormatter struct{}

type DefaultExternalFormatter DefaultStdoutFormatter

func (d DefaultFileFormatter) Format(item writable.Writable) string {
	if item == nil {
		return ""
	}

	bytes, _ := item.Json()

	return string(bytes)
}
