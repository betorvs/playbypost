package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/cli"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	Version  string = "development"
	commands        = []*discordgo.ApplicationCommand{
		{
			Name: "join",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Join to Play by Post",
		},
		{
			Name: "play-by-post",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "options",
					Description: "Options select menu",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "solo-start",
					Description: "Solo list select menu",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "solo-next",
					Description: "Solo next select menu",
				},
			},
			Description: "Main Play by Post command",
		},
	}
)

type app struct {
	logger   *slog.Logger
	web      *cli.Cli
	session  *discordgo.Session
	admToken string
	guildID  string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	logger.Info("starting discord bot test", "version", Version)
	token := os.Getenv("DISCORD_TOKEN")
	guildID := os.Getenv("DISCORD_GUILD_ID")
	playbypost := utils.GetEnv("PLAYBYPOST_SERVER", "http://localhost:3000")
	adminUser := utils.GetEnv("ADMIN_USER", "admin")
	adminFile := utils.GetEnv("CREDS_FILE", "./creds")
	adminToken, err := read(adminFile)
	if err != nil {
		logger.Error("error reading creds file", "error", err.Error())
	}
	logger.Info("debug", "token", adminToken)
	// new instances
	play := cli.NewHeaders(playbypost, adminUser, adminToken)
	mux := http.NewServeMux()
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("error creating discord session", "error", err)
		os.Exit(1)
	}
	// create internal app
	a := app{
		logger:   logger,
		web:      play,
		session:  discord,
		admToken: adminToken,
		guildID:  guildID,
	}
	// bot handlers
	discord.AddHandler(a.messageCreate)
	discord.AddHandler(a.interactionCommand)
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	// connect to discord
	err = discord.Open()
	if err != nil {
		logger.Error("error opening discord session", "error", err)
		os.Exit(1)
	}

	// register commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, guildID, v)
		if err != nil {
			logger.Error("error creating command", "error", err)
			os.Exit(2)
		}
		logger.Info("command created", "command", cmd.Name)
		registeredCommands[i] = cmd
	}

	// server config
	server := &http.Server{
		Addr:    ":8091",
		Handler: mux,
	}
	// web handlers
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"status\":\"OK\"}")
	})
	mux.HandleFunc("POST /api/v1/event", a.events)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen and serve error", "error", err)
			os.Exit(1)
		}
		logger.Info("stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	for _, v := range registeredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, guildID, v.ID)
		if err != nil {
			logger.Error("error deleting command", "error", err)
			os.Exit(2)
		}
	}
	logger.Info("commands deleted")
	logger.Info("shutting down bot...")
	discord.Close()

	ctxTimeout, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	if err := server.Shutdown(ctxTimeout); err != nil {
		logger.Error("shutdown error", "error", err)
	}
	logger.Info("graceful shutdown complete.")
}

func (a *app) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	a.logger.Info("message received", "message", m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}
	switch m.Content {
	case "hello":
		message := "Hello, I am Play by Post bot. How can I help you?"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println("Error sending message: ", err)
		}
	default:
		return
	}
}

func (a *app) interactionCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildid := i.GuildID
	if guildid != a.guildID {
		a.logger.Error("guild id mismatch", "guildid", guildid, "from_env", a.guildID)
	}
	var userid, username string
	if guildid != "" {
		userid = i.Member.User.ID
		username = i.Member.User.Username
	} else {
		userid = i.User.ID
		username = i.User.Username
	}
	a.logger.Info("user interaction", "user", userid, "type", i.Type)
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "join":
			// Do something
			if userid != "" {
				attachment := ""
				body, err := a.web.AddChatInformation(userid, username, i.ChannelID, types.Discord)
				if err != nil {
					a.logger.Error("error adding user info", "error", err.Error())
					attachment = fmt.Sprintf("Sorry, it did not work %s", username)
				}
				var msg types.Msg
				err = json.Unmarshal(body, &msg)
				if err != nil {
					a.logger.Error("error json unmarshal", "error", err.Error())
				}
				if strings.Contains(msg.Msg, "already added") {
					attachment = fmt.Sprintf("Already subscribed. Great, %s", username)
				}
				if strings.EqualFold(msg.Msg, "added") {
					attachment = fmt.Sprintf("Let's play %s", username)
				}
				response := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content:    attachment,
						Flags:      discordgo.MessageFlagsEphemeral,
						Components: nil,
					},
				}
				// should be via response
				err = s.InteractionRespond(i.Interaction, response)
				if err != nil {
					a.logger.Error("error responding to interaction", "error", err)
				}
			}

		case "play-by-post":
			var response *discordgo.InteractionResponse
			switch i.ApplicationCommandData().Options[0].Name {
			case "options", "opt":
				a.logger.Info("options")
				// post command
				msg, err := a.postCommand(userid, "opt", i.ChannelID)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
				}
				textPickItem := "Pick an item "
				if msg.Msg != "" {
					textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
				}
				if len(msg.Opt) > 0 {
					// create select menu
					options := []discordgo.SelectMenuOption{}
					for _, v := range msg.Opt {
						options = append(options, discordgo.SelectMenuOption{
							Label: v.Name,
							Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, i.ChannelID, userid, v.Name, v.ID),
						})
					}
					selectMenu := discordgo.SelectMenu{
						CustomID:    "opt",
						Placeholder: textPickItem,
						Options:     options,
					}
					// create action row
					actionRow := discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{selectMenu},
					}
					// send response back to discord
					response = &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content:    "Select an option",
							Flags:      discordgo.MessageFlagsEphemeral,
							Components: []discordgo.MessageComponent{actionRow},
						},
					}
				} else {
					response = &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content:    "No options available",
							Flags:      discordgo.MessageFlagsEphemeral,
							Components: nil,
						},
					}
				}

			case "solo-start":
				a.logger.Info("solo-start")
				// post command
				msg, err := a.postCommand(userid, "solo-start", i.ChannelID)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
				}
				textPickItem := "Pick an item "
				if msg.Msg != "" {
					textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
				}
				a.logger.Info("msg", "msg", msg)
				if len(msg.Opt) > 0 {
					// create select menu
					options := []discordgo.SelectMenuOption{}
					for _, v := range msg.Opt {
						options = append(options, discordgo.SelectMenuOption{
							Label: v.Name,
							Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, i.ChannelID, userid, v.Name, v.ID),
						})
					}
					selectMenu := discordgo.SelectMenu{
						CustomID:    "choice",
						Placeholder: textPickItem,
						Options:     options,
					}
					// create action row
					actionRow := discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{selectMenu},
					}
					// send response back to discord
					response = &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content:    "Select an option",
							Flags:      discordgo.MessageFlagsEphemeral,
							Components: []discordgo.MessageComponent{actionRow},
						},
					}
				} else {
					response = &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content:    "No options available",
							Flags:      discordgo.MessageFlagsEphemeral,
							Components: nil,
						},
					}
				}

			case "solo-next":
				a.logger.Info("solo-next")
				// post command
				msg, err := a.postCommand(userid, "solo-next", i.ChannelID)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
				}
				textPickItem := "Pick an item "
				if msg.Msg != "" {
					textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
				}
				if len(msg.Opt) > 0 {
					// create select menu
					options := []discordgo.SelectMenuOption{}
					for _, v := range msg.Opt {
						options = append(options, discordgo.SelectMenuOption{
							Label: v.Name,
							Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, i.ChannelID, userid, v.Name, v.ID),
						})
					}
					selectMenu := discordgo.SelectMenu{
						CustomID:    "choice",
						Placeholder: textPickItem,
						Options:     options,
					}
					// create action row
					actionRow := discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{selectMenu},
					}
					// send response back to discord
					response = &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content:    "Select an option",
							Flags:      discordgo.MessageFlagsEphemeral,
							Components: []discordgo.MessageComponent{actionRow},
						},
					}
				} else {
					response = &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content:    "No options available",
							Flags:      discordgo.MessageFlagsEphemeral,
							Components: nil,
						},
					}
				}

			}
			// send response back to discord
			if len(response.Data.Components) > 0 {
				a.logger.Info("response", "content", fmt.Sprintf("%+v", response.Data.Components[0]))
			}
			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				a.logger.Error("error responding to interaction", "error", err)
			}

		}
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case "opt":
			data := i.MessageComponentData()
			var userid, text, channel, display string
			splitted := strings.Split(data.Values[0], ";")
			a.logger.Info("splitted", "values", splitted, "len", len(splitted))
			if len(splitted) == 5 {
				channel = splitted[1]
				userid = splitted[2]
				text = fmt.Sprintf("cmd;%s;%s", splitted[3], splitted[4])
				display = splitted[3]
			}
			errorMessage := ""
			msg, err := a.postCommand(userid, text, channel)
			if err != nil {
				a.logger.Error("error posting to backend", "error", err.Error())
				errorMessage = "error posting to backend, try again in a few minutes"
			}
			returnMessage := msg.Msg
			if errorMessage != "" {
				returnMessage = errorMessage
			}
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Selected: %s; and answer: %s", display, returnMessage),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
			err = s.InteractionRespond(i.Interaction, response)
			if err != nil {
				a.logger.Error("error responding to interaction", "error", err)
			}

		case "choice":
			data := i.MessageComponentData()
			var userid, text, channel, display string
			splitted := strings.Split(data.Values[0], ";")
			a.logger.Info("splitted", "values", splitted, "len", len(splitted))
			if len(splitted) == 5 {
				channel = splitted[1]
				userid = splitted[2]
				text = fmt.Sprintf("choice;%s;%s", splitted[3], splitted[4])
				display = splitted[3]
			}
			errorMessage := ""
			msg, err := a.postCommand(userid, text, channel)
			if err != nil {
				a.logger.Error("error posting to backend", "error", err.Error())
				errorMessage = "error posting to backend, try again in a few minutes"
			}
			returnMessage := msg.Msg
			if errorMessage != "" {
				returnMessage = errorMessage
			}
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Selected: %s; and answer: %s", display, returnMessage),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
			err = s.InteractionRespond(i.Interaction, response)
			if err != nil {
				a.logger.Error("error responding to interaction", "error", err)
			}

		}
	}
}

func (a *app) events(w http.ResponseWriter, r *http.Request) {
	headerToken := r.Header.Get(types.HeaderToken)
	if headerToken != a.admToken {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "{\"msg\":\"unauthenticated\"}")
		return
	}

	// handle event
	obj := types.Event{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{\"msg\":\"json decode error\"}")
		return
	}
	if obj.Channel == "" || obj.UserID == "" || obj.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{\"msg\":\"missing required fields\"}")
		return
	}
	attachment := discordgo.MessageEmbed{}
	attachment.Description = obj.Message
	if obj.Kind != "" {
		switch obj.Kind {
		case types.EventAnnounce:
			emoji := ":mega:"
			attachment.Title = fmt.Sprintf("%s Announce", emoji)

		case types.EventSuccess:
			emoji := ":white_check_mark:"
			attachment.Title = fmt.Sprintf("%s Last Results", emoji)

		case types.EventFailure:
			emoji := ":x:"
			attachment.Title = fmt.Sprintf("%s Last Results", emoji)

		case types.EventDead:
			emoji := ":skull:"
			attachment.Title = fmt.Sprintf("%s Last Results", emoji)

		case types.EventInformation:
			emoji := ":information_source:"
			attachment.Title = fmt.Sprintf("%s Message", emoji)

		case types.EventEnd:
			emoji := ":stop_sign:"
			attachment.Title = fmt.Sprintf("%s End", emoji)
		}
	}
	res, err := a.session.ChannelMessageSendEmbed(obj.Channel, &attachment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.logger.Error("error sending message", "error", err)
		return
	}
	a.logger.Info("message sent", "message", res)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "{\"msg\":\"Accepted\"}")
}

func (a *app) postCommand(userid, text, channel string) (types.Composed, error) {
	var msg types.Composed
	body, err := a.web.PostCommand(userid, text, channel)
	if err != nil {
		a.logger.Error("post command", "error", err.Error())
		return msg, err
	}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		a.logger.Error("error decoding message from backend", "error", err.Error())
	}
	return msg, nil
}

func read(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
