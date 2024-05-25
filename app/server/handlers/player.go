package handlers

// const (
// 	active   string = "active"
// 	finished string = "finished"
// )

// type PlayerApi struct {
// 	ValidUsers    mappingUsers
// 	ValidStory    mappingStories
// 	ValidChannels mappingChannels
// 	logger        *slog.Logger
// 	s             *server.SvrWeb
// 	db            db.DBClient
// 	ctx           context.Context
// 	dice          rpg.Roll
// 	rpg           *rpg.RPGSystem
// }

// type mappingUsers struct {
// 	users map[string]int
// 	mu    *sync.Mutex
// }

// func (m *mappingUsers) Add(index string, value int) {
// 	m.mu.Lock()
// 	m.users[index] = value
// 	m.mu.Unlock()
// }

// func (m *mappingUsers) load(ctx context.Context, db db.DBClient) error {
// 	obj, err := db.GetUsers(ctx, true)
// 	if err != nil {
// 		return err
// 	}
// 	for _, v := range obj {
// 		if v.Username != "" {
// 			m.Add(v.UserID, v.ID)
// 		}
// 	}
// 	return nil
// }

// type mappingStories struct {
// 	stories map[string]int
// 	mu      *sync.Mutex
// }

// func (m *mappingStories) Add(index string, value int) {
// 	m.mu.Lock()
// 	m.stories[index] = value
// 	m.mu.Unlock()
// }

// func (m *mappingStories) load(ctx context.Context, db db.DBClient) error {
// 	obj, err := db.GetStory(ctx, true)
// 	if err != nil {
// 		return err
// 	}
// 	for _, v := range obj {
// 		if v.Title != "" {
// 			m.Add(v.Title, v.ID)
// 		}
// 	}
// 	return nil
// }

// type mappingChannels struct {
// 	channels map[string]int
// 	mu       *sync.Mutex
// }

// func (m *mappingChannels) Add(index string, value int) {
// 	m.mu.Lock()
// 	m.channels[index] = value
// 	m.mu.Unlock()
// }

// func (m *mappingChannels) load(ctx context.Context, db db.DBClient) error {
// 	obj, err := db.GetStoryChannels(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	for k, v := range obj {
// 		m.Add(k, v)
// 	}
// 	return nil
// }

// // type mappingPlayers struct {
// // 	players map[string]int
// // 	mu      *sync.Mutex
// // }

// // func (m *mappingPlayers) Add(index string, value int) {
// // 	m.mu.Lock()
// // 	m.players[index] = value
// // 	m.mu.Unlock()
// // }

// // func (m *mappingPlayers) load(ctx context.Context, db db.DBClient) error {
// // 	obj, err := db.GetPlayersByStoryID(ctx, true)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	for _, v := range obj {
// // 		if v.Title != "" {
// // 			m.Add(v.Title, v.ID)
// // 		}
// // 	}
// // 	return nil
// // }

// func NewPlayerApi(ctx context.Context, dice rpg.Roll, db db.DBClient, l *slog.Logger, s *server.SvrWeb, rpg *rpg.RPGSystem) *PlayerApi {
// 	users := mappingUsers{
// 		users: make(map[string]int),
// 		mu:    &sync.Mutex{},
// 	}
// 	users.load(ctx, db)
// 	story := mappingStories{
// 		stories: make(map[string]int),
// 		mu:      &sync.Mutex{},
// 	}
// 	story.load(ctx, db)
// 	channels := mappingChannels{
// 		channels: make(map[string]int),
// 		mu:       &sync.Mutex{},
// 	}
// 	channels.load(ctx, db)
// 	return &PlayerApi{
// 		ValidUsers:    users,
// 		ValidStory:    story,
// 		ValidChannels: channels,
// 		ctx:           ctx,
// 		db:            db,
// 		logger:        l,
// 		s:             s,
// 		dice:          dice,
// 		rpg:           rpg,
// 	}
// }

// func (a *PlayerApi) checkAuth(r *http.Request) bool {
// 	headerUserID := r.Header.Get(types.HeaderUserID)
// 	headerStory := r.Header.Get(types.HeaderStory)
// 	headerStoryChannel := r.Header.Get(types.HeaderStoryChannel)
// 	headerPlayerID := r.Header.Get(types.HeaderPlayerID)
// 	switch {
// 	case headerUserID != "" && headerStoryChannel != "":
// 		return false
// 	case headerUserID != "" && headerStory != "" && headerPlayerID != "":
// 		return false
// 	}

// 	return true
// }

// func (p *PlayerApi) ReloadCache(w http.ResponseWriter, r *http.Request) {
// 	p.logger.Info("reloading cache")
// 	err := p.ValidUsers.load(p.ctx, p.db)
// 	if err != nil {
// 		p.s.ErrJSON(w, http.StatusBadRequest, "reload users failed")
// 		return
// 	}
// 	err = p.ValidStory.load(p.ctx, p.db)
// 	if err != nil {
// 		p.s.ErrJSON(w, http.StatusBadRequest, "reload story failed")
// 		return
// 	}
// 	p.s.JSON(w, types.Msg{Msg: "OK"})
// }

// func (p *PlayerApi) Commands(w http.ResponseWriter, r *http.Request) {
// 	if p.checkAuth(r) {
// 		p.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
// 		return
// 	}
// 	headerUserID := r.Header.Get(types.HeaderUserID)
// 	headerStory := r.Header.Get(types.HeaderStory)
// 	headerStoryChannel := r.Header.Get(types.HeaderStoryChannel)
// 	headerPlayerID := r.Header.Get(types.HeaderPlayerID)
// 	var playerID int
// 	if headerPlayerID != "" {
// 		id, err := strconv.Atoi(headerPlayerID)
// 		if err != nil {
// 			p.logger.Error("command strconv error", "error", err.Error())
// 			p.s.ErrJSON(w, http.StatusBadRequest, "bad headers")
// 			return
// 		}
// 		playerID = id
// 	}

// 	obj := types.Command{}
// 	err := json.NewDecoder(r.Body).Decode(&obj)
// 	if err != nil {
// 		p.logger.Error("command json newdecoder", "error", err.Error())
// 		p.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
// 		return
// 	}
// 	var userID int
// 	if u, ok := p.ValidUsers.users[headerUserID]; ok {
// 		userID = u
// 	} else {
// 		p.logger.Info("user not found in cache")
// 		user, err := p.db.GetUserByUserID(p.ctx, headerUserID)
// 		if err != nil {
// 			p.s.ErrJSON(w, http.StatusBadRequest, "invalid user")
// 			return
// 		}
// 		userID = user.ID
// 	}
// 	var storyID int
// 	switch {
// 	case headerStory != "":
// 		{
// 			if u, ok := p.ValidStory.stories[headerStory]; ok {
// 				storyID = u
// 			} else {
// 				p.logger.Info("story not found in cache")
// 				story, err := p.db.GetStoryIDByTitle(p.ctx, headerStory)
// 				if err != nil {
// 					p.s.ErrJSON(w, http.StatusBadRequest, "invalid user")
// 					return
// 				}
// 				storyID = story
// 			}
// 		}
// 	case headerStoryChannel != "":
// 		{
// 			if u, ok := p.ValidChannels.channels[headerStoryChannel]; ok {
// 				storyID = u
// 			} else {
// 				p.logger.Info("story not found by channel id in cache")
// 				// story, err := p.db.GetStoryIDByTitle(p.ctx, headerStory)
// 				// if err != nil {
// 				// 	p.s.ErrJSON(w, http.StatusBadRequest, "invalid user")
// 				// 	return
// 				// }
// 				// storyID = story
// 			}
// 		}
// 	}

// 	if userID == 0 || storyID == 0 {
// 		p.s.ErrJSON(w, http.StatusForbidden, "invalid credentials")
// 		return
// 	}
// 	var player *rules.Creature
// 	if playerID != 0 {
// 		player, err = p.db.GetPlayer(p.ctx, playerID, false, p.rpg)
// 		if err != nil {
// 			p.s.ErrJSON(w, http.StatusBadRequest, "invalid player")
// 			return
// 		}
// 	} else {
// 		player, err = p.db.GetPlayerByUserID(p.ctx, userID, false, p.rpg)
// 		if err != nil {
// 			p.s.ErrJSON(w, http.StatusBadRequest, "invalid player")
// 			return
// 		}
// 	}

// 	initiative, initiativeID, err := p.db.GetRunningInitiativeByStoryID(p.ctx, storyID)
// 	if err != nil {
// 		p.s.ErrJSON(w, http.StatusBadRequest, "initiative not found or with error")
// 		return
// 	}
// 	next := initiative.Next()
// 	if initiative.Participants[next].Name == player.Name {
// 		// requires update turn position in initiative
// 		info := initiative.NextInfo()
// 		p.logger.Info("initiative info", "current", next, "next_player", initiative.Participants[info].Name)

// 		// check action
// 		command, err := parser.TextToCommand(obj.Text)
// 		if err != nil {
// 			p.logger.Error("error validatind player command", "error", err.Error())
// 			p.s.ErrJSON(w, http.StatusBadRequest, "invalid action")
// 			return
// 		}

// 		// saving position before leaving
// 		err = p.db.UpdateNextPlayer(p.ctx, initiativeID, initiative.Position)
// 		if err != nil {
// 			p.logger.Error("error updating initiative", "error", err.Error())
// 			p.s.ErrJSON(w, http.StatusBadRequest, "update initiative error")
// 			return
// 		}
// 		// return message
// 		msg := fmt.Sprintf("userID found '%d' and story found '%d' and player '%s' command '%v' -> '%s'", userID, storyID, player.Name, command.Act.String(), command.NotAct)
// 		//
// 		p.s.JSON(w, types.Msg{Msg: msg})
// 		return
// 	}
// 	msg := fmt.Sprintf("not your turn player %s", player.Name)
// 	p.s.JSON(w, types.Msg{Msg: msg})
// }
