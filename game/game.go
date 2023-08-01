package game

import (
	"Balloon/cfg"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/wonderivan/logger"
)

const (
	ERR_PARAM  = 1
	ERR_USER   = 2
	ERR_ENERGY = 3
	ERR_NAME   = 4
	ERR_RANK   = 5
	ERR_TIMES  = 6
)

var (
	errMsg = map[int]string{
		ERR_PARAM:  "error param",
		ERR_USER:   "error user_id",
		ERR_ENERGY: "zero energy",
		ERR_NAME:   "name error",
		ERR_RANK:   "rank empty",
		ERR_TIMES:  "times over limit",
	}
)

func makeUserId() string {
	s := fmt.Sprintf("tiny_%d", time.Now().Unix())
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func logRequest(r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	logstr := strings.ReplaceAll(string(dump), "\r", "")
	logs := strings.Split(logstr, "\n")
	logger.Info("request: %s %s", logs[0], logs[len(logs)-1])
}

func replyErr(w http.ResponseWriter, code int) {
	logger.Info(`response: {"code":%d, "data":"%s"}`, code, errMsg[code])
	fmt.Fprintf(w, `{"code":%d, "data":"%s"}`, code, errMsg[code])
}

func reply(w http.ResponseWriter, data string) {
	logger.Info(`response: {"code":0, "data":%s}`, data)
	fmt.Fprintf(w, `{"code":0, "data":%s}`, data)
}

func register(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	userId := makeUserId()
	invitor := r.PostFormValue("invitor")
	addUserData(userId, invitor)
	if invitor != "" {
		user := getUserData(invitor)
		if user != nil {
			user.Energy += 1
			energyChange(invitor, user.Energy, time.Now().Unix())
		}
	}
	reply(w, fmt.Sprintf(`{"user_id":"%s"}`, userId))
}

func changeName(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	userId := r.PostFormValue("user_id")
	name := r.PostFormValue("name")
	user := getUserData(userId)
	if user == nil {
		replyErr(w, ERR_USER)
		return
	}
	if len(name) > 32 || strings.ContainsAny(name, " ,./<>?;'\":~!@#$%^&*()_+=-`") {
		replyErr(w, ERR_NAME)
		return
	}
	user.Name = name
	userRename(userId, name)
	reply(w, `"ok"`)
}

func playerInfo(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	userId := r.PostFormValue("user_id")
	user := getUserData(userId)
	if user == nil {
		replyErr(w, ERR_USER)
		return
	}
	// // 恢复体力
	// if recoverEnery(user) {
	// 	energyChange(userId, user.Energy, user.updateTime)
	// }
	reply(w, getSelfJson(user))
}

func startGame(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	userId := r.PostFormValue("user_id")
	user := getUserData(userId)
	if user == nil {
		replyErr(w, ERR_USER)
		return
	}
	// 先恢复体力
	//recoverEnery(user)

	// 再判断体力是否足够
	if user.Energy == 0 {
		replyErr(w, ERR_ENERGY)
		return
	}
	if user.Energy == MAX_ENERGY {
		user.updateTime = getUpdateTime()
	}
	user.Energy -= 1
	energyChange(userId, user.Energy, user.updateTime)

	invitor := getUserData(user.invitor)
	if invitor != nil {
		if invitor.otherTimes < 3 {
			invitor.otherTimes += 3
			invitor.Energy += 1
			addEnergyByWay(invitor.UserId, invitor.Energy, invitor.otherTimes, "other_times")
		}
	}
	explores := []int{}
	explore := 100
	for i := 0; i < 3; i++ {
		explore = getNextExplode(explore, i+1)
		explores = append(explores, explore)
	}
	data, _ := json.Marshal(explores)
	reply(w, string(data))
}

func pushScore(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	userId := r.PostFormValue("user_id")
	scoreStr := r.PostFormValue("score")
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		replyErr(w, ERR_PARAM)
		return
	}
	user := getUserData(userId)
	if user == nil {
		replyErr(w, ERR_USER)
		return
	}

	rank := changeUserScore(userId, score)
	if rank > 0 {
		reply(w, fmt.Sprintf(`{"rank":%d,"score":%d,"total":%d}`, rank, user.Score, len(rankDatas)))
	} else {
		reply(w, `"ok"`)
	}
}

func rankInfo(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	userId := r.PostFormValue("user_id")
	user := getUserData(userId)
	if user == nil {
		replyErr(w, ERR_USER)
		return
	}

	if len(rankDatas) == 0 {
		replyErr(w, ERR_RANK)
		return
	}

	rank := getRank(user)
	count := len(rankDatas)
	if count > 100 {
		count = 100
	}
	rankInfo, _ := json.Marshal(rankDatas[:count])
	reply(w, fmt.Sprintf(`{"rank_info":%s, "self":%d, "total":%d}`, string(rankInfo), rank, len(rankDatas)))
}

func addEnergy(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	userId := r.PostFormValue("user_id")
	user := getUserData(userId)
	if user == nil {
		replyErr(w, ERR_USER)
		return
	}

	addType := r.PostFormValue("type")
	if addType == "ad" {
		if user.adTimes >= 3 {
			replyErr(w, ERR_TIMES)
			return
		}
		user.adTimes += 1
		user.Energy += 1
		addEnergyByWay(userId, user.Energy, user.adTimes, "ad_times")
		reply(w, fmt.Sprintf(`{"energy":%d, "ad_times":%d}`, user.Energy, user.adTimes))
	} else if addType == "invite" {
		if user.inviteTimes >= 1 {
			replyErr(w, ERR_TIMES)
			return
		}
		user.inviteTimes += 1
		user.Energy += 1
		addEnergyByWay(userId, user.Energy, user.inviteTimes, "invite_times")
		reply(w, fmt.Sprintf(`{"energy":%d, "invite_times":%d}`, user.Energy, user.inviteTimes))
	} else {
		replyErr(w, ERR_PARAM)
	}
}

func doReset() {
	rankDatas = []*UserInfo{}
	for _, data := range userDatas {
		data.Score = 0
		data.scoreTime = 0
		data.adTimes = 0
		data.inviteTimes = 0
		data.otherTimes = 0
		if data.Energy < MAX_ENERGY {
			data.Energy = MAX_ENERGY
		}
	}
	resetScores()
	setResetTime()
}

func reset() {
	todayStr := time.Now().Format("2006-01-02")
	today, _ := time.ParseInLocation("2006-01-02", todayStr, time.Local)
	nextDay := today.AddDate(0, 0, 1)

	resetTime := getResetTime()
	if resetTime < today.Unix() {
		doReset()
	}

	waitDuration := time.Until(nextDay)
	time.Sleep(waitDuration)
	for {
		doReset()
		time.Sleep(time.Hour * 24)
	}
}

func Init(cfg *cfg.ConfigInfo) bool {
	if !connect(cfg) {
		return false
	}

	userDatas = loadUserData()
	if userDatas == nil {
		return false
	}
	initRank()
	go reset()

	http.HandleFunc("/token", register)
	http.HandleFunc("/name", changeName)
	http.HandleFunc("/info", playerInfo)
	http.HandleFunc("/start", startGame)
	http.HandleFunc("/score", pushScore)
	http.HandleFunc("/rank", rankInfo)
	http.HandleFunc("/energy", addEnergy)
	http.ListenAndServe("0.0.0.0:"+cfg.Server.Port, nil)
	return true
}
