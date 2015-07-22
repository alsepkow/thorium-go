
DROP TABLE "characters";

CREATE TABLE "characters" (
	"id" SERIAL PRIMARY KEY,
	"uid" INTEGER references account_data,
	"name" TEXT,
	"game_data" JSON
);
