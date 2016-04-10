package swearbot

import (
	"fmt"
	"os"
	"log"
	"sort"
	"errors"
	"io/ioutil"
	"encoding/json"
)

type Stats struct {
	Months map[string]*Month
}

type Month struct {
	Year int
	Month int
	Users []*User
}

type User struct {
	Name string
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

func (sb *SwearBot) AddSwears(month int, year int, name string, swears int) error {
	stats, err := sb.readStats()
	if err != nil {
		return err
	}
	addSwears(stats, month, year, name, swears)
	return sb.writeStats(stats)
}

func (sb *SwearBot) GetMonthlyRank(month int, year int) ([]*User, error) {
	stats, err := sb.readStats()
	if err != nil {
		return nil, err
	}
	return getMonthlyRank(stats, month, year), nil
}

func addSwears(stats *Stats, m int, y int, name string, swears int) {
	monthKey := getMonthKey(m, y)
	month := stats.Months[monthKey]
	if month == nil {
		month = &Month {
			Year: y,
			Month: m,
			Users: []*User{},
		}
		stats.Months[monthKey] = month
	}
	user := findUser(month.Users, name)
	if user == nil {
		user = &User {
			Name: name,
			SwearCount: 0,
		}
		month.Users = append(month.Users, user)
	}
	user.SwearCount += swears
}

func getMonthlyRank(stats *Stats, m int, y int) []*User {
	monthKey := getMonthKey(m, y)
	month := stats.Months[monthKey]
	if month == nil {
		return []*User{}
	}
	sort.Sort(BySwearCount(month.Users))
	return month.Users
}

func (sb *SwearBot) createStatsFileIfNotExist() error {
	if _, err := os.Stat(sb.statsFileName); os.IsNotExist(err) {
		stats := &Stats {
			Months: make(map[string]*Month),
		}
		return sb.writeStats(stats)
	}
	return nil
}

func (sb *SwearBot) readStats() (*Stats, error) {
	createErr := sb.createStatsFileIfNotExist()
	if createErr != nil {
		log.Printf("Stats: Cannot create stats file '%s': %s\n", sb.statsFileName, createErr)
		return nil, errors.New(sb.config.OnStatsFileCreateErr)
	}
	bytes, fileReadErr := ioutil.ReadFile(sb.statsFileName)
	if fileReadErr != nil {
		log.Printf("Stats: Cannot read stats from file '%s': %s\n", sb.statsFileName, fileReadErr)
		return nil, errors.New(sb.config.OnStatsFileReadErr)
	}
	var stats Stats
	unmarshalErr := json.Unmarshal(bytes, &stats)
	if unmarshalErr != nil {
		log.Printf("Stats: Error when unmarshaling stats from JSON: %s\n", unmarshalErr)
		return nil, errors.New(sb.config.OnStatsUnmarshalErr)
	}
	return &stats, nil
}

func (sb *SwearBot) writeStats(stats *Stats) error {
	bytes, marshalErr := json.Marshal(stats)
	if marshalErr != nil {
		log.Printf("Stats: Error when marshaling stats to JSON: %s\n", marshalErr)
		return errors.New(sb.config.OnStatsMarshalErr)
	}
	saveErr := ioutil.WriteFile(sb.statsFileName, bytes, 0666)
	if saveErr != nil {
		log.Printf("Stats: Cannot write stats to file '%s': %s\n", saveErr)
		return errors.New(sb.config.OnStatsSaveErr)
	}
	return nil
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