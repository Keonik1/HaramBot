package db

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
	_ "github.com/golang-migrate/migrate/v4/database"
	_ "github.com/mattn/go-sqlite3"

	"haram_bot/tools"
)

type Command struct {
	ActionType string            `json:"action_type"`
	Params     map[string]string `json:"params"`
}

type ServerConfig struct {
	NSFWChannels []string           `json:"nsfw_channels"`
	MuteRoles    []string           `json:"mute_roles"`
	Children     []string           `json:"children"`
	Birthdays    map[string]string  `json:"birthdays"` // userID -> date
	Commands     map[string]Command `json:"commands"`
}

var InitialServerConfig = ServerConfig{
	NSFWChannels: []string{},
	MuteRoles:    []string{},
	Children:     []string{},
	Birthdays:    map[string]string{},
	Commands:     map[string]Command{},
}

func Connect() (*sql.DB, error) {
	log.Println("Attempt connect to db.")
	return sql.Open("sqlite3", "./bot.db")
}

func InitServersTable(s *discordgo.Session, database *sql.DB) {
	var after string
	tools.LogInfo("Initializing bot DB.")
	query := "CREATE TABLE IF NOT EXISTS server_configs (server_id TEXT PRIMARY KEY, config TEXT NOT NULL);"
	if _, err := database.Exec(query); err != nil {
		log.Fatalf("ERROR: failed to create table: %v", err)
	}

	query = `
		INSERT INTO server_configs (server_id, config)
		VALUES (?, ?)
		ON CONFLICT(server_id) DO NOTHING
	`
	statement, err := database.Prepare(query)
	if err != nil {
		log.Fatalf("ERROR: server table creation error: %v", err)
	}
	defer statement.Close()

	for {
		servers, err := s.UserGuilds(100, "", after, false)
		if err != nil {
			log.Fatalf("ERROR: cant get servers: %v", err)
		}

		if len(servers) == 0 {
			log.Println("INFO: Bot not installed for no one server.")
			break
		}

		for _, server := range servers {
			log.Printf("INFO: Processing server:\n\tServer ID: %s\n\tServer name: %s\n", server.ID, server.Name)
			if err := InitServerConfig(statement, server.ID, InitialServerConfig); err != nil {
				log.Printf("Failed to insert config for server %s: %v", server.ID, err)
			}
		}
		after = servers[len(servers)-1].ID

		if len(servers) < 100 {
			break
		}
	}

	log.Println("INFO: Finish DB initializing.")
}

func saveServerConfig(stmt *sql.Stmt, serverID string, config ServerConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(serverID, string(data))
	return err
}

func InitServerConfig(stmt *sql.Stmt, serverID string, config ServerConfig) error {
	return saveServerConfig(stmt, serverID, config)
}

func UpdateServerConfig(db *sql.DB, serverID string, config ServerConfig) error {
	query := `
		INSERT INTO server_configs (server_id, config)
		VALUES (?, ?)
		ON CONFLICT(server_id) DO UPDATE SET config = excluded.config
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return saveServerConfig(stmt, serverID, config)
}

// func RunMigrations() error {
// 	m, err := migrate.New(
// 		"file://migrations",
// 		"sqlite3://bot.db",
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	err = m.Up()
// 	if err != nil && err != migrate.ErrNoChange {
// 		return err
// 	}

// 	return nil
// }
