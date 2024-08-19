CREATE TABLE writers (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password text NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE story (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title VARCHAR(50) UNIQUE NOT NULL,
  notes text,
  announcement text,
  writer_id int NOT NULL REFERENCES writers(id),
  rpg VARCHAR(25),
  UNIQUE(title, writer_id)
);

CREATE TABLE story_keys (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  encoding_key VARCHAR(16) NOT NULL,
  story_id int NOT NULL REFERENCES story(id)
);

CREATE TABLE access_story_keys (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  writer_id int NOT NULL REFERENCES writers(id),
  story_keys_id int NOT NULL REFERENCES story_keys(id)
);

CREATE TABLE encounters (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title VARCHAR(50) UNIQUE NOT NULL,
  notes text,
  announcement text,
  story_id int NOT NULL REFERENCES story(id),
  writer_id int NOT NULL REFERENCES writers(id),
  UNIQUE(title, story_id)
);

CREATE TABLE tasks (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  description VARCHAR(50) UNIQUE NOT NULL,
  kind int NOT NULL DEFAULT 0,
  ability VARCHAR(50) NOT NULL,
  skill VARCHAR(50) NOT NULL,
  target int NOT NULL DEFAULT 0
);

CREATE TABLE chat_information (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  userid VARCHAR(50) NOT NULL,
  channel VARCHAR(50) NOT NULL,
  username VARCHAR(50) NOT NULL,
  chat VARCHAR(50) NOT NULL
);

CREATE TABLE users (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  userid VARCHAR(50) UNIQUE NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE stage (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  display_text VARCHAR(50) NOT NULL,
  encoding_key VARCHAR(16) NOT NULL,
  finished BOOLEAN NOT NULL DEFAULT FALSE,
  storyteller_id int NOT NULL REFERENCES users(id),
  story_id int NOT NULL REFERENCES story(id)
);

CREATE TABLE stage_channel (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  channel VARCHAR(50) UNIQUE NOT NULL,
  stage_id int NOT NULL REFERENCES stage(id),
  active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE stage_encounters (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  display_text VARCHAR(50) NOT NULL,
  phase int NOT NULL DEFAULT 0,
  stage_id int NOT NULL REFERENCES stage(id),
  storyteller_id int NOT NULL REFERENCES users(id),
  encounters_id int NOT NULL REFERENCES encounters(id),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  UNIQUE(display_text, encounters_id)
);

CREATE TABLE stage_running_tasks (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  display_text VARCHAR(50) NOT NULL,
  stage_id int NOT NULL REFERENCES stage(id),
  storyteller_id int NOT NULL REFERENCES users(id),
  stage_encounters_id int NOT NULL REFERENCES stage_encounters(id),
  task_id int NOT NULL REFERENCES tasks(id),
  UNIQUE(display_text, stage_encounters_id, task_id)
);

CREATE TABLE stage_next_encounter (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  display_text VARCHAR(50) NOT NULL,
  stage_id int NOT NULL REFERENCES stage(id),
  current_encounter_id int NOT NULL REFERENCES stage_encounters(id),
  next_encounter_id int NOT NULL REFERENCES stage_encounters(id),
  UNIQUE(display_text, current_encounter_id, next_encounter_id)
);

CREATE TABLE stage_encounter_activities (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  stage_id int NOT NULL REFERENCES stage(id),
  encounter_id int NOT NULL REFERENCES stage_encounters(id),
  actions JSONB,
  processed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE players (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  character_name VARCHAR(255) NOT NULL,
  stage_id int NOT NULL REFERENCES stage(id),
  player_id int NOT NULL REFERENCES users(id),
  destroyed BOOLEAN NOT NULL DEFAULT FALSE,
  abilities JSONB,
  skills JSONB,
  extensions JSONB,
  rpg VARCHAR(50)
);

CREATE TABLE non_players (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  npc_name VARCHAR(255) NOT NULL,
  stage_id int NOT NULL REFERENCES stage(id),
  storyteller_id int NOT NULL REFERENCES users(id),
  destroyed BOOLEAN NOT NULL DEFAULT FALSE,
  abilities JSONB,
  skills JSONB,
  extensions JSONB,
  rpg VARCHAR(50)
);

CREATE TABLE stage_encounters_participants_players (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  players_id int NOT NULL REFERENCES players(id),
  stage_encounters_id int NOT NULL REFERENCES stage_encounters(id),
  UNIQUE(players_id, stage_encounters_id)
);

CREATE TABLE stage_encounters_participants_non_players (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  non_players_id int NOT NULL REFERENCES non_players(id),
  stage_encounters_id int NOT NULL REFERENCES stage_encounters(id),
  UNIQUE(non_players_id, stage_encounters_id)
);

CREATE TABLE initiative (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title VARCHAR(50) UNIQUE NOT NULL,
  next_player int NOT NULL,
  stage_encounters_id int NOT NULL REFERENCES stage_encounters(id)
);

CREATE TABLE initiative_participants (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  initiative_id int NOT NULL REFERENCES initiative(id),
  participant_name VARCHAR(255) NOT NULL,
  participant_bonus int NOT NULL,
  participant_result int NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE
);