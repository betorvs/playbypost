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
			Description: "Join to Play by Post to be a Storyteller or a Player",
		},
		{
			Name: "play-by-post",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "help",
					Description: "show help message",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "options",
					Description: "Options menu for play by post Story used by Storyteller and Players",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "solo-start",
					Description: "Call it to to start a solo game",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "solo-next",
					Description: "Get next select menu for your solo game",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "solo-describe",
					Description: "Describe your solo game options",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "didatic-join",
					Description: "Join to a Didatic Adventure",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "didatic-start",
					Description: "Call it to to start a didatic game",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "didatic-next",
					Description: "Get next select menu for your didatic game",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "didatic-describe",
					Description: "Describe your didatic game options",
				},
			},
			Description: "Main Play by Post command",
		},
		{
			Name:        "iamwriter",
			Description: "Identify yourself as a writer",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "writer_username",
					Description: "The username of the writer you want to associate with",
					Required:    true,
				},
			},
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
	adminToken, err := utils.Read(adminFile)
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

	err = a.web.Ping()
	if err != nil {
		logger.Error("error connecting with backend", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("connected to backend")

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
		_, _ = fmt.Fprint(w, "{\"status\":\"OK\"}")
	})
	mux.HandleFunc("POST /api/v1/event", a.events)
	mux.HandleFunc("GET /api/v1/validate", a.validate)

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
	err = discord.Close()
	if err != nil {
		logger.Error("error closing discord session", "error", err)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	if err := server.Shutdown(ctxTimeout); err != nil {
		logger.Error("shutdown error", "error", err)
	}
	logger.Info("graceful shutdown complete.")
}

func (a *app) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	a.logger.Info("message received", "message", m.Content)
	switch {
	case strings.Contains(strings.ToLower(m.Content), "hello"):
		message := "Hello, I am Play by Post bot. How can I help you? Try `help` to get more options. ;)"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println("Error sending message: ", err)
		}
	case strings.Contains(strings.ToLower(m.Content), "help"):
		content, embed := helpMessage()
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: content,
			Embed:   embed[0],
		})
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

		case "iamwriter":
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			writerUsername := optionMap["writer_username"].StringValue()

			// Call backend to associate writer and user
			var responseContent string
			writer, err := a.web.GetWriterByUsername(writerUsername)
			if err != nil {
				a.logger.Error("error getting writer by username", "error", err.Error())
				responseContent = fmt.Sprintf("Sorry, it did not work to find writer %s", writerUsername)
			} else {
				user, err := a.web.GetUserByUserID(userid)
				if err != nil {
					a.logger.Error("error getting user by userid", "error", err.Error())
					responseContent = fmt.Sprintf("Sorry, it did not work to find user %s", username)
				} else {
					_, err := a.web.CreateWriterUserAssociation(writer.ID, user.ID)
					if err != nil {
						a.logger.Error("error creating writer user association", "error", err.Error())
						responseContent = fmt.Sprintf("Sorry, it did not work to associate writer %s with user %s", writerUsername, username)
					} else {
						responseContent = fmt.Sprintf("Successfully associated you with writer: %s", writerUsername)
					}
				}
			}

			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:    responseContent,
					Flags:      discordgo.MessageFlagsEphemeral,
					Components: nil,
				},
			}
			err = s.InteractionRespond(i.Interaction, response)
			if err != nil {
				a.logger.Error("error responding to interaction", "error", err)
			}

		case "play-by-post":
			var response *discordgo.InteractionResponse
			text := i.ApplicationCommandData().Options[0].Name
			switch text {
			case "help":
				content, embed := helpMessage()
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Title:      "Play by Post Help",
						Content:    content,
						Flags:      discordgo.MessageFlagsEphemeral,
						Components: nil,
						Embeds:     embed,
					},
				}
			case "options", types.Opt:
				a.logger.Info("options")
				// post command
				msg, err := a.web.PostCommandComposed(userid, types.Opt, i.ChannelID)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
				}
				textPickItem := "Pick an item "
				if msg.Msg != "" {
					textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
				}
				if len(msg.Opts) > 0 {
					// create select menu
					options := []discordgo.SelectMenuOption{}
					for _, v := range msg.Opts {
						options = append(options, discordgo.SelectMenuOption{
							Label: v.Name,
							Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, i.ChannelID, userid, v.Value, v.ID),
						})
					}
					selectMenu := discordgo.SelectMenu{
						CustomID:    types.Opt,
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

			case types.SoloDescribe, types.DidaticDescribe: // solo-describe
				a.logger.Info("describe", "text", text)
				title := "Play by Post Solo Adventures"
				if strings.Contains(text, types.Didatic) {
					title = "Play by Post Didatic Adventures"
				}
				// post command
				msg, err := a.web.PostCommandComposed(userid, text, i.ChannelID)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
				}

				embed := discordgo.MessageEmbed{
					Title:       title,
					Description: "List all published adventures",
				}
				if len(msg.Opts) > 0 {
					for _, v := range msg.Opts {
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:  v.Name,
							Value: v.Value,
						})
					}
				} else {
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:  "No adventures available",
						Value: "-",
					})
				}
				slice := []*discordgo.MessageEmbed{&embed}

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Title:      "Play by Post Describe",
						Content:    "describe",
						Flags:      discordgo.MessageFlagsEphemeral,
						Components: nil,
						Embeds:     slice,
					},
				}

			case types.SoloStart, types.DidaticStart, types.DidaticJoin: // solo-start
				a.logger.Info("start or join", "text", text)
				// post command
				msg, err := a.web.PostCommandComposed(userid, text, i.ChannelID)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
				}
				textPickItem := "Pick an item "
				if msg.Msg != "" {
					textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
				}
				a.logger.Info("msg", "msg", msg)
				startOpt := types.Choice
				if strings.Contains(text, types.Didatic) {
					startOpt = types.Decision
				}
				if len(msg.Opts) > 0 {
					// create select menu
					options := []discordgo.SelectMenuOption{}
					for _, v := range msg.Opts {
						options = append(options, discordgo.SelectMenuOption{
							Label: v.Name,
							Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, i.ChannelID, userid, v.Value, v.ID),
						})
					}
					selectMenu := discordgo.SelectMenu{
						CustomID:    startOpt,
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
					noOptions := "No options available"
					if msg.Msg != "" {
						noOptions = msg.Msg
					}
					response = &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content:    noOptions,
							Flags:      discordgo.MessageFlagsEphemeral,
							Components: nil,
						},
					}
				}

			case types.SoloNext, types.DidaticNext: // solo-next
				a.logger.Info("next", "text", text)
				// post command
				msg, err := a.web.PostCommandComposed(userid, text, i.ChannelID)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
				}
				textPickItem := "Pick an item "
				if msg.Msg != "" {
					textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
				}
				startOpt := types.Choice
				if strings.Contains(text, types.Didatic) {
					startOpt = types.Decision
				}
				if len(msg.Opts) > 0 {
					// create select menu
					options := []discordgo.SelectMenuOption{}
					for _, v := range msg.Opts {
						options = append(options, discordgo.SelectMenuOption{
							Label: v.Name,
							Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, i.ChannelID, userid, v.Value, v.ID),
						})
					}
					selectMenu := discordgo.SelectMenu{
						CustomID:    startOpt,
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
			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				a.logger.Error("error responding to interaction", "error", err)
			}

		}
	case discordgo.InteractionMessageComponent:
		customID := i.MessageComponentData().CustomID
		switch customID {
		case types.Opt: // opt
			data := i.MessageComponentData()
			startOpt := types.Cmd
			var errorMessage, returnMessage string
			channel, userid, text, display, err := cli.ParserValues(data.Values[0], startOpt)
			if err != nil {
				a.logger.Error("error parsing values", "error", err)
				errorMessage = "error parsing values from backend"
			} else {
				msg, err := a.web.PostCommandComposed(userid, text, channel)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
					errorMessage = "error posting to backend, try again in a few minutes"
				}
				returnMessage = msg.Msg
			}
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

		case types.Choice, types.Decision: // choice
			data := i.MessageComponentData()
			startOpt := types.Choice
			if customID == types.Decision {
				startOpt = types.Decision
			}
			var errorMessage, returnMessage string
			channel, userid, text, display, err := cli.ParserValues(data.Values[0], startOpt)
			if err != nil {
				a.logger.Error("error parsing values", "error", err)
				errorMessage = "error parsing values from backend"
			} else {
				msg, err := a.web.PostCommandComposed(userid, text, channel)
				if err != nil {
					a.logger.Error("error posting to backend", "error", err.Error())
					errorMessage = "error posting to backend, try again in a few minutes"
				}
				returnMessage = msg.Msg
			}
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
		_, _ = fmt.Fprint(w, "{\"msg\":\"unauthenticated\"}")
		return
	}

	// handle event
	obj := types.Event{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "{\"msg\":\"json decode error\"}")
		return
	}
	if obj.Channel == "" || obj.UserID == "" || obj.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "{\"msg\":\"missing required fields\"}")
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
	if obj.ImageURL != "" {
		attachment.Image = &discordgo.MessageEmbedImage{URL: obj.ImageURL}
	}
	res, err := a.session.ChannelMessageSendEmbed(obj.Channel, &attachment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.logger.Error("error sending message", "error", err)
		return
	}
	a.logger.Info("message sent", "message", res)

	w.WriteHeader(http.StatusAccepted)
	_, _ = fmt.Fprint(w, "{\"msg\":\"Accepted\"}")
}

func (a *app) validate(w http.ResponseWriter, r *http.Request) {
	headerToken := r.Header.Get(types.HeaderToken)
	if headerToken != a.admToken {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w, "{\"msg\":\"unauthenticated\"}")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "{\"msg\":\"authenticated\"}")
}

func helpMessage() (string, []*discordgo.MessageEmbed) {
	header := "help message from Play by Post Bot"
	content := "Play by Post Bot helps you play roleplaying games using text messages here in Slack. You can play a shared table RPG session using the playbypost slash command, or you can play a solo adventure using the solo commands. Or, if you are a student, you can use the didatic command to play an interesting adventure and learn something special."
	embed := &discordgo.MessageEmbed{
		Title:       "Play by Post Help",
		Description: content,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Join",
				Value: "Use `/join` to register yourself to play as Storyteller or as Player",
			},
			{
				Name:  "options",
				Value: "Use `/play-by-post options` to get your options as a Player or Storyteller",
			},
			{
				Name:  "solo-start",
				Value: "Use `/play-by-post solo-start` to start a solo adventure",
			},
			{
				Name:  "solo-next",
				Value: "Use `/play-by-post solo-next` to get your options in your solo adventure",
			},
			{
				Name:  "solo-describe",
				Value: "Use `/play-by-post solo-describe` to get a list of solo adventure with a description",
			},
			{
				Name:  "didatic-start",
				Value: "Use `/play-by-post didatic-start` to start a didatic adventure",
			},
			{
				Name:  "didatic-join",
				Value: "Use `/play-by-post didatic-join` to join in a didatic adventure",
			},
			{
				Name:  "didatic-next",
				Value: "Use `/play-by-post didatic-next` to get your options in your didatic adventure",
			},
			{
				Name:  "didatic-describe",
				Value: "Use `/play-by-post didatic-describe` to get a list of didatic adventure with a description",
			},
		},
	}
	slice := []*discordgo.MessageEmbed{embed}
	return header, slice
}
