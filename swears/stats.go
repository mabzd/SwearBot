package swears

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
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

func (sw *Swears) AddSwearCount(month int, year int, name string, count int) error {
	data, err := sw.readStats()
	if err != nil {
		return err
	}
	addSwearCount(data, month, year, name, count)
	return sw.writeStats(data)
}

func (sw *Swears) GetMonthlyRank(month int, year int) ([]*UserStats, error) {
	data, err := sw.readStats()
	if err != nil {
		return nil, err
	}
	return getMonthlyRank(data, month, year), nil
}

func (sw *Swears) createStatsFileIfNotExist() error {
	if _, err := os.Stat(sw.config.StatsFileName); os.IsNotExist(err) {
		data := &AllStats{
			Months: make(map[string]*MonthStats),
		}
		return sw.writeStats(data)
	}
	return nil
}

func (sw *Swears) readStats() (*AllStats, error) {
	fileName := sw.config.StatsFileName
	createErr := sw.createStatsFileIfNotExist()
	if createErr != nil {
		log.Printf("Swears: Cannot create Swears file '%s': %s\n", fileName, createErr)
		return nil, errors.New(sw.config.OnStatsFileCreateErr)
	}
	bytes, fileReadErr := ioutil.ReadFile(fileName)
	if fileReadErr != nil {
		log.Printf("Swears: Cannot read Swears from file '%s': %s\n", fileName, fileReadErr)
		return nil, errors.New(sw.config.OnStatsFileReadErr)
	}
	var data AllStats
	unmarshalErr := json.Unmarshal(bytes, &data)
	if unmarshalErr != nil {
		log.Printf("Swears: Error when unmarshaling Swears from JSON: %s\n", unmarshalErr)
		return nil, errors.New(sw.config.OnStatsUnmarshalErr)
	}
	return &data, nil
}

func (sw *Swears) writeStats(data *AllStats) error {
	fileName := sw.config.StatsFileName
	bytes, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		log.Printf("Swears: Error when marshaling Swears to JSON: %s\n", marshalErr)
		return errors.New(sw.config.OnStatsMarshalErr)
	}
	saveErr := ioutil.WriteFile(fileName, bytes, 0666)
	if saveErr != nil {
		log.Printf("Swears: Cannot write Swears to file '%s': %s\n", fileName, saveErr)
		return errors.New(sw.config.OnStatsSaveErr)
	}
	return nil
}

func addSwearCount(data *AllStats, month int, year int, userId string, count int) {
	monthKey := getMonthKey(month, year)
	monthStats := data.Months[monthKey]
	if monthStats == nil {
		monthStats = &MonthStats{
			Year:  year,
			Month: month,
			Users: []*UserStats{},
		}
		data.Months[monthKey] = monthStats
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

func getMonthlyRank(data *AllStats, month int, year int) []*UserStats {
	monthKey := getMonthKey(month, year)
	monthStats := data.Months[monthKey]
	if monthStats == nil {
		return []*UserStats{}
	}
	sort.Sort(BySwearCount(monthStats.Users))
	return monthStats.Users
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
