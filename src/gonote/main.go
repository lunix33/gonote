package main

import (
	"fmt"
	"gonote/db"
	"gonote/models/setting"
	"gonote/route"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to DB
	dbID, err := db.Connect()
	if err != nil {
		panic(err)
	}

	sets := setting.GetAll(&dbID)

	dbSetup(sets[setting.DBVersion], &dbID)
	listen := formatInterfacePort(
		sets[setting.Interface], sets[setting.Port], &dbID)

	db.Close(dbID)

	// Register the web routes.
	route.RegisterRoute()

	// Start the web server on designated interface and port.
	log.Printf("Listening on: %s\n", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}

// dbSetup ensure the database is properly initialized and up to date
// `dbID` is the ID of the database.
func dbSetup(dbVerSet *setting.Setting, dbID *string) {
	// Validate the db version setting.
	if dbVerSet == nil {
		panic("unable to validate the application database version")
	}

	// Convert the version to workable number.
	dbVersion, convErr := strconv.ParseInt(dbVerSet.Value, 10, 0)
	if convErr != nil {
		panic("unable to parse the database version")
	}

	// Apply migration.
	db.MigrateFrom(dbVersion, 0, dbID)
}

// formatInterfacePort allow to format (or default) the interface and port of the application.
// `dbID` is the ID of the database.
// Returns a string with the interface and port on which the application can run.
func formatInterfacePort(interfaceSet *setting.Setting, portSet *setting.Setting, dbID *string) string {
	// Validate the interface setting.
	if interfaceSet == nil || interfaceSet.Value == "" {
		interfaceSet = new(setting.Setting)
		interfaceSet.Key = setting.Interface
		interfaceSet.Set("localhost", dbID)
		log.Printf("Unable to get the application interface. Defaulting to: %s", interfaceSet.Value)
	}

	// Validate the port setting.
	if portSet == nil {
		portSet = new(setting.Setting)
		portSet.Key = setting.Port
		portSet.Set("8080", dbID)
		log.Printf("Unable to get the application port. Defaulting to: %s", portSet.Value)
	}

	return fmt.Sprintf("%s:%s", interfaceSet.Value, portSet.Value)
}
