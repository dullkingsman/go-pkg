package roga

import (
	"encoding/json"
	"github.com/dullkingsman/go-pkg/utils"
	"os"
	"strconv"
)

type DefaultWriter struct{ Writer }

func (d DefaultWriter) WriteOperationsToStdout(items []Operation, r *Roga) {
	for _, item := range items {
		var additional = ""

		if item.Measurements != nil {
			var prefix = "\n                    " + " " + utils.GetEquivalentWhiteSpace(item.Name)

			var jsonBytes, err = json.MarshalIndent(item.Measurements, prefix, "    ")

			if err != nil {
				utils.LogError("roga:operation-writer(stdout)", "could not marshal additional values")
			} else {
				additional += prefix + utils.BeautifyIndentedJson(jsonBytes)
			}
		}

		utils.LogInfo(item.Name, "finished in "+utils.GreyString(
			strconv.FormatInt(item.EssentialMeasurements.EndTime.Sub(
				item.EssentialMeasurements.StartTime,
			).Milliseconds(), 64)+"ms"+additional,
		))
	}
}

func (d DefaultWriter) WriteOperationsToFile(items []Operation, file *os.File, r *Roga) {
	var jsonItems []byte

	for _, item := range items {
		var jsonBytes, err = json.Marshal(item)

		if err != nil {
			utils.LogError("roga:operation-writer(file)", "could not marshal operation "+utils.GreyString(item.Id.String()))
			continue
		}

		jsonItems = append(jsonItems, '\n')
		jsonItems = append(jsonItems, jsonBytes...)
	}

	if len(jsonItems) == 0 {
		return
	}

	var _, err = file.Write(jsonItems)

	if err != nil {
		utils.LogError("roga:operation-writer(file)", "could not write operations to file: "+err.Error())
	}

}

func (d DefaultWriter) WriteOperationsToExternal(items []Operation, r *Roga) {
	// TODO implementation
}

func (d DefaultWriter) WriteLogsToStdout(items []Log, r *Roga) {
	for _, item := range items {
		if item.Type != 0 {
			continue
		}

		var operation = r.buffers.operations[item.OperationId]

		switch item.Level {
		case LevelFatal:
			utils.LogError(operation.Name, item.String())
			os.Exit(1)
		case LevelError:
			utils.LogError(operation.Name, item.String())
		case LevelWarn:
			utils.LogWarn(operation.Name, item.String())
		case LevelInfo:
			utils.LogInfo(operation.Name, item.String())
		case LevelDebug:
			utils.LogDebug(operation.Name, item.String())
		default:
			utils.LogInfo(operation.Name, item.String())
		}
	}
}

func (d DefaultWriter) WriteLogsToFile(items []Log, normal *os.File, audit *os.File, event *os.File, r *Roga) {
	var jsonNormalItems []byte
	var jsonAuditItems []byte
	var jsonEventItems []byte

	for _, item := range items {
		var jsonBytes, err = json.Marshal(item)

		if err != nil {
			utils.LogError("roga:log-writer(file)", "could not marshal log "+utils.GreyString(item.Id.String()))
			continue
		}

		switch item.Type {
		case TypeNormal:
			jsonNormalItems = append(jsonNormalItems, '\n')
			jsonNormalItems = append(jsonNormalItems, jsonBytes...)
		case TypeAudit:
			jsonAuditItems = append(jsonAuditItems, '\n')
			jsonAuditItems = append(jsonAuditItems, jsonBytes...)
		case TypeEvent:
			jsonEventItems = append(jsonEventItems, '\n')
			jsonEventItems = append(jsonEventItems, jsonBytes...)
		}
	}

	if len(jsonNormalItems) == 0 && len(jsonAuditItems) == 0 && len(jsonEventItems) == 0 {
		return
	}

	if len(jsonNormalItems) > 0 {
		var _, err = normal.Write(jsonNormalItems)

		if err != nil {
			utils.LogError("roga:log-writer(file)", "could not write normal logs to file: "+err.Error())
		}
	}

	if len(jsonAuditItems) > 0 {
		var _, err = audit.Write(jsonAuditItems)

		if err != nil {
			utils.LogError("roga:log-writer(file)", "could not write audit logs to file: "+err.Error())
		}
	}

	if len(jsonEventItems) > 0 {
		var _, err = event.Write(jsonEventItems)

		if err != nil {
			utils.LogError("roga:log-writer(file)", "could not write event logs to file: "+err.Error())
		}
	}
}

func (d DefaultWriter) WriteLogsToExternal(items []Log, r *Roga) {
	// TODO implementation
}
