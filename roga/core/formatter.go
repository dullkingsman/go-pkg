package roga

import (
	"encoding/json"
	"fmt"
	"github.com/dullkingsman/go-pkg/utils"
	"strconv"
	"time"
)

type DefaultLogFormatter struct{ Formatter[Log] }

func (d DefaultLogFormatter) Format(log Log) string {
	return fmt.Sprintf("Log{Id: %v, Message: %s}", log.Id, log.Message)
}

type DefaultOperationFormatter struct{ Formatter[Operation] }

func (o DefaultOperationFormatter) Format(operation Operation) string {
	var additional = ""

	if operation.Measurements != nil {
		var prefix = "\n                    " + " " + utils.GetEquivalentWhiteSpace(operation.Name)

		var jsonBytes, err = json.MarshalIndent(operation.Measurements, prefix, "    ")

		if err != nil {
			utils.LogError("roga:operation-writer(stdout)", "could not marshal additional values")
		} else {
			additional += prefix + utils.BeautifyIndentedJson(jsonBytes)
		}
	}

	if operation.EssentialMeasurements.EndTime.IsZero() {
		return "started at " + utils.GreyString(operation.EssentialMeasurements.StartTime.Format(time.RFC3339)) + additional
	}

	var took = operation.EssentialMeasurements.EndTime.Sub(operation.EssentialMeasurements.StartTime).Milliseconds()

	return "finished in " + utils.GreyString(strconv.FormatInt(took, 10)+"ms"+additional)
}
