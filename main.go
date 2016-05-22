package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/pg.v4"
)

var db *pg.DB
var urlPrefix string
var championToRoleMapping = map[string][]string{}

type IndexResult struct {
	Prefix string
	Error  error
}

type ChampionRoleMapping struct {
	Champion string
	Role     string
}

type MatchupResults struct {
	Enemy         string
	Prefix        string
	Role          string
	SummonerName  string
	Matchups      []ChampionMatchup
	ResultsLength int
}

type MatchupPostData struct {
	SummonerName         string
	Enemy                string
	Role                 string
	Region               string
	RememberMe           bool
	Referral             string
	OriginalSummonerName string
}

func (m *MatchupPostData) FormatFields() {
	m.OriginalSummonerName = m.SummonerName
	m.SummonerName = NormalizeSummonerName(strings.TrimSpace(m.SummonerName))[0]
	m.Region = strings.ToLower(m.Region)
	m.Enemy = NormalizeChampion(m.Enemy)
}

func main() {

	// log.SetFlags(log.Lshortfile)
	log.SetOutput(ioutil.Discard)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	username := viper.GetString("postgres.username")
	password := viper.GetString("postgres.password")
	urlPrefix = viper.GetString("urlPrefix")

	db = pg.Connect(&pg.Options{
		User:     username,
		Password: password,
	})

	err = createSchema(db)
	if err != nil {
		panic(err)
	}

	err = setupRiotApi()
	if err != nil {
		panic(err)
	}

	championToRoleMapping = make(map[string][]string)
	sql := `SELECT champion, role FROM champion_matchups GROUP BY role, champion`
	var championToRole []ChampionRoleMapping
	_, err = db.Query(&championToRole, sql)
	if err != nil {
		fmt.Println("Failed to get champions and their roles")
		panic(err)
	}
	for _, c := range championToRole {
		if championToRoleMapping[c.Champion] == nil {
			championToRoleMapping[c.Champion] = make([]string, 0)
		}
		championToRoleMapping[c.Champion] = append(championToRoleMapping[c.Champion], strings.ToLower(c.Role))
	}

	serveSingle(urlPrefix+"sitemap.xml", "./static/sitemap.xml")
	serveSingle(urlPrefix+"favicon.ico", "./static/favicon.ico")
	serveSingle(urlPrefix+"robots.txt", "./static/robots.txt")
	serveSingle(urlPrefix+"humans.txt", "./static/humans.txt")

	http.HandleFunc(urlPrefix+"", Index)
	http.HandleFunc(urlPrefix+"matchup", GetMatchup)
	fs := http.FileServer(http.Dir("static"))
	http.Handle(urlPrefix+"static/", http.StripPrefix(urlPrefix+"static/", fs))
	fmt.Println("Server started")
	// analytics
	http.HandleFunc(urlPrefix+"analytics/index", AnalyzeIndex)
	http.HandleFunc(urlPrefix+"analytics/matchup", AnalyzeMatchup)
	http.HandleFunc(urlPrefix+"analytics/external", AnalyzeExternalLink)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		defer un(trace("/ "))
		result := IndexResult{urlPrefix, nil}
		t, _ := template.ParseFiles("index.html")
		a := AnalyticsPage{"/", r.Referer(), time.Now().UTC()}
		err := db.Create(&a)
		if err != nil {
			fmt.Println("Failed to save page analyics", err)
		}
		t.Execute(w, result)
	} else {
		http.Redirect(w, r, "/", 301)
	}
}

func GetMatchup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer un(trace("/matchup"))

		decoder := json.NewDecoder(r.Body)
		var m MatchupPostData
		err := decoder.Decode(&m)
		if err != nil {
			fmt.Println("Failed to decode matchup data", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.FormatFields()

		invalidKey := validateMatchupPayload(m)
		if invalidKey != "" {
			http.Error(w, "Enter a valid "+invalidKey, http.StatusBadRequest)
			return
		}

		if CHAMPION_KEYS[m.Enemy] == "" {
			http.Error(w, "Enter a valid Enemy Champion", http.StatusBadRequest)
			return
		}

		a := AnalyticsPage{"/matchup", r.Referer(), time.Now().UTC()}
		err = db.Create(&a)
		if err != nil {
			fmt.Println("Failed to save index analyics", err)
		}

		summoner, err := getOrCreateSummoner(m.Region, m.SummonerName, db)
		if err != nil {
			errors.New("Summoner does not exist")
			if err.Error() == "Summoner does not exist" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			if err.Error() == "Failed to get summoner" || err.Error() == "Failed to get champion masteries" {

				http.Error(w, err.Error()+" from riot's api", http.StatusInternalServerError)
				return
			}
			// t, _ := template.ParseFiles("index.html")
			fmt.Println("Error get or create summoner:", m.SummonerName, err)
			// result := IndexResult{urlPrefix, err}
			// t.Execute(w, result)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		matchups, pms, err := getMatchups(summoner, CHAMPION_KEYS[m.Enemy], m.Role, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for index, matchup := range matchups {
			matchups[index].SetChampion(CHAMPION_KEYS_BY_KEY_PROPER_CASING[matchup.Champion])
			for _, p := range pms {
				if p.Champion == matchup.Champion {
					matchups[index].UpdatePersonalData(p)
				}
			}
		}
		// fmt.Println(CHAMPION_KEYS_BY_KEY_PROPER_CASING[CHAMPION_KEYS[enemy]], urlPrefix, m.Role, summonerName, matchups)
		result := MatchupResults{CHAMPION_KEYS_BY_KEY_PROPER_CASING[CHAMPION_KEYS[m.Enemy]], urlPrefix, m.Role, m.OriginalSummonerName, matchups, len(matchups)}
		t, _ := template.ParseFiles("matchups.html")
		t.Execute(w, result)
		// json.NewEncoder(w).Encode(result)
	} else {
		// StatusNotFound
		http.NotFound(w, r)
	}
}

func validateMatchupPayload(m MatchupPostData) string {
	if m.SummonerName == "" {
		return "Summoner Name"
	} else if m.Enemy == "" {
		return "Enemy Champion"
	} else if m.Role == "" {
		return "Role"
	} else if !RIOT_REGIONS[m.Region] {
		return "Region"
	}
	return ""
}
