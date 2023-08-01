package game

import (
	"Balloon/cfg"
	"database/sql"
	"fmt"
	"time"

	"Balloon/db"
)

const (
	MAX_ENERGY = 5
)

func addUser(user *UserInfo) {
	now := time.Now().Unix()
	seq := fmt.Sprintf("insert into user(user_id, name, score, energy, update_time, score_time, create_time, ad_times, invite_times, other_times, invitor) values('%s','%s',0,%d,%d,0,%d,%d,%d,%d,'%s')", user.UserId, user.Name, user.Energy, user.updateTime, now, user.adTimes, user.inviteTimes, user.otherTimes, user.invitor)
	db.ExecSql(seq)
}

func userRename(userId string, name string) {
	seq := fmt.Sprintf("update user set name='%s' where user_id='%s'", name, userId)
	db.ExecSql(seq)
}

func energyChange(userId string, energy int, uptime int64) {
	seq := fmt.Sprintf("update user set energy=%d, update_time=%d where user_id='%s'", energy, uptime, userId)
	db.ExecSql(seq)
}

func addEnergyByWay(userId string, energy int, times int, way string) {
	seq := fmt.Sprintf("update user set energy=%d, %s=%d where user_id='%s'", energy, way, times, userId)
	db.ExecSql(seq)
}

func saveScore(userId string, score int, scoretime int64) {
	seq := fmt.Sprintf("update user set score=%d, score_time=%d where user_id='%s'", score, scoretime, userId)
	db.ExecSql(seq)
}

func resetScores() {
	db.ExecSql("update user set score=0, score_time=0, ad_times=0, invite_times=0, other_times=0")
	db.ExecSql(fmt.Sprintf("update user set energy= %d where energy < %d", MAX_ENERGY, MAX_ENERGY))
}

func loadUserData() map[string]*UserInfo {
	datas := map[string]*UserInfo{}
	rows := db.Query("select user_id, name, score, energy, update_time, score_time, ad_times, invite_times, other_times, invitor from user")
	if rows == nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var user UserInfo
		rows.Scan(&user.UserId, &user.Name, &user.Score, &user.Energy, &user.updateTime, &user.scoreTime, &user.adTimes, &user.inviteTimes, &user.otherTimes, &user.invitor)
		datas[user.UserId] = &user
	}
	return datas
}

func getResetTime() int64 {
	row := db.QueryOne("select reset_time from server")

	var t sql.NullInt64
	err := row.Scan(&t)
	if err == sql.ErrNoRows {
		db.ExecSql("insert into server(reset_time) values(0)")
		return 0
	} else if err != nil {
		return -1
	} else {
		return t.Int64
	}
}

func setResetTime() {
	db.ExecSql(fmt.Sprintf("update server set reset_time=%d", time.Now().Unix()))
}

func connect(cfg *cfg.ConfigInfo) bool {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Db)
	return db.InitOpenConns(dbInfo, cfg.Db.MaxConnection)
}
