package cmd

var (
	username      string // cmd: users
	password      string // cmd: users
	title         string // cmd: story, encounter
	displayText   string // cmd: encounter, stage
	announcement  string // cmd: story, encounter
	notes         string // cmd: story, encounter
	writerID      int    // cmd: story
	name          string // cmd: player
	playerid      int    // cmd: player
	storyid       int    // cmd: player, encounter
	encounterid   int    // cmd: encounter, initiative
	listPlayersID []int  // cmd: encounter
	isNPC         bool   //cmd: encounter
	adminToken    string //cmd: all
	adminUser     string //cmd: all
	server        string //cmd: all
	description   string //cmd: tasks
	ability       string //cmd: task
	skill         string //cmd: task
	kind          int    //cmd: task
	target        int    //cmd: task
	storyID       int    //cmd: stage
	userID        string //cmd: stage
	channel       string // cmd: initiative
)
