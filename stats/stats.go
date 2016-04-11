package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type Stats struct {
	statsFileName string
	config        StatsConfig
}

type StatsConfig struct {
	OnStatsFileCreateErr string
	OnStatsFileReadErr   string
	OnStatsUnmarshalErr  string
	OnStatsMarshalErr    string
	OnStatsSaveErr       string
}

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

func NewStats(statsFileName string, config StatsConfig) *Stats {
	return &Stats{
		statsFileName: statsFileName,
		config:        config,
	}
}

func (st *Stats) AddSwearCount(month int, year int, name string, count int) error {
	data, err := st.readStats()
	if err != nil {
		return err
	}
	addSwearCount(data, month, year, name, count)
	return st.writeStats(data)
}

func (st *Stats) GetMonthlyRank(month int, year int) ([]*UserStats, error) {
	data, err := st.readStats()
	if err != nil {
		return nil, err
	}
	return getMonthlyRank(data, month, year), nil
}

func (st *Stats) createStatsFileIfNotExist() error {
	if _, err := os.Stat(st.statsFileName); os.IsNotExist(err) {
		data := &AllStats{
			Months: make(map[string]*MonthStats),
		}
		return st.writeStats(data)
	}
	return nil
}

func (st *Stats) readStats() (*AllStats, error) {
	createErr := st.createStatsFileIfNotExist()
	if createErr != nil {
		log.Printf("Stats: Cannot create stats file '%s': %s\n", st.statsFileName, createErr)
		return nil, errors.New(st.config.OnStatsFileCreateErr)
	}
	bytes, fileReadErr := ioutil.ReadFile(st.statsFileName)
	if fileReadErr != nil {
		log.Printf("Stats: Cannot read stats from file '%s': %s\n", st.statsFileName, fileReadErr)
		return nil, errors.New(st.config.OnStatsFileReadErr)
	}
	var data AllStats
	unmarshalErr := json.Unmarshal(bytes, &data)
	if unmarshalErr != nil {
		log.Printf("Stats: Error when unmarshaling stats from JSON: %s\n", unmarshalErr)
		return nil, errors.New(st.config.OnStatsUnmarshalErr)
	}
	return &data, nil
}

func (st *Stats) writeStats(data *AllStats) error {
	bytes, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		log.Printf("Stats: Error when marshaling stats to JSON: %s\n", marshalErr)
		return errors.New(st.config.OnStatsMarshalErr)
	}
	saveErr := ioutil.WriteFile(st.statsFileName, bytes, 0666)
	if saveErr != nil {
		log.Printf("Stats: Cannot write stats to file '%s': %s\n", saveErr)
		return errors.New(st.config.OnStatsSaveErr)
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
	user := findUser(monthStats.Users, userId)
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

func findUser(users []*UserStats, userId string) *UserStats {
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
