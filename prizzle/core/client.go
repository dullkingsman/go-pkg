package prizzle

import (
	"context"
	"database/sql"
	"github.com/dullkingsman/go-pkg/utils"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var DatabaseConnections []*DB

// LoadDatabaseConnections loads database connections.
//
// ___
//
// NOTE: The function takes in two optional arguments.
//
//	The first argument is the environment variable key for the master database url
//	The second argument is the environment variable key for the secondary database urls
func LoadDatabaseConnections(driver string, urlEnvKeys ...string) {
	var masterDbUrlKey = "DB_URL"
	var secondaryDbUrlsKey = "SECONDARY_DB_URLS"

	if len(urlEnvKeys) > 0 {
		masterDbUrlKey = urlEnvKeys[0]
	}

	if len(urlEnvKeys) > 1 {
		secondaryDbUrlsKey = urlEnvKeys[1]
	}

	var masterDbUrl = os.Getenv(masterDbUrlKey)
	var _secondaryDbUrls = os.Getenv(secondaryDbUrlsKey)

	if masterDbUrl == "" {
		utils.LogFatal("prizzle-database-connections-loader", "MASTER_DB_URL not set in environment file")
	}

	var secondaryDbUrls []string

	if _secondaryDbUrls != "" {
		secondaryDbUrls = strings.Split(_secondaryDbUrls, ",")

		if len(secondaryDbUrls) == 0 {
			utils.LogWarn("prizzle-database-connections-loader", "SECONDARY_DB_URLS not set in environment file")
		}
	}

	var dbUrls = []string{masterDbUrl}

	dbUrls = append(dbUrls, secondaryDbUrls...)

	var (
		urlsString        = "urls"
		connectionsString = "connections"
	)

	if len(dbUrls) == 1 {
		urlsString = "url"
		connectionsString = "connection"
	}

	utils.LogInfo("prizzle-database-connection-loader", "found "+utils.YellowString(strconv.Itoa(len(dbUrls)))+" database "+urlsString)

	utils.LogInfo("prizzle-database-connection-loader", "creating database "+connectionsString+"...")

	for _, connString := range dbUrls {
		var db, err = sql.Open(driver, connString)

		if err != nil {
			utils.LogFatal("prizzle-database-connection-loader", "database connection error: "+err.Error())
			return
		}

		db.SetMaxOpenConns(10)

		db.SetMaxIdleConns(5)

		if err := db.Ping(); err != nil {
			utils.LogFatal("prizzle-database-connection-loader", "database ping error: "+err.Error())
		}

		var dbInfo = ""

		if info, parseError := url.Parse(connString); parseError == nil {
			dbInfo = " " + utils.GreyString(strings.TrimSuffix(info.Path, "/")) + " at " + utils.GreyString(info.Host)
		}

		utils.LogSuccess("prizzle-database-connection-loader", "successfully connected to database%s", dbInfo)

		DatabaseConnections = append(DatabaseConnections, &DB{DB: db})
	}
}

func GetRandomReadConnection() *DB {
	if len(DatabaseConnections) == 1 {
		return GetMasterConnection()
	}

	if len(DatabaseConnections) == 2 {
		return DatabaseConnections[1]
	}

	index := rand.Intn(len(DatabaseConnections))

	if index == 0 {
		index++
	}

	return DatabaseConnections[index]
}

func GetTransactor(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	var tx, err = DatabaseConnections[0].BeginTx(ctx, opts)

	if err != nil {
		utils.LogError("prizzle-database-transactor-creator", "could not begin transaction: "+err.Error())
		return nil, err
	}

	return &Tx{Tx: tx}, nil
}

func GetMasterConnection() *DB {
	return DatabaseConnections[0]
}

func CloseDbConnections() {
	var connectionsString = "connections"

	if len(DatabaseConnections) == 1 {
		connectionsString = "connection"
	}

	utils.LogInfo("prizzle-database-connection-cleaner", "closing database "+connectionsString+"...")

	for _, conn := range DatabaseConnections {

		err := conn.Close()

		if err != nil {
			utils.LogError("prizzle-database-connection-cleaner", "could not close db connection: "+err.Error())
		}
	}

	utils.LogInfo("prizzle-database-connection-cleaner", "finished closing database "+connectionsString)
}
