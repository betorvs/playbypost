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
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var (
	Version string = "development"
)

type app struct {
	logger   *slog.Logger
	web      *cli.Cli
	slack    *slack.Client
	admToken string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// log := slog.NewLogLogger(logger.Handler(), slog.LevelInfo)
	logger.Info("starting slack bot test", "version", Version)
	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	playbypost := utils.GetEnv("PLAYBYPOST_SERVER", "http://localhost:3000")
	adminUser := utils.GetEnv("ADMIN_USER", "admin")
	adminFile := utils.GetEnv("CREDS_FILE", "./creds")
	adminToken, err := utils.Read(adminFile)
	if err != nil {
		logger.Error("error reading creds file", "error", err.Error())
	}
	logger.Info("debug", "token", adminToken)
	play := cli.NewHeaders(playbypost, adminUser, adminToken)
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	socket := socketmode.New(
		client,
		socketmode.OptionDebug(false),
		// socketmode.OptionLog(log),
	)
	slack := slack.New(token)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8091",
		Handler: mux,
	}
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "{\"status\":\"OK\"}")
	})

	// create internal app
	a := app{
		logger:   logger,
		web:      play,
		slack:    slack,
		admToken: adminToken,
	}

	mux.HandleFunc("POST /api/v1/event", a.events)
	mux.HandleFunc("GET /api/v1/validate", a.validate)

	err = a.web.Ping()
	if err != nil {
		logger.Error("error connecting with backend", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("connected to backend")

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen and serve error", "error", err)
			os.Exit(1)
		}
		logger.Info("stopped serving new connections.")
	}()

	socketmodeHandler := socketmode.NewSocketmodeHandler(socket)
	socketmodeHandler.Handle(socketmode.EventTypeEventsAPI, a.middlewareEventsAPI)
	socketmodeHandler.Handle(socketmode.EventTypeConnecting, a.middlewareConnecting)
	socketmodeHandler.Handle(socketmode.EventTypeConnected, a.middlewareConnected)
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, a.middlewareConnectionError)
	socketmodeHandler.Handle(socketmode.EventTypeHello, a.middlewareHello)

	socketmodeHandler.Handle(socketmode.EventTypeSlashCommand, a.middlewareSlashCommand)
	socketmodeHandler.HandleSlashCommand("/play-by-post", a.middlewareSlashCommand)
	socketmodeHandler.Handle(socketmode.EventTypeInteractive, a.middlewareInteractive)

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()
	go func() {
		err := socketmodeHandler.RunEventLoopContext(ctx)
		if err != nil {
			logger.Error("socket run event loop message", "error", err.Error())
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctxTimeout, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	if err := server.Shutdown(ctxTimeout); err != nil {
		logger.Error("shutdown error", "error", err)
	}
	logger.Info("graceful shutdown complete.")
}

func (a *app) middlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	a.logger.Info("Connecting to Slack with Socket Mode...")
}

func (a *app) middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	a.logger.Error("Connection failed. Retrying later...")
}

func (a *app) middlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	a.logger.Info("Connected to Slack with Socket Mode.")
}
func (a *app) middlewareHello(evt *socketmode.Event, client *socketmode.Client) {
	a.logger.Info("Hello message received with Socket Mode.")
}

func (a *app) middlewareEventsAPI(evt *socketmode.Event, client *socketmode.Client) {
	a.logger.Info("middlewareEventsAPI")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)

	a.logger.Info(fmt.Sprintf("Event received: %+v\n", eventsAPIEvent))
	attachment := slack.Attachment{}
	attachment.Color = "#4af030"
	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			a.logger.Info("We have been mentionned", "channel", ev.Channel, "user_id", ev.User)
			user, err := a.slack.GetUserInfo(ev.User)
			if err != nil {
				a.logger.Error("get user info error", "error", err.Error())
			}
			a.logger.Info(fmt.Sprintf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email))
			attachment.Text = fmt.Sprintf("Hello %s", user.Profile.RealName)
			switch {
			case strings.Contains(ev.Text, "help"):
				payload := helpMessage()
				_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionBlocks(payload...)) // PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
				if err != nil {
					a.logger.Error("failed to post message", "error", err.Error())
				}

			case strings.Contains(ev.Text, "join"):
				body, err := a.web.AddChatInformation(user.ID, user.Profile.RealName, ev.Channel, types.Slack)
				if err != nil {
					a.logger.Error("error adding user info", "error", err.Error())
					attachment.Text = fmt.Sprintf("Sorry, it did not work %s", user.Profile.RealName)
				} else {
					a.logger.Info("user join to playbypost", "username", user.Profile.RealName)
					var msg types.Msg
					err = json.Unmarshal(body, &msg)
					if err != nil {
						a.logger.Error("error json unmarshal", "error", err.Error())
					}
					if strings.Contains(msg.Msg, "already added") {
						attachment.Text = fmt.Sprintf("Already subscribed. Great, %s", user.Profile.RealName)
					}
					if strings.EqualFold(msg.Msg, "added") {
						attachment.Text = fmt.Sprintf("Let's play %s", user.Profile.RealName)
					}
				}
				_, _, err = client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
				if err != nil {
					a.logger.Error("failed to post message", "error", err.Error())
				}

			}

		case *slackevents.MemberJoinedChannelEvent:
			a.logger.Info(fmt.Sprintf("user %q joined to channel %q", ev.User, ev.Channel))

		}
	default:
		client.Debugf("unsupported Events API event received")
	}

	a.logger.Info("end of function")
}

func (a *app) middlewareSlashCommand(evt *socketmode.Event, client *socketmode.Client) {
	cmd, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	a.logger.Info(fmt.Sprintf("slash command from %v by %s", cmd.ChannelID, cmd.UserID))
	text := strings.ToLower(cmd.Text)
	noOptions := "No options for you"
	switch text {
	case "help":
		payload := helpMessage()
		payload2 := map[string]interface{}{
			"blocks": payload,
		}
		client.Ack(*evt.Request, payload2)
		return
	case "/iamwriter":
		writerUsername := strings.TrimSpace(strings.Replace(cmd.Text, "/iamwriter", "", 1))

		// Call backend to associate writer and user
		var responseContent string
		writer, err := a.web.GetWriterByUsername(writerUsername)
		if err != nil {
			a.logger.Error("error getting writer by username", "error", err.Error())
			responseContent = fmt.Sprintf("Sorry, it did not work to find writer %s", writerUsername)
		} else {
			user, err := a.web.GetUserByUserID(cmd.UserID)
			if err != nil {
				a.logger.Error("error getting user by userid", "error", err.Error())
				responseContent = fmt.Sprintf("Sorry, it did not work to find user %s", cmd.UserName)
			} else {
				_, err := a.web.CreateWriterUserAssociation(writer.ID, user.ID)
				if err != nil {
					a.logger.Error("error creating writer user association", "error", err.Error())
					responseContent = fmt.Sprintf("Sorry, it did not work to associate writer %s with user %s", writerUsername, cmd.UserName)
				} else {
					responseContent = fmt.Sprintf("Successfully associated you with writer: %s", writerUsername)
				}
			}
		}

		payload := map[string]interface{}{
			"text":          responseContent,
			"response_type": "ephemeral",
		}
		client.Ack(*evt.Request, payload)
		return
	case types.Opt: // opt
		msg, err := a.web.PostCommandComposed(cmd.UserID, types.Opt, cmd.ChannelID)
		if err != nil {
			a.logger.Error("error posting to backend", "error", err.Error())
		}

		textPickItem := "Pick an item "
		if msg.Msg != "" {
			textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
		}
		if len(msg.Opts) > 0 {
			a.logger.Info(fmt.Sprintf("options: %v", msg.Opts))
			options := []*slack.OptionBlockObject{}
			for _, v := range msg.Opts {
				options = append(options, &slack.OptionBlockObject{
					Text: &slack.TextBlockObject{
						Type:  "plain_text",
						Text:  v.Name,
						Emoji: boolPtr(true),
					},
					Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, cmd.ChannelID, cmd.UserID, v.Value, v.ID),
				})
			}
			payload := map[string]interface{}{
				"blocks": []slack.Block{
					slack.SectionBlock{
						Type: "section",
						Text: &slack.TextBlockObject{
							Type: "mrkdwn",
							Text: textPickItem,
						},
						Accessory: &slack.Accessory{
							SelectElement: &slack.SelectBlockElement{
								Type: "static_select",
								Placeholder: &slack.TextBlockObject{
									Type:  "plain_text",
									Text:  "Select an item",
									Emoji: boolPtr(true),
								},
								ActionID: "static_select-action",
								Options:  options,
							},
						},
					},
				},
			}

			client.Ack(*evt.Request, payload)
			return
		}
	case types.SoloDescribe, types.DidaticDescribe: // "solo-describe"
		a.logger.Info("describe", "text", text)
		msg, err := a.web.PostCommandComposed(cmd.UserID, text, cmd.ChannelID)
		if err != nil {
			a.logger.Error("error posting to backend", "error", err.Error())
		}
		title := "Play by Post Solo Adventures"
		if strings.Contains(text, types.Didatic) {
			title = "Play by Post Didatic Adventures"
		}
		payload := []slack.Block{
			slack.SectionBlock{
				Type: "header",
				Text: &slack.TextBlockObject{
					Type: "plain_text",
					Text: title,
				},
			},
		}
		if len(msg.Opts) > 0 {
			for _, v := range msg.Opts {
				payload = append(payload, slack.SectionBlock{
					Type: "section",
					Text: &slack.TextBlockObject{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*%s*\n%s", v.Name, v.Value),
					},
				})
			}
		} else {
			payload = append(payload, slack.SectionBlock{
				Type: "section",
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "No adventures available",
				},
			})
		}
		payload2 := map[string]interface{}{
			"blocks": payload,
		}
		client.Ack(*evt.Request, payload2)
		return

	case types.SoloStart, types.DidaticStart, types.DidaticJoin: // "solo-start"
		a.logger.Info("start or join", "text", text)
		msg, err := a.web.PostCommandComposed(cmd.UserID, text, cmd.ChannelID)
		if err != nil {
			a.logger.Error("error posting to backend", "error", err.Error())
		}

		textPickItem := "Pick an item "
		if msg.Msg != "" {
			textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
		}
		if len(msg.Opts) > 0 {
			a.logger.Info(fmt.Sprintf("options: %v", msg.Opts))
			options := []*slack.OptionBlockObject{}
			for _, v := range msg.Opts {
				options = append(options, &slack.OptionBlockObject{
					Text: &slack.TextBlockObject{
						Type:  "plain_text",
						Text:  v.Name,
						Emoji: boolPtr(true),
					},
					Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, cmd.ChannelID, cmd.UserID, v.Value, v.ID),
				})
			}
			payload := map[string]interface{}{
				"blocks": []slack.Block{
					slack.SectionBlock{
						Type: "section",
						Text: &slack.TextBlockObject{
							Type: "mrkdwn",
							Text: textPickItem,
						},
						Accessory: &slack.Accessory{
							SelectElement: &slack.SelectBlockElement{
								Type: "static_select",
								Placeholder: &slack.TextBlockObject{
									Type:  "plain_text",
									Text:  "Select an item",
									Emoji: boolPtr(true),
								},
								ActionID: "static_select-action",
								Options:  options,
							},
						},
					},
				},
			}

			client.Ack(*evt.Request, payload)
			return
		}
		if msg.Msg != "" && text == types.DidaticJoin {
			noOptions = msg.Msg
		}

	case types.SoloNext, types.DidaticNext: // "solo-next"
		a.logger.Info("game next", "text", text)
		msg, err := a.web.PostCommandComposed(cmd.UserID, text, cmd.ChannelID)
		if err != nil {
			a.logger.Error("error posting to backend", "error", err.Error())
		}

		textPickItem := "Pick an item "
		if msg.Msg != "" {
			textPickItem = fmt.Sprintf("%s Pick an item", msg.Msg)
		}
		if len(msg.Opts) > 0 {
			a.logger.Info(fmt.Sprintf("options: %v", msg.Opts))
			options := []*slack.OptionBlockObject{}
			for _, v := range msg.Opts {
				options = append(options, &slack.OptionBlockObject{
					Text: &slack.TextBlockObject{
						Type:  "plain_text",
						Text:  v.Name,
						Emoji: boolPtr(true),
					},
					Value: fmt.Sprintf(`cuni;%s;%s;%s;%d`, cmd.ChannelID, cmd.UserID, v.Value, v.ID),
				})
			}
			payload := map[string]interface{}{
				"blocks": []slack.Block{
					slack.SectionBlock{
						Type: "section",
						Text: &slack.TextBlockObject{
							Type: "mrkdwn",
							Text: textPickItem,
						},
						Accessory: &slack.Accessory{
							SelectElement: &slack.SelectBlockElement{
								Type: "static_select",
								Placeholder: &slack.TextBlockObject{
									Type:  "plain_text",
									Text:  "Select an item",
									Emoji: boolPtr(true),
								},
								ActionID: "static_select-action",
								Options:  options,
							},
						},
					},
				},
			}

			client.Ack(*evt.Request, payload)
			return
		}

	default:
		a.logger.Info("No options for you")
	}

	payload2 := map[string]interface{}{
		"blocks": []slack.Block{
			slack.SectionBlock{
				Type: "section",
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: noOptions,
				},
			},
		},
	}

	client.Ack(*evt.Request, payload2)
}

func (a *app) middlewareInteractive(evt *socketmode.Event, client *socketmode.Client) {
	interactiveEvent, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		a.logger.Debug(fmt.Sprintf("Ignored %+v\n", evt))
		return
	}
	values := ""
	for _, action := range interactiveEvent.ActionCallback.BlockActions {
		a.logger.Debug(fmt.Sprintf("action: %+v\n", action))
		if action.SelectedOption.Value != "" {
			a.logger.Debug(fmt.Sprintln("value: ", action.SelectedOption.Value))
			values = action.SelectedOption.Value
			break
		}
	}

	a.logger.Info("value received", "values", values)
	// check if command or solo
	startOpt := types.Cmd
	if strings.Contains(values, types.Solo) {
		startOpt = types.Choice
	}
	if strings.Contains(values, types.Didatic) {
		startOpt = types.Decision
	}
	var errorMessage, returnMessage string
	channel, userid, text, display, err := cli.ParserValues(values, startOpt)
	if err != nil {
		a.logger.Error("error parsing values", "error", err.Error())
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
	a.logger.Info("values", "userid", userid, "text", text, "channel", channel)
	attachment := slack.Attachment{
		Text: fmt.Sprintf("Selected: %s; and answer: %s", display, returnMessage),
	}

	_, _, err = a.slack.PostMessage(
		interactiveEvent.Container.ChannelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionResponseURL(interactiveEvent.ResponseURL, slack.ResponseTypeInChannel),
		slack.MsgOptionReplaceOriginal(interactiveEvent.ResponseURL),
	)
	if err != nil {
		a.logger.Error("error sending message to slack", "error", err.Error())
	}

	client.Ack(*evt.Request)
}

func helpMessage() []slack.Block {
	h := []slack.Block{
		slack.SectionBlock{
			Type: "header",
			Text: &slack.TextBlockObject{
				Type: "plain_text",
				Text: "Play by Post Help",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Play by Post Bot helps you play roleplaying games using text messages here in Slack. You can play a shared table RPG session using the playbypost slash command, or you can play a solo adventure using the solo commands. Or, if you are a student, you can use the didatic command to play an interesting adventure and learn something special.",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `@play-by-post join` to register yourself to play as Storyteller or as Player",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `/play-by-post options` to get your options as a Player or Storyteller",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `/play-by-post solo-start` to start a solo adventure",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `/play-by-post solo-next` to get your options in your solo adventure",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `/play-by-post solo-describe` to get a list of solo adventure with a description",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `/play-by-post didatic-start` to start a didatic adventure",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `/play-by-post didatic-join` to join in a didatic adventure",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `/play-by-post didatic-next` to get your options in your didatic adventure",
			},
		},
		slack.SectionBlock{
			Type: "section",
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "Use: `/play-by-post didatic-describe` to get a list of didatic adventure with a description",
			},
		},
	}
	// }
	return h
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
	attachment := slack.Attachment{}
	attachment.Text = obj.Message
	if obj.Kind != "" {
		switch obj.Kind {
		case types.EventAnnounce:
			emoji := ":mega:"
			attachment.Text = fmt.Sprintf("%s %s", emoji, obj.Message)

		case types.EventSuccess:
			emoji := ":white_check_mark:"
			attachment.Text = fmt.Sprintf("%s %s", emoji, obj.Message)

		case types.EventFailure:
			emoji := ":x:"
			attachment.Text = fmt.Sprintf("%s %s", emoji, obj.Message)

		case types.EventDead:
			emoji := ":skull:"
			attachment.Text = fmt.Sprintf("%s %s", emoji, obj.Message)

		case types.EventInformation:
			emoji := ":information_source:"
			attachment.Text = fmt.Sprintf("%s %s", emoji, obj.Message)

		case types.EventEnd:
			emoji := ":end:"
			attachment.Text = fmt.Sprintf("%s %s", emoji, obj.Message)
		}
	}
	if obj.ImageURL != "" {
		attachment.ImageURL = obj.ImageURL
	}
	switch {
	case strings.ToLower(obj.UserID) == "all":
		// send message to all in channel
		a.logger.Info("send message to all in channel", "channel", obj.Channel)
		_, _, err := a.slack.PostMessage(obj.Channel, slack.MsgOptionAttachments(attachment))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "{\"msg\":\"cannot send message to slack\"}")
			return
		}
	default:
		// send private message
		a.logger.Info("send private message", "channel", obj.Channel, "user", obj.UserID)
		_, err := a.slack.PostEphemeral(obj.Channel, obj.UserID, slack.MsgOptionAttachments(attachment))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "{\"msg\":\"cannot send private message to slack\"}")
			return
		}
	}

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

func boolPtr(b bool) *bool { return &b }
