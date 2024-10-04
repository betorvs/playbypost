package cmd

const (
	formatJSON = "json"
	formatLog  = "log"
)

var (
	username           string // cmd: users
	password           string // cmd: users
	title              string // cmd: story, encounter
	displayText        string // cmd: encounter, stage
	announcement       string // cmd: story, encounter
	notes              string // cmd: story, encounter
	writerID           int    // cmd: story
	name               string // cmd: player
	playerID           int    // cmd: player
	playerUserID       string // cmd: player
	storyID            int    // cmd: player, encounter, stage
	stageTitle         string // cmd: player, stage
	stageID            int    // cmd: player
	encounterID        int    // cmd: encounter, initiative
	listPlayersID      []int  // cmd: encounter
	isNPC              bool   //cmd: encounter
	firstEncounter     bool   //cmd: encounter
	lastEncounter      bool   //cmd: encounter
	adminToken         string //cmd: all
	adminUser          string //cmd: all
	server             string //cmd: all
	description        string //cmd: tasks
	ability            string //cmd: task
	skill              string //cmd: task
	kind               int    //cmd: task
	target             int    //cmd: task
	userID             string //cmd: stage
	storyTitle         string //cmd: stage
	channel            string // cmd: initiative
	solo               bool   // cmd: auto-play
	autoPlayID         int    // cmd: auto-play
	nextEncounterID    int    // cmd: auto-play
	encounterTitle     string // cmd: auto-play, stage
	nextEncounterTitle string // cmd: auto-play
	chatUserID         string // cmd: chat
	chatChannelID      string // cmd: chat
	chatUserName       string // cmd: chat
	random             bool   // cmd: db
	objectiveKind      string // cmd: auto-play
	objectiveValues    []int  // cmd: auto-play
	outputFormat       string // cmd: list all
	objectKind         string // cmd: validator
	objectID           int    // cmd: validator
)
