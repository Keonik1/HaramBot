package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"haram_bot/cmd"
	"haram_bot/parse_env"
)

// Bot parameters
var (
	ServerID       string
	BotToken       string
	RemoveCommands bool
	AppID          string
)

var s *discordgo.Session

func init() {
	var dotEnvFilePath string = parse_env.GetEnvString("ENV_FILE_PATH", ".env")
	godotenv.Load(dotEnvFilePath)
	BotToken = parse_env.GetEnvString("BOT_TOKEN")
	ServerID = parse_env.GetEnvString("SERVER_ID", "")
	AppID = parse_env.GetEnvString("APP_ID", "")
	RemoveCommands = parse_env.GetEnvBool("RM_CMD_ON_SHUTDOWN", true)

	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	var commandHandlers = cmd.GetExampleHandlers()

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	log.Println("Adding commands...")
	var commands = cmd.GetExampleCommands()
	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, ServerID, commands)

	if err != nil {
		log.Fatalf("Cannot register commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	cmd.HelloModCommand() //example. TODO: delete
	cmd.HelloConfCommand()
	cmd.HelloUserCommand()
	<-stop

	if RemoveCommands {
		log.Println("Removing commands...")
		for _, cmd := range createdCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, ServerID, cmd.ID)
			if err != nil {
				log.Fatalf("Cannot delete %q command: %v", cmd.Name, err)
			}
		}
	}
	log.Println("Gracefully shutting down")
}
