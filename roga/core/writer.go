package roga

import (
	"bufio"
	"encoding/json"
	"github.com/dullkingsman/go-pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

type DefaultWriter struct{ Writer }

func (d DefaultWriter) WriteOperationsToStdout(items []Operation, r *Roga) {
	for _, item := range items {
		utils.LogInfo(item.Name+"("+item.Id.String()+")", item.String(r))
	}
}

func (d DefaultWriter) WriteOperationsToFile(items []Operation, file *os.File, r *Roga) {

	var instanceConfig = defaultRogaConfig.Instance.Inner()
	if instanceConfig.logsFormat == ".bson" {
		operationWriter := bufio.NewWriter(file)

		defer operationWriter.Flush()

		var itemsWritten int
		for _, item := range items {
			bsonByte, err := bson.Marshal(item)
			if err != nil {
				utils.LogError("roga:operation-writer(file)", "could not marshal operation "+utils.GreyString(item.Id.String()))
				continue
			}

			_, writeErr := operationWriter.Write(bsonByte)
			if writeErr != nil {
				utils.LogError("roga:log-writer(file)", "could not write log to buffer: "+writeErr.Error())
				continue
			}
			itemsWritten++
		}
	} else {
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
}

func (d DefaultWriter) WriteOperationsToExternal(items []Operation, r *Roga) {
	// TODO implementation
}

func (d DefaultWriter) WriteLogsToStdout(items []Log, r *Roga) {
	for _, item := range items {
		var operation = r.buffers.operations.Read(item.OperationId.String())

		var operationName = "root"

		if operation != nil {
			operationName = operation.Name
		}

		switch item.Level {
		case LevelFatal:
			utils.LogError(operationName, item.String(r))
			os.Exit(1)
		case LevelError:
			utils.LogError(operationName, item.String(r))
		case LevelWarn:
			utils.LogWarn(operationName, item.String(r))
		case LevelInfo:
			utils.LogInfo(operationName, item.String(r))
		case LevelDebug:
			utils.LogDebug(operationName, item.String(r))
		default:
			utils.LogInfo(operationName, item.String(r))
		}
	}
}

func (d DefaultWriter) WriteLogsToFile(items []Log, normal *os.File, audit *os.File, event *os.File, r *Roga) {

	var instanceConfig = defaultRogaConfig.Instance.Inner()
	if instanceConfig.logsFormat == ".bson" {

		normalWriter := bufio.NewWriter(normal)
		auditWriter := bufio.NewWriter(audit)
		eventWriter := bufio.NewWriter(event)

		// Ensure all buffered data is written to the files before the function exits.
		defer normalWriter.Flush()
		defer auditWriter.Flush()
		defer eventWriter.Flush()

		var itemsWritten int
		for _, item := range items {

			bsonBytes, err := bson.Marshal(item)
			if err != nil {
				utils.LogError("roga:log-writer(file)", "could not marshal log "+utils.GreyString(item.Id.String()))
				continue
			}
			var writeErr error
			switch item.Type {
			case TypeNormal:
				_, writeErr = normalWriter.Write(bsonBytes)
			case TypeAudit:
				_, writeErr = auditWriter.Write(bsonBytes)
			case TypeEvent:
				_, writeErr = eventWriter.Write(bsonBytes)
			}
			if writeErr != nil {
				utils.LogError("roga:log-writer(file)", "could not write log to buffer: "+writeErr.Error())
				continue
			}
			itemsWritten++

		}

	} else {
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
}

func (d DefaultWriter) WriteLogsToExternal(items []Log, r *Roga) {
	// TODO implementation
}
