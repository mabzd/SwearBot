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

type StatsData struct {
	Months map[string]*Month
}

type Month struct {
	Year  int
	Month int
	Users []*User
}

type User struct {
	Name       string
	SwearCount int
}

type BySwearCount []*User

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

func (st *Stats) GetMonthlyRank(month int, year int) ([]*User, error) {
	data, err := st.readStats()
	if err != nil {
		return nil, err
	}
	return getMonthlyRank(data, month, year), nil
}

func (st *Stats) createStatsFileIfNotExist() error {
	if _, err := os.Stat(st.statsFileName); os.IsNotExist(err) {
		data := &StatsData{
			Months: make(map[string]*Month),
		}
		return st.writeStats(data)
	}
	return nil
}

func (st *Stats) readStats() (*StatsData, error) {
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
	var data StatsData
	unmarshalErr := json.Unmarshal(bytes, &data)
	if unmarshalErr != nil {
		log.Printf("Stats: Error when unmarshaling stats from JSON: %s\n", unmarshalErr)
		return nil, errors.New(st.config.OnStatsUnmarshalErr)
	}
	return &data, nil
}

func (st *Stats) writeStats(data *StatsData) error {
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

func addSwearCount(data *StatsData, m int, y int, name string, count int) {
	monthKey := getMonthKey(m, y)
	month := data.Months[monthKey]
	if month == nil {
		month = &Month{
			Year:  y,
			Month: m,
			Users: []*User{},
		}
		data.Months[monthKey] = month
	}
	user := findUser(month.Users, name)
	if user == nil {
		user = &User{
			Name:       name,
			SwearCount: 0,
		}
		month.Users = append(month.Users, user)
	}
	user.SwearCount += count
}

func getMonthlyRank(data *StatsData, m int, y int) []*User {
	monthKey := getMonthKey(m, y)
	month := data.Months[monthKey]
	if month == nil {
		return []*User{}
	}
	sort.Sort(BySwearCount(month.Users))
	return month.Users
}

func findUser(users []*User, name string) *User {
	for _, user := range users {
		if user.Name == name {
			return user
		}
	}
	return nil
}

func getMonthKey(month int, year int) string {
	return fmt.Sprintf("%d.%d", month, year)
}
