package main

import (
	"log"
	"os"
	"os/signal"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"haram_bot/cmd"
	"haram_bot/db"
	"haram_bot/parse_env"
	"haram_bot/tools"
)

var (
	ServerID       string
	BotToken       string
	RemoveCommands bool
	AppID          string
	LogLevel       string
	LogDestination string
)

var s *discordgo.Session

func init() {
	var dotEnvFilePath string = parse_env.GetEnvString("ENV_FILE_PATH", ".env")
	godotenv.Load(dotEnvFilePath)
	BotToken = parse_env.GetEnvString("BOT_TOKEN")
	ServerID = parse_env.GetEnvString("SERVER_ID", "")
	AppID = parse_env.GetEnvString("APP_ID", "")
	RemoveCommands = parse_env.GetEnvBool("RM_CMD_ON_SHUTDOWN", true)
	LogLevel = parse_env.GetEnvString("LOG_LEVEL", "info")
	LogDestination = parse_env.GetEnvString("LOG_DESTINATION", "console")

	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	if err := tools.InitLogger(LogLevel, LogDestination); err != nil {
		log.Fatalf("ERROR: logger initialization error: %v", err)
	}
	tools.LogTrace("Logger active")

	// Test logger
	// for i := 0; i < 3010; i++ {
	// 	tools.LogInfo("Log i = %v", i)
	// 	time.Sleep(3 * time.Millisecond)
	// }

	database, err := db.Connect()
	if err != nil {
		tools.LogFatal("Failed to connect to database: %v", err)
	}
	defer database.Close()

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		tools.LogInfo("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	tools.RegisterHandlers(s,
		cmd.GetExampleHandlers(),
		cmd.GetModHandlers(),
		cmd.GetConfHandlers(),
		cmd.GetHelpHandlers(),
	)
	err = s.Open()
	if err != nil {
		tools.LogFatal("Cannot open the session: %v", err)
	}
	defer s.Close()

	db.InitServersTable(s, database)

	tools.LogInfo("Adding commands...")
	commands := slices.Concat(
		cmd.GetExampleCommands(),
		cmd.GetModCommands(),
		cmd.GetConfCommands(),
		cmd.GetHelpCommands(),
	)
	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, ServerID, commands)

	if err != nil {
		tools.LogFatal("Cannot register commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	tools.LogInfo("Press Ctrl+C to exit")
	// cmd.HelloModCommand() //example. TODO: delete
	cmd.HelloConfCommand()
	cmd.HelloUserCommand()
	<-stop

	if RemoveCommands {
		tools.LogInfo("Removing commands...")
		for _, cmd := range createdCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, ServerID, cmd.ID)
			if err != nil {
				tools.LogError("Cannot delete %q command: %v", cmd.Name, err)
			}
		}
	}
	tools.LogInfo("Gracefully shutting down")
}
