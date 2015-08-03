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

CREATE OR REPLACE FUNCTION get_available_machine()
	RETURNS TABLE (
		"remote_address" TEXT,
		"service_listen_port" INTEGER,
		"most_recent_key" TEXT
	) AS
$$
BEGIN
	RETURN QUERY
	SELECT m.remote_address, m.service_listen_port, mm.most_recent_key
	FROM machines m
	  JOIN machines_metadata mm USING (machine_id)
	WHERE mm.cpu_usage_pct < 80.0
	AND mm.network_usage_pct < 80.0
	ORDER BY RANDOM()
	LIMIT 1;
END
$$ language plpgsql;
