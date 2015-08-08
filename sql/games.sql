
DROP TABLE "game_servers";
DROP TABLE "games";


CREATE TABLE IF NOT EXISTS "games" (
	"game_id" SERIAL PRIMARY KEY,
	"map_name" TEXT NOT NULL,
	"game_mode" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "game_servers" (
	"game_id" SERIAL PRIMARY KEY references games(game_id),
	"machine_id" SERIAL references machines(machine_id),
	"port" INTEGER
);

CREATE TABLE IF NOT EXISTS "account_data" (
	"user_id" SERIAL PRIMARY KEY,
	"username" TEXT NOT NULL,
	"password" BYTEA NOT NULL,
	"salt" BYTEA NOT NULL,
	"algorithm" TEXT NOT NULL,
	"createdon" TIMESTAMP NOT NULL,
	"lastlogin" TIMESTAMP NOT NULL
);
