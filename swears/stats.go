package swears

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

const (
	StatsFileCreateErr = 11
	StatsFileReadErr   = 12
	StatsUnmarshalErr  = 13
	StatsMarshalErr    = 14
	StatsSaveErr       = 15
)

type AllStats struct {
	Months map[string]*MonthStats
}

type MonthStats struct {
	Year  int
	Month int
	Users []*UserStats
}

type UserStats struct {
	UserId     string
	SwearCount int
}

type BySwearCount []*UserStats

func (a BySwearCount) Len() int {
	return len(a)
}

func (a BySwearCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a BySwearCount) Less(i, j int) bool {
	return a[i].SwearCount > a[j].SwearCount
}

func (sw *Swears) AddSwearCount(month int, year int, name string, count int) int {
	stats, err := readStats(sw.config.StatsFileName)
	if err != Success {
		return err
	}

	addSwearCount(stats, month, year, name, count)
	return writeStats(sw.config.StatsFileName, stats)
}

func (sw *Swears) GetMonthlyRank(month int, year int) ([]*UserStats, int) {
	stats, err := readStats(sw.config.StatsFileName)
	if err != Success {
		return nil, err
	}

	return getMonthlyRank(stats, month, year), Success
}

func (sw *Swears) GetTotalRank() ([]*UserStats, int) {
	stats, err := readStats(sw.config.StatsFileName)
	if err != Success {
		return nil, err
	}

	return getTotalRank(stats), Success
}

func createStatsFileIfNotExist(fileName string) int {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		stats := &AllStats{
			Months: make(map[string]*MonthStats),
		}
		return writeStats(fileName, stats)
	}

	return Success
}

func readStats(fileName string) (*AllStats, int) {
	createErr := createStatsFileIfNotExist(fileName)
	if createErr != Success {
		log.Println("Swears: Stats file creation failed.")
		return nil, StatsFileCreateErr
	}

	bytes, fileReadErr := ioutil.ReadFile(fileName)
	if fileReadErr != nil {
		log.Printf("Swears: Cannot read Swears from file '%s': %s\n", fileName, fileReadErr)
		return nil, StatsFileReadErr
	}

	var stats AllStats
	unmarshalErr := json.Unmarshal(bytes, &stats)
	if unmarshalErr != nil {
		log.Printf("Swears: Error when unmarshaling Swears from JSON: %s\n", unmarshalErr)
		return nil, StatsUnmarshalErr
	}

	return &stats, Success
}

func writeStats(fileName string, stats *AllStats) int {
	bytes, marshalErr := json.Marshal(stats)
	if marshalErr != nil {
		log.Printf("Swears: Error when marshaling Swears to JSON: %s\n", marshalErr)
		return StatsMarshalErr
	}

	saveErr := ioutil.WriteFile(fileName, bytes, 0666)
	if saveErr != nil {
		log.Printf("Swears: Cannot write Swears to file '%s': %s\n", fileName, saveErr)
		return StatsSaveErr
	}

	return Success
}

func addSwearCount(stats *AllStats, month int, year int, userId string, count int) {
	monthKey := getMonthKey(month, year)
	monthStats := stats.Months[monthKey]
	if monthStats == nil {
		monthStats = &MonthStats{
			Year:  year,
			Month: month,
			Users: []*UserStats{},
		}
		stats.Months[monthKey] = monthStats
	}
	user := getUserStatsById(monthStats.Users, userId)
	if user == nil {
		user = &UserStats{
			UserId:     userId,
			SwearCount: 0,
		}
		monthStats.Users = append(monthStats.Users, user)
	}
	user.SwearCount += count
}

func getMonthlyRank(stats *AllStats, month int, year int) []*UserStats {
	monthKey := getMonthKey(month, year)
	monthStats := stats.Months[monthKey]
	if monthStats == nil {
		return []*UserStats{}
	}
	sort.Sort(BySwearCount(monthStats.Users))
	return monthStats.Users
}

func getTotalRank(stats *AllStats) []*UserStats {
	userIdToSwears := make(map[string]int)
	for _, monthStats := range stats.Months {
		for _, userStats := range monthStats.Users {
			userId := userStats.UserId
			swearCount := userIdToSwears[userId]
			userIdToSwears[userId] = swearCount + userStats.SwearCount
		}
	}
	totalRank := toUserStats(userIdToSwears)
	sort.Sort(BySwearCount(totalRank))
	return totalRank
}

func toUserStats(userIdToSwears map[string]int) []*UserStats {
	userStats := []*UserStats{}
	for userId, swears := range userIdToSwears {
		userStats = append(userStats, &UserStats{UserId: userId, SwearCount: swears})
	}
	return userStats
}

func getUserStatsById(users []*UserStats, userId string) *UserStats {
	for _, user := range users {
		if user.UserId == userId {
			return user
		}
	}
	return nil
}

func getMonthKey(month int, year int) string {
	return fmt.Sprintf("%d.%d", month, year)
}
