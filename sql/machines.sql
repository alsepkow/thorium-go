DELETE FROM machines;
DROP TABLE machines_metadata;
DROP TABLE machines;

CREATE TABLE IF NOT EXISTS "machines" (
	"machine_id" SERIAL PRIMARY KEY,
	"remote_address" TEXT,
	"service_listen_port" INTEGER
);

CREATE TABLE IF NOT EXISTS "machines_metadata" (
	"machine_id" INTEGER PRIMARY KEY references machines(machine_id) ON DELETE CASCADE,
	"most_recent_key" TEXT,
	"last_heartbeat" TIMESTAMP,
	"cpu_usage_pct" REAL,
	"network_usage_pct" REAL,
	"player_occupancy_pct" REAL
);

