CREATE TABLE users (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  userid VARCHAR(50) UNIQUE NOT NULL,
  encoding_key VARCHAR(16) NOT NULL,
  password text NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE story (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title VARCHAR(50) UNIQUE NOT NULL,
  notes text,
  announcement text,
  finished BOOLEAN NOT NULL DEFAULT FALSE,
  master_id int NOT NULL REFERENCES users(id),
  rpg VARCHAR(50)
);

CREATE TABLE story_channels (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  story_id int NOT NULL REFERENCES story(id),
  channel VARCHAR(25),
  UNIQUE(channel, story_id)
);

CREATE TABLE players (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  character_name VARCHAR(255) NOT NULL,
  story_id int NOT NULL REFERENCES story(id),
  player_id int NOT NULL REFERENCES users(id),
  destroyed BOOLEAN NOT NULL DEFAULT FALSE,
  abilities JSONB,
  skills JSONB,
  rpg VARCHAR(50)
);

CREATE TABLE non_players (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  npc_name VARCHAR(255) NOT NULL,
  story_id int NOT NULL REFERENCES story(id),
  master_id int NOT NULL REFERENCES users(id),
  destroyed BOOLEAN NOT NULL DEFAULT FALSE,
  abilities JSONB,
  skills JSONB,
  rpg VARCHAR(50)
);

CREATE TABLE encounters (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title VARCHAR(50) UNIQUE NOT NULL,
  display_text VARCHAR(25) NOT NULL,
  notes text,
  announcement text,
  reward text,
  xp int NOT NULL DEFAULT 0,
  phase int NOT NULL DEFAULT 0,
  finished BOOLEAN NOT NULL DEFAULT FALSE,
  story_id int NOT NULL REFERENCES story(id),
  master_id int NOT NULL REFERENCES users(id),
  UNIQUE(display_text, story_id)
);

CREATE TABLE encounters_participants_players (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  players_id int NOT NULL REFERENCES players(id),
  encounters_id int NOT NULL REFERENCES encounters(id),
  UNIQUE(players_id, encounters_id)
);

CREATE TABLE encounters_participants_non_players (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  non_players_id int NOT NULL REFERENCES non_players(id),
  encounters_id int NOT NULL REFERENCES encounters(id),
  UNIQUE(non_players_id, encounters_id)
);


CREATE TABLE initiative (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title VARCHAR(50) UNIQUE NOT NULL,
  next_player int NOT NULL,
  encounters_id int NOT NULL REFERENCES encounters(id)
);

CREATE TABLE initiative_participants (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  initiative_id int NOT NULL REFERENCES initiative(id),
  participant_name VARCHAR(255) NOT NULL,
  participant_bonus int NOT NULL,
  participant_result int NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE tasks (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  encounters_id int NOT NULL REFERENCES encounters(id),
  title VARCHAR(50) UNIQUE NOT NULL,
  display_text VARCHAR(25) NOT NULL,
  kind int NOT NULL DEFAULT 0,
  checks VARCHAR(100) NOT NULL,
  target int NOT NULL DEFAULT 0,
  options JSONB,
  finished BOOLEAN NOT NULL DEFAULT FALSE,
  UNIQUE(display_text, encounters_id)
);