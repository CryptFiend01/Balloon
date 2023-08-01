package game

import (
	"encoding/json"
	"math"
	"math/rand"
	"sort"
	"time"
)

type UserInfo struct {
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	Energy      int    `json:"energy"`
	updateTime  int64
	scoreTime   int64
	adTimes     int
	inviteTimes int
	otherTimes  int
	invitor     string
}

type SelfInfo struct {
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	Energy      int    `json:"energy"`
	AdTimes     int    `json:"add_times"`
	InviteTimes int    `json:"invite_times"`
}

var (
	rankDatas []*UserInfo
	userDatas map[string]*UserInfo
)

func sortRank() {
	sort.Slice(rankDatas, func(i, j int) bool {
		if rankDatas[i].Score > rankDatas[j].Score {
			return true
		} else if rankDatas[i].Score == rankDatas[j].Score {
			return rankDatas[i].scoreTime < rankDatas[j].scoreTime
		} else {
			return false
		}
	})
}

func getRank(user *UserInfo) int {
	if len(rankDatas) == 0 || user.Score == 0 {
		return -1
	}
	left := 0
	right := len(rankDatas) - 1
	if rankDatas[left] == user {
		return left + 1
	} else if rankDatas[right] == user {
		return right + 1
	}

	for {
		mid := left + (right-left)/2
		if rankDatas[mid] == user {
			return mid + 1
		}
		if user.Score < rankDatas[mid].Score {
			left = mid
		} else if user.Score > rankDatas[mid].Score {
			right = mid
		} else {
			if rankDatas[mid].scoreTime < user.scoreTime {
				left = mid
			} else {
				right = mid
			}
		}
	}
}

func initRank() {
	for _, data := range userDatas {
		if data.Score > 0 {
			rankDatas = append(rankDatas, data)
		}
	}
	sortRank()
}

func getUpdateTime() int64 {
	now := time.Now()
	return now.Unix() - int64(now.Minute()*60) - int64(now.Second())
}

func addUserData(userId string, invitor string) bool {
	user := getUserData(userId)
	if user != nil {
		return false
	}

	user = &UserInfo{
		UserId:      userId,
		Name:        userId,
		Score:       0,
		Energy:      MAX_ENERGY,
		updateTime:  getUpdateTime(),
		scoreTime:   0,
		adTimes:     0,
		inviteTimes: 0,
		otherTimes:  0,
		invitor:     invitor,
	}
	userDatas[userId] = user
	addUser(user)
	return true
}

func getUserData(userId string) *UserInfo {
	user, ok := userDatas[userId]
	if ok {
		return user
	} else {
		return nil
	}
}

func getSelfJson(user *UserInfo) string {
	self := SelfInfo{
		UserId:      user.UserId,
		Name:        user.Name,
		Score:       user.Score,
		Energy:      user.Energy,
		AdTimes:     user.adTimes,
		InviteTimes: user.inviteTimes,
	}
	data, _ := json.Marshal(self)
	return string(data)
}

func changeUserScore(userId string, score int) int {
	user := getUserData(userId)
	if user == nil {
		return -1
	}
	if score > 0 {
		if user.Score == 0 {
			rankDatas = append(rankDatas, user)
		}
		user.Score += score
		user.scoreTime = time.Now().Unix()
		saveScore(userId, score, user.scoreTime)
		sortRank()
		return getRank(user)
	} else {
		return 0
	}
}

// func recoverEnery(user *UserInfo) bool {
// 	if user.Energy >= MAX_ENERGY {
// 		return false
// 	}
// 	now := time.Now()
// 	uptime := time.Unix(user.updateTime, 0)
// 	hour := int(now.Sub(uptime).Hours())
// 	if hour > 0 {
// 		user.Energy += hour
// 		if user.Energy > MAX_ENERGY {
// 			user.Energy = MAX_ENERGY
// 		}
// 		user.updateTime = getUpdateTime()
// 		return true
// 	} else {
// 		return false
// 	}
// }

func getNextExplode(nextExplode int, times int) int {
	randomValue := rand.Float64()
	increaseValue := math.Floor(100.0/randomValue - 100.0)
	nextExplode += int(increaseValue)
	if nextExplode > 150 && times == 1 {
		nextExplode = 150
	} else if nextExplode > 300 && times == 2 {
		nextExplode = 300
	}
	return nextExplode
}
