package main

import (
	"fmt"
	"gonote/db"
	"gonote/mngment"
	"gonote/router"
	"gonote/util"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func main() {
	// Connect to DB
	var listen string
	db.MustConnect(nil, func(c *db.Conn) {
		sets := mngment.GetAllSettings(c)

		dbSetup(sets[mngment.DBVersionSetting], c)
		listen = formatInterfacePort(
			sets[mngment.InterfaceSetting], sets[mngment.PortSetting], c)
	})

	// Register the web routes.
	router.RegisterRoute()

	// Start token cleanup routine
	go mngment.TokenCleanupRoutine()

	// Start the web server on designated interface and port.
	log.Printf("Listening on: %s\n", listen)
	util.LogErr(http.ListenAndServe(listen, nil))
}

// dbSetup ensure the database is properly initialized and up to date
//
// "dbVerSet" the setting for the database version.
// "c" is the database connection.
func dbSetup(dbVerSet *mngment.Setting, c *db.Conn) {
	// Validate the db version setting.
	if dbVerSet == nil {
		panic(errors.New("unable to validate the application database version"))
	}

	// Convert the version to workable number.
	dbVersion, convErr := strconv.ParseInt(dbVerSet.Value, 10, 0)
	if convErr != nil {
		panic(errors.Wrap(convErr, "unable to parse the database version"))
	}

	// Apply migration.
	db.MigrateFrom(dbVersion, 0, c)
}

// formatInterfacePort allow to format (or default) the interface and port of the application.
//
// "interfaceSet" is the structure for the interface setting.
// "portSet" is the structure for the port setting.
// "c" is the database connection
//
// Returns a string with the interface and port on which the application can run.
func formatInterfacePort(interfaceSet *mngment.Setting, portSet *mngment.Setting, c *db.Conn) string {
	// Validate the interface setting.
	if interfaceSet == nil || interfaceSet.Value == "" {
		interfaceSet = &mngment.Setting{
			Key: mngment.InterfaceSetting}
		interfaceSet.Set("localhost", c)
		log.Printf("Unable to get the application interface. Defaulting to: %s", interfaceSet.Value)
	}

	// Validate the port setting.
	if portSet == nil {
		portSet = &mngment.Setting{
			Key: mngment.PortSetting}
		portSet.Set("8080", c)
		log.Printf("Unable to get the application port. Defaulting to: %s", portSet.Value)
	}

	return fmt.Sprintf("%s:%s", interfaceSet.Value, portSet.Value)
}
