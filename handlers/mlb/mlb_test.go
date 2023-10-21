package mlb

import (
	"fmt"
	"testing"
)

func TestFetchTeams(t *testing.T) {
	teams, err := fetchTeams()
	if err != nil {
		t.Error(err)
	}
	t.Log(teams)
}

func TestFetchLeagueTeams(t *testing.T) {
	for _, league := range []string{"AL", "NL"} {
		t.Run(league, func(t *testing.T) {
			teams, err := fetchLeagueTeams(league)
			if err != nil {
				t.Error(err)
			}
			if len(teams) != 15 {
				t.Errorf("expected 15 teams, got %d", len(teams))
			}
		})
	}
}

func TestLastUpdatedAt(t *testing.T) {
	got, err := lastUpdatedAt()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("got %v", got)
}
