CREATE TABLE IF NOT EXISTS "game_machines" (
	"machine_id" SERIAL PRIMARY KEY,
	"ip_address" TEXT,
	"port" INTEGER,
	"session_key" TEXT
);

CREATE TABLE IF NOT EXISTS "games" (
	"game_id" SERIAL PRIMARY KEY,
	"map_name" TEXT NOT NULL,
	"max_players" INTEGER NOT NULL,
	"is_verified" BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS "active_games" (
	"game_id" SERIAL PRIMARY KEY references games(game_id),
	"machine_id" SERIAL references game_machines(machine_id),
	"port" INTEGER
);

CREATE TABLE IF NOT EXISTS "account_data" (
	"user_id" SERIAL PRIMARY KEY,
	"username" VARCHAR(128) NOT NULL,
	"password" BYTEA NOT NULL,
	"salt" BYTEA NOT NULL,
	"algorithm" VARCHAR(128) NOT NULL,
	"createdon" TIMESTAMP NOT NULL,
	"lastlogin" TIMESTAMP NOT NULL
);
