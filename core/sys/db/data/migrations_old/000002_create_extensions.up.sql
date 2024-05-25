CREATE TABLE extension_d20_e35_pc (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL REFERENCES players(id),
    level_total int NOT NULL,
    hit_points int NOT NULL,
    armor_class int NOT NULL,
    class JSONB,
    multiclass BOOLEAN NOT NULL DEFAULT FALSE,
    race text NOT NULL,
    size text NOT NULL,
    weapon JSONB
);

CREATE TABLE extension_d10_homemade_pc (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL REFERENCES players(id),
    health int NOT NULL,
    defense int NOT NULL,
    willpower int NOT NULL,
    initiative int NOT NULL,
    size int NOT NULL,
    armor int NOT NULL,
    weapon JSONB
);

CREATE TABLE extension_d10_oldschool_pc (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL REFERENCES players(id),
    health int NOT NULL,
    willpower int NOT NULL,
    initiative int NOT NULL,
    size int NOT NULL,
    armor int NOT NULL,
    conscience_conviction int NOT NULL,
    self_control_instinct int NOT NULL,
    courage int NOT NULL,
    weapon JSONB
);

CREATE TABLE extension_d20_e35_npc (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL REFERENCES non_players(id),
    level_total int NOT NULL,
    hit_points int NOT NULL,
    armor_class int NOT NULL,
    class JSONB,
    multiclass BOOLEAN NOT NULL DEFAULT FALSE,
    race text NOT NULL,
    size text NOT NULL,
    weapon JSONB
);

CREATE TABLE extension_d10_homemade_npc (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL REFERENCES non_players(id),
    health int NOT NULL,
    defense int NOT NULL,
    willpower int NOT NULL,
    initiative int NOT NULL,
    size int NOT NULL,
    armor int NOT NULL,
    weapon JSONB
);

CREATE TABLE extension_d10_oldschool_npc (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL REFERENCES non_players(id),
    health int NOT NULL,
    willpower int NOT NULL,
    initiative int NOT NULL,
    size int NOT NULL,
    armor int NOT NULL,
    conscience_conviction int NOT NULL,
    self_control_instinct int NOT NULL,
    courage int NOT NULL,
    weapon JSONB
);