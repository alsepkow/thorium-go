

CREATE TABLE IF NOT EXISTS "games" (
	"game_id" SERIAL PRIMARY KEY,
	"ip_endpoint" TEXT DEFAULT NULL,
	"map_name" TEXT NOT NULL,
	"max_players" INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS "game_machines" (
	"machine_id" SERIAL PRIMARY KEY,
	"ip_endpoint" TEXT,
	"session_key" TEXT
);
