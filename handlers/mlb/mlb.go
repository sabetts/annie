package mlb

import (
	"encoding/json"
	"fmt"
	"goirc/bot"
	"goirc/fetch"
	"sort"
	"strings"
	"time"
)

type TeamEndData struct {
	PoffTitle float64
	WsWin     float64
	CsWin     float64
}

type Team struct {
	TeamID   int `json:"teamId"`
	AbbName  string
	League   string
	Division string
	EndData  TeamEndData
}

type TeamList []Team

func (tl TeamList) String() string {
	arr := []string{}
	for i, team := range tl {
		arr = append(arr, fmt.Sprintf("%s:%.0f%%", team.AbbName, 100*team.EndData.PoffTitle))
		if i == 5 { // 6 teams in each division make the poffs
			arr = append(arr, "|")
		}
	}
	return strings.Join(arr, " ")
}

func fetchTeams() (TeamList, error) {
	date := time.Now().Format(time.DateOnly)

	url := fmt.Sprintf("https://www.fangraphs.com/api/playoff-odds/odds?dateEnd=%s&dateDelta=&projectionMode=2&standingsType=lg", date)

	_, bytes, err := fetch.Get(url, time.Minute)
	if err != nil {
		return nil, err
	}

	var teams TeamList
	err = json.Unmarshal(bytes, &teams)
	if err != nil {
		return nil, err
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].EndData.PoffTitle > teams[j].EndData.PoffTitle
	})

	return teams, nil
}

func fetchLeagueTeams(league string) (TeamList, error) {
	teams, err := fetchTeams()
	if err != nil {
		return nil, err
	}

	var lt TeamList

	for _, t := range teams {
		if t.League == league {
			lt = append(lt, t)
		}
	}

	return lt, nil
}

func PlayoffOdds(params bot.HandlerParams) error {
	teams, err := fetchLeagueTeams("AL")
	if err != nil {
		return err
	}
	al := fmt.Sprintf("AL: %s", teams.String())

	teams, err = fetchLeagueTeams("NL")
	if err != nil {
		return err
	}
	nl := fmt.Sprintf("NL: %s", teams.String())

	params.Privmsgf(params.Target, "%s", al)
	params.Privmsgf(params.Target, "%s", nl)

	return nil
}
