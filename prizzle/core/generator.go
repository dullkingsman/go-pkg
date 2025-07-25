package prizzle

import (
	"database/sql"
	"encoding/json"
	"github.com/dullkingsman/go-pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

func GenerateClientModel(driver string, client *DB, schemaFilePath string) {
	var (
		enums map[string]Enum
		err   error
	)

	if driver == "postgres" {
		enums, err = GetEnumsInfo(client)

		if err != nil {
			utils.LogFatal("prizzle-generator", "could not get enums info: "+err.Error())
		}
	}

	tables, err := GetTablesInfo(driver, client)

	if err != nil {
		utils.LogFatal("prizzle-generator", "could not get tables info: "+err.Error())
	}

	if len(tables) == 0 {
		utils.LogFatal("prizzle-generator", "no tables found")
	}

	if driver == "postgres" && len(enums) == 0 {
		utils.LogWarn("prizzle-generator", "no enums found")
	}

	GenerateDefinitionModel(driver, enums, tables, schemaFilePath)

	GenerateQueryModel(tables, schemaFilePath)
}

func GenerateDefinitionModel(
	driver string,
	enums map[string]Enum,
	tables map[string]Table,
	schemaFilePath string,
) {
	utils.LogInfo("prizzle-definition-model-generator", "generating definition model...")

	var buffer = "package client\n\n"

	buffer += "import prizzle \"github.com/dullkingsman/go-pkg/prizzle/core\"\n\n"

	var _enums = ""
	var _tables = ""

	_enums += "// ------------------------------------------------------------------------------------------------------------------\n"
	_enums += "// ENUMS ------------------------------------------------------------------------------------------------------------\n"
	_enums += "// ------------------------------------------------------------------------------------------------------------------\n\n"

	for name, enum := range enums {
		var def = ""
		var init = ""
		var values = ""

		var goName = utils.AnyToPascalCase(name)

		for _, value := range enum.Values {
			var valueName = utils.AnyToPascalCase(strings.ToLower(value))

			def += "\t" + valueName + " " + goName + "\n"
			init += "\t" + valueName + ": " + goName + "(\"" + value + "\"),\n"
			values += "\t" + goName + "(\"" + value + "\"),\n"
		}

		_enums += "// " + goName + " --------------------------------------------------------------------------------------------\n\n"

		_enums += "type " + goName + " string\n\n"

		_enums += "type _" + goName + " struct {\n"
		_enums += def
		_enums += "}\n"

		_enums += "var " + goName + "value = _" + goName + " {\n"
		_enums += init
		_enums += "}\n"

		_enums += "var " + goName + "Values = []" + goName + " {\n"
		_enums += values
		_enums += "}\n"

		_enums += "func (e " + goName + ") String() string { return string(e) }\n\n\n"
	}

	_tables += "// ------------------------------------------------------------------------------------------------------------------\n"
	_tables += "// Tables -----------------------------------------------------------------------------------------------------------\n"
	_tables += "// ------------------------------------------------------------------------------------------------------------------\n\n"

	for name, table := range tables {
		var tableName = utils.AnyToPascalCase(name)

		_tables += "// " + tableName + " ------------------------------------------------------------------------------------------------\n\n"

		_tables += "type Inner" + tableName + " struct {\n"

		for _, column := range table.Columns {
			var col = utils.AnyToPascalCase(column.Name)

			var _type = PgTypeToGoType(column)

			if driver == "sqlite3" {
				_type = SqliteTypeToGoType(column)
			}

			_tables += "\t" + col + " " + _type + " " + "`json:\"" + utils.LowercaseFirstLetter(col) + ",omitempty\"`" + "\n"
		}

		_tables += "}\n\n\n"
	}

	buffer += _enums
	buffer += "\n"
	buffer += _tables

	formatted, err := utils.FormatAsGoCode(buffer)

	if err != nil {
		utils.LogFatal("prizzle-definition-model-generator", "could not format code: "+err.Error())
	}

	var filePath = filepath.Join(filepath.Dir(schemaFilePath) + "/client/definition.go")

	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		utils.LogFatal("prizzle-definition-model-generator", "could not create directories: "+err.Error())
	}

	if err := utils.WriteToFile(filePath, formatted); err != nil {
		utils.LogFatal("prizzle-definition-model-generator", "could not write to file: "+err.Error())
	}

	utils.LogSuccess("prizzle-definition-model-generator", "generated definition model")
}

func GenerateQueryModel(schema map[string]Table, schemaFilePath string) {
	utils.LogInfo("prizzle-query-model-generator", "generating query model...")

	var buffer = "package client\n\n"

	buffer += "import prizzle \"github.com/dullkingsman/go-pkg/prizzle/core\"\n\n"

	var types = ""
	var extensions = ""
	var values = ""

	for name, table := range schema {
		var tableName = utils.AnyToPascalCase(name)

		var typeColumns = ""
		var valueColumns = ""

		for _, column := range table.Columns {
			var col = utils.AnyToPascalCase(column.Name)

			typeColumns += "\t" + col + " prizzle.SqlName\n"

			valueColumns += "\t" + col + ": prizzle.SqlName(\"" + column.Name + "\"),\n"
		}

		types += "type _" + utils.AnyToLowerSnakeCase(name) + " struct {\n"
		types += "\tprizzle.SqlTable\n"
		types += typeColumns
		types += "}\n\n"

		extensions += "// " + utils.AnyToLowerSnakeCase(name) + " extensions\n"
		extensions += "func (t _" + utils.AnyToLowerSnakeCase(name) + ") GetSqlTable() prizzle.SqlTable { return t.SqlTable }\n"
		extensions += "func (t _" + utils.AnyToLowerSnakeCase(name) + ") As(alias string) prizzle.EmbedsSqlTable { t.SqlTable = t.SqlTable.As(alias); return t }\n"
		extensions += "func (t _" + utils.AnyToLowerSnakeCase(name) + ") Aliased(alias string) prizzle.EmbedsSqlTable { t.SqlTable = t.SqlTable.Aliased(alias); return t }\n"
		extensions += "func (t _" + utils.AnyToLowerSnakeCase(name) + ") Namespaced(alias string) prizzle.EmbedsSqlTable { t.SqlTable = t.SqlTable.Namespaced(alias); return t }\n\n"

		values += "var " + tableName + " = _" + utils.AnyToLowerSnakeCase(name) + "{\n\tSqlTable: prizzle.SqlTable{\n\t\tName: \"" + name + "\",\n\t},\n" + valueColumns + "}\n"
	}

	buffer += types

	buffer += "// --------------------------------------------------------------------------------------------------------------------\n"
	buffer += "// EXTENSIONS ---------------------------------------------------------------------------------------------------------\n"
	buffer += "// --------------------------------------------------------------------------------------------------------------------\n\n"
	buffer += extensions

	buffer += "// --------------------------------------------------------------------------------------------------------------------\n"
	buffer += "// QUERY TABLES -------------------------------------------------------------------------------------------------------\n"
	buffer += "// --------------------------------------------------------------------------------------------------------------------\n\n"
	buffer += values

	formatted, err := utils.FormatAsGoCode(buffer)

	if err != nil {
		utils.LogFatal("prizzle-query-model-generator", "could not format code: "+err.Error())
	}

	var filePath = filepath.Join(filepath.Dir(schemaFilePath) + "/client/query.go")

	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		utils.LogFatal("prizzle-query-model-generator", "could not create directories: "+err.Error())
	}

	if err := utils.WriteToFile(filePath, formatted); err != nil {
		utils.LogFatal("prizzle-query-model-generator", "could not write to file: "+err.Error())
	}

	utils.LogSuccess("prizzle-query-model-generator", "generated query model")
}

func GetTablesInfo(driver string, dbClient DatabaseClient, schemas ...string) (map[string]Table, error) {
	utils.LogInfo("prizzle-table-info-extractor", "getting tables info...")

	var tmp = "public"

	if len(schemas) > 0 {
		tmp = strings.Join(schemas, ", ")
	}

	var query = `WITH column_info AS (
    SELECT
    CASE
        WHEN c.table_schema = 'public' THEN c.table_name
        ELSE CONCAT(c.table_schema, '.', c.table_name)
        END AS table_name,
	c.table_schema as table_schema,
    c.column_name AS name,
    CASE
        WHEN data_type = 'USER-DEFINED' THEN
            udt_name
        ELSE
            data_type
        END AS type,
    c.is_nullable = 'YES' AS nullable,
    CASE
        WHEN c.data_type = 'USER-DEFINED' AND t.typcategory = 'E' THEN
            TRUE
        ELSE
            FALSE
        END AS is_enum,
	CASE
        WHEN c.data_type = 'USER-DEFINED' THEN
            n.nspname
        ELSE
            'public'
        END AS column_schema
    FROM
        information_schema.columns c
    LEFT JOIN
        pg_catalog.pg_type t
    ON
        c.udt_name = t.typname
	LEFT JOIN
		pg_catalog.pg_namespace n
    ON
        t.typnamespace = n.oid
	WHERE
    	(
        	c.table_schema = 'public' AND
        	c.table_name <> 'pg_stat_statements' AND
        	c.table_name <> '_prisma_migrations' AND
        	c.table_name <> 'pg_stat_statements_info' AND
        	c.table_name <> 'geometry_columns' AND
        	c.table_name <> 'geography_columns' AND
        	c.table_name <> 'spatial_ref_sys'
    	) OR
    	c.table_schema IN (` + tmp + `)
)
SELECT DISTINCT ON (table_name, table_schema)
    table_name as name,
    table_schema as schema,
    json_agg(json_build_object(
        'name', name,
        'type', type,
		'schema', column_schema,
        'nullable', nullable,
        'is_enum', is_enum
    )) AS columns
FROM
    column_info
GROUP BY
    table_name,
    table_schema
ORDER BY
    table_name;`

	if driver == "sqlite3" {
		query = `WITH column_info AS (
    SELECT
        m.name AS table_name,
        p.name AS column_name,
        p.type AS data_type,
        p.[notnull] = 0 AS nullable,
        FALSE AS is_enum -- SQLite does not support enums
    FROM
        sqlite_master m
            JOIN
        pragma_table_info(m.name) p
        ON
            m.type = 'table'
    WHERE
        m.name NOT LIKE 'sqlite_%' -- Exclude SQLite internal tables
      AND m.name NOT IN ('_prisma_migrations') -- Exclude specific tables
)
SELECT
    table_name AS name,
    json_group_array(
            json_object(
                    'name', column_name,
                    'type', data_type,
                    'nullable', nullable,
                    'is_enum', is_enum
            )
    ) AS columns
FROM
    column_info
GROUP BY
    table_name
ORDER BY
    table_name;`
	}

	rows, err := dbClient.Query(query)

	if err != nil {
		utils.LogError("prizzle-table-info-extractor", "could not query db: "+err.Error())
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.LogError("prizzle-table-info-extractor", "could not close rows: "+err.Error())
		}
	}(rows)

	var schema = make(map[string]Table)

	for rows.Next() {
		var table = Table{}
		var columns []uint8

		if err := rows.Scan(&table.Name, &columns); err != nil {
			utils.LogError("prizzle-table-info-extractor", "could not scan rows: "+err.Error())
			return nil, err
		}

		if driver == "sqlite3" {
			var _columns []SqliteColumn

			if err := json.Unmarshal(columns, &_columns); err != nil {
				utils.LogFatal("prizzle-table-info-extractor", "could not unmarshal columns: "+err.Error())
				return nil, err
			}

			for _, column := range _columns {
				table.Columns = append(table.Columns, Column{
					Name:     column.Name,
					Type:     column.Type,
					Nullable: column.Nullable == 1,
				})
			}
		} else {
			if err := json.Unmarshal(columns, &table.Columns); err != nil {
				utils.LogError("prizzle-table-info-extractor", "could not unmarshal columns: "+err.Error())
				return nil, err
			}
		}

		schema[table.Name] = table
	}

	if len(schema) == 0 {
		utils.LogError("prizzle-table-info-extractor", "no tables found")
		return schema, nil
	} else {
		utils.LogInfo("prizzle-table-info-extractor", "found tables:")

		for _, table := range schema {
			var columnList = ""
			var first = true

			for _, column := range table.Columns {
				if !first {
					columnList += ", "
				} else {
					first = false
				}

				columnList += utils.GreyString(column.Name)
			}

			utils.LogInfo("", utils.GreyString(table.Name)+" with columns: "+columnList)
		}
	}

	utils.LogSuccess("prizzle-table-info-extractor", "got tables info")

	return schema, nil
}

func GetEnumsInfo(dbClient DatabaseClient, schemas ...string) (map[string]Enum, error) {
	utils.LogInfo("prizzle-enum-info-extractor", "getting enums info...")

	var tmp = "public"

	if len(schemas) > 0 {
		tmp = strings.Join(schemas, ", ")
	}

	var query = `SELECT
    CASE
        WHEN n.nspname = 'public' THEN t.typname
        ELSE CONCAT(n.nspname, '_', t.typname)
        END AS name,
    json_agg(e.enumlabel) AS values
FROM
    pg_type t
        JOIN
    pg_enum e
    ON
        t.oid = e.enumtypid
        JOIN
    pg_catalog.pg_namespace n
    ON
        n.oid = t.typnamespace
WHERE
    n.nspname IN (` + tmp + `)
GROUP BY
    n.nspname, t.typname`

	rows, err := dbClient.Query(query)

	if err != nil {
		utils.LogError("prizzle-enum-info-extractor", "could not query db: "+err.Error())
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.LogError("prizzle-enum-info-extractor", "could not close rows: "+err.Error())
		}
	}(rows)

	var schema = make(map[string]Enum)

	for rows.Next() {
		var enum = Enum{}
		var values []byte

		if err := rows.Scan(&enum.Name, &values); err != nil {
			utils.LogError("prizzle-enum-info-extractor", "could not scan rows: "+err.Error())
			return nil, err
		}

		if err := json.Unmarshal(values, &enum.Values); err != nil {
			utils.LogError("prizzle-enum-info-extractor", "could not unmarshal values: "+err.Error())
			return nil, err
		}

		schema[enum.Name] = enum
	}

	utils.LogSuccess("prizzle-enum-info-extractor", "got enums info")

	return schema, nil
}

func SqliteTypeToGoType(column Column) string {
	switch strings.ToUpper(column.Type) {
	case "TINYINT":
		if column.Nullable {
			return "*int8"
		}
		return "int8"
	case "SMALLINT", "INT2":
		if column.Nullable {
			return "*int16"
		}
		return "int16"
	case "MEDIUMINT":
		if column.Nullable {
			return "*int32"
		}
		return "int32"
	case "INTEGER", "INT", "BIGINT", "INT8":
		if column.Nullable {
			return "*int64"
		}
		return "int64"

	case "UNSIGNED BIG INT":
		if column.Nullable {
			return "*uint64"
		}
		return "uint64"

	case "REAL", "DOUBLE", "DOUBLE PRECISION", "FLOAT", "NUMERIC", "DECIMAL":
		if column.Nullable {
			return "*float64"
		}
		return "float64"

	case "TEXT", "CLOB", "CHARACTER", "CHAR", "NCHAR", "NATIVE CHARACTER", "VARCHAR", "VARYING CHARACTER", "NVARCHAR":
		if column.Nullable {
			return "*string"
		}
		return "string"

	case "DATE", "DATETIME":
		if column.Nullable {
			return "*prizzle.DateTime"
		}
		return "prizzle.DateTime"

	case "BLOB":
		if column.Nullable {
			return "*[]byte"
		}
		return "[]byte"

	case "BOOLEAN":
		if column.Nullable {
			return "*bool"
		}
		return "bool"

	default:
		if column.Nullable {
			return "*string"
		}
		return "string"
	}
}

func PgTypeToGoType(column Column, schema ...string) string {
	switch strings.ToLower(column.Type) {
	case "bigint":
		if column.Nullable {
			return "*int64"
		}
		return "int64"
	case "integer":
		if column.Nullable {
			return "*int"
		}
		return "int"
	case "smallint":
		if column.Nullable {
			return "*int16"
		}
		return "int16"
	case "boolean":
		if column.Nullable {
			return "*bool"
		}
		return "bool"
	case "text", "varchar", "character varying":
		if column.Nullable {
			return "*string"
		}
		return "string"
	case "timestamp", "timestamp without time zone", "timestamp with time zone":
		if column.Nullable {
			return "*prizzle.DateTime"
		}
		return "prizzle.DateTime"
	case "numeric", "decimal", "double precision":
		if column.Nullable {
			return "*float64"
		}
		return "float64"
	}

	if column.IsEnum {
		var _type = column.Type

		if len(schema) > 0 {
			if schema[0] != "public" {
				_type = utils.CapitalizeFirstLetter(schema[0]) + utils.CapitalizeFirstLetter(_type)
			}
		}

		if column.Nullable {
			return "*" + _type
		}

		return _type
	}

	if column.Nullable {
		return "*string"
	}

	return "string"
}

type Enum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type Table struct {
	Name    string   `json:"name"`
	Schema  string   `json:"schema"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Schema   string `json:"schema"`
	Nullable bool   `json:"nullable"`
	IsEnum   bool   `json:"is_enum"`
}

type SqliteColumn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable int    `json:"nullable"`
}
