package thordb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=thoriumnet password=thoriumtest dbname=thoriumnet host=localhost")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func CheckExists(gameId int) (bool, error) {
	res, err := db.Exec("SELECT * FROM public.games WHERE game_id = $1", gameId)
	if err != nil {
		return false, err
	}

	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	exists := rows > 0
	return exists, err
}

func RegisterNewGame(mapName string, maxPlayers int) (int, error) {
	var gameId int
	err := db.QueryRow("INSERT INTO games (map_name, max_players) VALUES ( $1, $2 ) RETURNING game_id", mapName, maxPlayers).Scan(&gameId)
	return gameId, err
}

func RegisterActiveGame(gameId int, machineId int, port int) (bool, error) {
	res, err := db.Exec("INSERT INTO active_games (game_id, machine_id, port) VALUES ( $1, $2, $3 )", gameId, machineId, port)
	if err != nil {
		return false, err
	}

	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	exists := rows > 0
	return exists, err
}

func RegisterMachine(machineIp string, port int) (int, error) {
	var machineId int
	err := db.QueryRow("INSERT INTO game_machines (ip_address, port) VALUES ( $1, $2 ) RETURNING machine_id", machineIp, port).Scan(&machineId)
	return machineId, err
}

func UnregisterMachine(machineId int) (bool, error) {

	success := false

	res, err := db.Exec("DELETE FROM game_machines WHERE machine_id = $1", machineId)
	if err != nil {
		return success, err
	}

	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		return success, err
	}

	success = rows > 0
	return success, err

}

func CheckUsernameExists(username string) (bool, error) {
	fmt.Println("Checking username")
	var exists bool
	rows, err := db.Query("select username FROM account_data;")
	if err != nil {
		fmt.Println("error selecting username: ", err)
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Println("in rows.Next()")
		var tmpuser string
		err := rows.Scan(&tmpuser)
		if len(tmpuser) == 0 {
			return false, err
		}
		if err != nil {
			fmt.Println("error scanning row ", err)
		}
		if username == tmpuser {
			exists = true
			break
		}
	}
	return exists, err

}

func RegisterAccount(username string, hashpw []byte, salt []byte, alg string, created time.Time, lastlogin time.Time) (int, error) {
	var userId int
	err := db.QueryRow("INSERT INTO account_data (username, password, salt, algorithm, createdon, lastlogin) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id", username, hashpw, salt, alg, created, lastlogin).Scan(&userId)
	if err != nil {
		fmt.Println("error inserting account data: ", err)
	}
	return userId, err
}
