package prizzle

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/dullkingsman/go-pkg/utils"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

func LoadDatabaseCluster(driver string, config ClusterConfig) Cluster {
	var cluster = Cluster{
		WriteNodes: map[string]*DB{},
		ReadNodes:  map[string]*DB{},
	}

	var (
		urlsString        = "nodes"
		connectionsString = "connections"
	)

	if len(config.ReadNodesConfig) == 0 {
		urlsString = "node"
		connectionsString = "connection"
	}

	utils.LogInfo("prizzle-database-connection-loader", "found "+utils.YellowString(strconv.Itoa(1+len(config.ReadNodesConfig)))+" database "+urlsString)

	utils.LogInfo("prizzle-database-connection-loader", "creating database "+connectionsString+"...")

	for key, nodeConfig := range config.ReadNodesConfig {
		cluster.ReadNodes[key] = LoadDatabaseClusterNode(driver, nodeConfig)
	}

	return cluster
}

func LoadDatabaseClusterNode(driver string, config ClusterNodeConfig) *DB {
	var db, err = sql.Open(driver, config.Url)

	if err != nil {
		utils.LogFatal("prizzle-database-connection-loader", "database connection error: "+err.Error())
		return nil
	}

	if config.MaxOpenConnections != nil {
		db.SetMaxOpenConns(*config.MaxOpenConnections)
	}

	if config.MaxIdleConnections != nil {
		db.SetMaxIdleConns(*config.MaxIdleConnections)
	}

	if config.MaxLifetime != nil {
		db.SetConnMaxLifetime(*config.MaxLifetime)
	}

	if config.MaxIdleTime != nil {
		db.SetConnMaxIdleTime(*config.MaxIdleTime)
	}

	if err := db.Ping(); err != nil {
		utils.LogFatal("prizzle-database-connection-loader", "database ping error: "+err.Error())
	}

	var dbInfo = ""

	if driver == "sqlite3" {
		dbInfo = " " + utils.GreyString(config.Url)
	} else {
		if info, parseError := url.Parse(config.Url); parseError == nil {
			dbInfo = " " + utils.GreyString(strings.TrimSuffix(info.Path, "/")) + " at " + utils.GreyString(info.Host)
		}
	}

	utils.LogSuccess("prizzle-database-connection-loader", "successfully connected to database"+dbInfo)

	return &DB{db}
}

func CloseDatabaseClusterNode(connection *DB) {
	err := connection.Close()

	if err != nil {
		utils.LogError("prizzle-database-connection-cleaner", "could not close db connection: "+err.Error())
	}
}

func (c *Cluster) GetRandomWriteNodeKey() string {
	var index = 0

	if len(c.writeNodeKeys) > 1 {
		index = rand.Intn(len(c.writeNodeKeys))
	}

	var i = 0

	for _, key := range c.writeNodeKeys {
		if i == index {
			return key
		}

		i++
	}

	return c.writeNodeKeys[0]
}

func (c *Cluster) GetRandomReadNodeKey() string {
	if len(c.readNodeKeys) == 0 {
		return c.writeNodeKeys[0]
	}

	var index = 0

	if len(c.readNodeKeys) > 1 {
		index = rand.Intn(len(c.readNodeKeys))
	}

	var i = 0

	for _, key := range c.readNodeKeys {
		if i == index {
			return key
		}

		i++
	}

	return c.writeNodeKeys[0]
}

func (c *Cluster) GetReadConnection(key string) *DB {
	return c.ReadNodes[key]
}

func (c *Cluster) GetWriteConnection(key string) *DB {
	return c.WriteNodes[key]
}

func (c *Cluster) GetRandomReadConnection() *DB {
	return c.ReadNodes[c.GetRandomReadNodeKey()]
}

func (c *Cluster) GetRandomWriteConnection() *DB {
	return c.WriteNodes[c.GetRandomWriteNodeKey()]
}

func (c *Cluster) GetTransactor(key string) (*Tx, error) {
	var tx, err = c.GetWriteConnection(key).Begin()

	if err != nil {
		utils.LogError("prizzle-database-transactor-creator", "could not begin transaction: "+err.Error())
		return nil, err
	}

	return &Tx{Tx: tx}, nil
}

func (c *Cluster) GetTransactorContext(ctx context.Context, key string, opts *sql.TxOptions) (*Tx, error) {
	var tx, err = c.GetWriteConnection(key).BeginTx(ctx, opts)

	if err != nil {
		utils.LogError("prizzle-database-transactor-creator", "could not begin transaction: "+err.Error())
		return nil, err
	}

	return &Tx{Tx: tx}, nil
}

func (c *Cluster) GetRandomTransactor() (*Tx, error) {
	return c.GetTransactor(c.GetRandomWriteNodeKey())
}

func (c *Cluster) GetRandomTransactorContext(ctx context.Context) (*Tx, error) {
	return c.GetTransactorContext(ctx, c.GetRandomWriteNodeKey(), nil)
}

func (c *Cluster) CloseDbConnections() {
	var connectionsString = "connections"

	if len(c.ReadNodes) == 0 && len(c.WriteNodes) == 0 {
		connectionsString = "connection"
	}

	utils.LogInfo("prizzle-database-connection-cleaner", "closing database "+connectionsString+"...")

	for key, _ := range c.WriteNodes {
		c.CloseWriteNode(key)
	}

	for key, _ := range c.ReadNodes {
		c.CloseReadNode(key)
	}

	utils.LogInfo("prizzle-database-connection-cleaner", "finished closing database "+connectionsString)
}

func (c *Cluster) CloseReadNode(key string) {
	if len(c.ReadNodes) == 0 {
		utils.LogInfo("prizzle-database-connection-cleaner", "no read nodes found")
		return
	}

	var conn, ok = c.ReadNodes[key]

	if !ok {
		utils.LogError("prizzle-database-connection-cleaner", "could not find the db connection "+utils.GreyString(key))
		return
	}

	CloseDatabaseClusterNode(conn)
}

func (c *Cluster) CloseWriteNode(key string) {
	if len(c.WriteNodes) == 0 {
		utils.LogInfo("prizzle-database-connection-cleaner", "no write nodes found")
		return
	}

	var conn, ok = c.WriteNodes[key]

	if !ok {
		utils.LogError("prizzle-database-connection-cleaner", "could not find the db connection "+utils.GreyString(key))
		return
	}

	CloseDatabaseClusterNode(conn)
}

func (client *DB) Conn(ctx context.Context) (*Conn, error) {
	var conn, err = client.DB.Conn(ctx)

	if err != nil {
		utils.LogError("prizzle-database-connection-creator", "could not create connection: "+err.Error())
		return nil, err
	}

	return &Conn{Conn: conn}, nil
}

func (client *DB) Stats() sql.DBStats {
	return client.DB.Stats()
}

func (client *DB) Ping() error {
	return client.DB.Ping()
}

func (client *DB) PingContext(ctx context.Context) error {
	return client.DB.PingContext(ctx)
}

func (client *DB) Driver() driver.Driver {
	return client.DB.Driver()
}

func (c *Conn) PingContext(ctx context.Context) error {
	return c.Conn.PingContext(ctx)
}

// todo check if you can take to query instead
func (c *Conn) Row(f func(driverConn any) error) error {
	return c.Conn.Raw(f)
}

func (c *Conn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	var tx, err = c.Conn.BeginTx(ctx, opts)

	if err != nil {
		utils.LogError("prizzle-database-transaction-creator", "could not begin transaction: "+err.Error())
		return nil, err
	}

	return &Tx{Tx: tx}, nil
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}
