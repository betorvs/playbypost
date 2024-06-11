package cmd

var (
	username      string // cmd: users
	userid        string // cmd: users
	password      string // cmd: users
	title         string // cmd: story, encounter
	displayText   string // cmd: encounter
	announcement  string // cmd: story, encounter
	notes         string // cmd: story, encounter
	masterID      int    // cmd: story
	name          string // cmd: player, initiative
	playerid      int    // cmd: player
	storyid       int    // cmd: player, encounter
	encounterid   int    // cmd: encounter, initiative
	listPlayersID []int  // cmd: encounter
	isNPC         bool   //cmd: encounter, initiative
	adminToken    string //cmd: all
	adminUser     string //cmd: all
	server        string //cmd: all
)
