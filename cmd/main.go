package main

import (
	"context"
	"fmt"
	"gazprom/pkg/commands"
	"gazprom/logging"
	"gazprom/pkg/utils"
	"os"
	"os/signal"
	"strings"
	"time"

	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	_ "github.com/lib/pq"
)

type Config struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string

}

func main() {
	logger := logging.GetLogger()

	if err := initConfig(); err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("error loading env variables: %s", err.Error())
	}

	cfg := Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
	}



	utils.DBConn, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
	cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		logger.Fatalf("failed to initialize database: %s", err.Error())
	}
	defer utils.DBConn.Close()

	err = utils.DBConn.Ping()
	if err != nil {
		logger.Fatalf("Connection to the database is not alive: %s", err.Error())
	}
	utils.DBConn.SetConnMaxLifetime(time.Minute * 3)
	utils.DBConn.SetMaxOpenConns(10)
	utils.DBConn.SetMaxIdleConns(10)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		// bot.WithDebug(),
	}

	
	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		logger.Fatalf("failed to create a bot instance: %s", err.Error())
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypePrefix, commands.AddHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/showall", bot.MatchTypeExact, commands.ShowallHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/search", bot.MatchTypeExact, commands.SearchHandler)
	// b.RegisterHandler(bot.HandlerTypeMessageText, "/edit", bot.MatchTypeExact, commands.SearchHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/remove", bot.MatchTypeExact, commands.SearchHandler)

	utils.StartNotifier(ctx, b)
	b.Start(ctx)
}

func sendHelpCommand(ctx context.Context, b *bot.Bot, chatID int64) {
	msg := "Hi, I am a birthday reminder bot. I can help you remember birthdays of your friends and family.\n\n"
	msg += "You can use the following commands to interact with me:\n"
	msg += "  • /add <name> - Add a new birthday\n"
	msg += "  • /showall - Show all birthdays\n"
	msg += "  • /search - Search for birthdays\n"
	msg += "  • /remove - Remove a birthday\n"
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		fmt.Println("update.Message is nil")
		return
	}

	if update.Message.ReplyToMessage == nil {
		sendHelpCommand(ctx, b, update.Message.Chat.ID)
		return
	}

	cmd := strings.Split(update.Message.ReplyToMessage.Text, "/")
	if len(cmd) != 2 {
		sendHelpCommand(ctx, b, update.Message.Chat.ID)
		return
	}
	if !slices.Contains([]string{"search", "edit", "remove"}, cmd[1]) {
		sendHelpCommand(ctx, b, update.Message.Chat.ID)
		return
	}

	update.Message.ReplyToMessage.Text = "/" + cmd[1]
	bd := utils.BirthdayInfo{
		ChatID: update.Message.Chat.ID,
		Name:   update.Message.Text,
	}
	commands.Search(ctx, b, update.Message, &bd, "name")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}