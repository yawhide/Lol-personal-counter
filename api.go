package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cznic/sortutil"
	"github.com/spf13/viper"
	"github.com/yawhide/go-lol"
	"gopkg.in/pg.v4"
)

var RIOT_REGIONS = map[string]bool{
	"br":   true,
	"eune": true,
	"euw":  true,
	"jp":   true,
	"kr":   true,
	"lan":  true,
	"las":  true,
	"na":   true,
	"oce":  true,
	"tr":   true,
	"ru":   true,
}
var CHAMPIONS = [...]string{"shaco", "drmundo", "rammus", "anivia", "irelia", "yasuo", "sona", "kassadin", "zac", "gnar", "karma", "corki", "gangplank", "janna", "jhin", "kindred", "braum", "ashe", "tryndamere", "jax", "morgana", "zilean", "singed", "evelynn", "twitch", "galio", "velkoz", "olaf", "annie", "karthus", "leblanc", "urgot", "amumu", "xinzhao", "chogath", "twistedfate", "fiddlesticks", "vladimir", "warwick", "teemo", "tristana", "sivir", "soraka", "ryze", "sion", "masteryi", "alistar", "missfortune", "nunu", "rengar", "volibear", "fizz", "graves", "ahri", "shyvana", "lux", "xerath", "thresh", "shen", "kogmaw", "jinx", "tahmkench", "riven", "talon", "malzahar", "kayle", "kalista", "reksai", "illaoi", "leona", "lulu", "gragas", "poppy", "fiora", "ziggs", "udyr", "viktor", "sejuani", "varus", "nautilus", "draven", "bard", "mordekaiser", "ekko", "yorick", "pantheon", "ezreal", "garen", "akali", "kennen", "vayne", "jayce", "lissandra", "cassiopeia", "rumble", "khazix", "darius", "hecarim", "skarner", "lucian", "heimerdinger", "nasus", "zed", "nidalee", "syndra", "jarvaniv", "quinn", "renekton", "maokai", "aurelionsol", "nocturne", "katarina", "leesin", "monkeyking", "azir", "brand", "diana", "elise", "nami", "aatrox", "orianna", "zyra", "trundle", "veigar", "taric", "caitlyn", "blitzcrank", "malphite", "vi", "swain", "taliyah"}
var CHAMPION_KEYS = map[string]string{"aatrox": "266", "ahri": "103", "akali": "84", "alistar": "12", "amumu": "32", "anivia": "34", "annie": "1", "ashe": "22", "aurelionsol": "136", "azir": "268", "bard": "432", "blitzcrank": "53", "brand": "63", "braum": "201", "caitlyn": "51", "cassiopeia": "69", "chogath": "31", "corki": "42", "darius": "122", "diana": "131", "draven": "119", "drmundo": "36", "ekko": "245", "elise": "60", "evelynn": "28", "ezreal": "81", "fiddlesticks": "9", "fiora": "114", "fizz": "105", "galio": "3", "gangplank": "41", "garen": "86", "gnar": "150", "gragas": "79", "graves": "104", "hecarim": "120", "heimerdinger": "74", "illaoi": "420", "irelia": "39", "janna": "40", "jarvaniv": "59", "jax": "24", "jayce": "126", "jhin": "202", "jinx": "222", "kalista": "429", "karma": "43", "karthus": "30", "kassadin": "38", "katarina": "55", "kayle": "10", "kennen": "85", "khazix": "121", "kindred": "203", "kogmaw": "96", "leblanc": "7", "leesin": "64", "leona": "89", "lissandra": "127", "lucian": "236", "lulu": "117", "lux": "99", "malphite": "54", "malzahar": "90", "maokai": "57", "masteryi": "11", "missfortune": "21", "monkeyking": "62", "wukong": "62", "mordekaiser": "82", "morgana": "25", "nami": "267", "nasus": "75", "nautilus": "111", "nidalee": "76", "nocturne": "56", "nunu": "20", "olaf": "2", "orianna": "61", "pantheon": "80", "poppy": "78", "quinn": "133", "rammus": "33", "reksai": "421", "renekton": "58", "rengar": "107", "riven": "92", "rumble": "68", "ryze": "13", "sejuani": "113", "shaco": "35", "shen": "98", "shyvana": "102", "singed": "27", "sion": "14", "sivir": "15", "skarner": "72", "sona": "37", "soraka": "16", "swain": "50", "syndra": "134", "tahmkench": "223", "taliyah": "163", "talon": "91", "taric": "44", "teemo": "17", "thresh": "412", "tristana": "18", "trundle": "48", "tryndamere": "23", "twistedfate": "4", "twitch": "29", "udyr": "77", "urgot": "6", "varus": "110", "vayne": "67", "veigar": "45", "velkoz": "161", "vi": "254", "viktor": "112", "vladimir": "8", "volibear": "106", "warwick": "19", "xerath": "101", "xinzhao": "5", "yasuo": "157", "yorick": "83", "zac": "154", "zed": "238", "ziggs": "115", "zilean": "26", "zyra": "143"}
var CHAMPION_KEYS_BY_KEY = map[string]string{"1": "annie", "2": "olaf", "3": "galio", "4": "twistedfate", "5": "xinzhao", "6": "urgot", "7": "leblanc", "8": "vladimir", "9": "fiddlesticks", "10": "kayle", "11": "masteryi", "12": "alistar", "13": "ryze", "14": "sion", "15": "sivir", "16": "soraka", "17": "teemo", "18": "tristana", "19": "warwick", "20": "nunu", "21": "missfortune", "22": "ashe", "23": "tryndamere", "24": "jax", "25": "morgana", "26": "zilean", "27": "singed", "28": "evelynn", "29": "twitch", "30": "karthus", "31": "chogath", "32": "amumu", "33": "rammus", "34": "anivia", "35": "shaco", "36": "drmundo", "37": "sona", "38": "kassadin", "39": "irelia", "40": "janna", "41": "gangplank", "42": "corki", "43": "karma", "44": "taric", "45": "veigar", "48": "trundle", "50": "swain", "51": "caitlyn", "53": "blitzcrank", "54": "malphite", "55": "katarina", "56": "nocturne", "57": "maokai", "58": "renekton", "59": "jarvaniv", "60": "elise", "61": "orianna", "62": "monkeyking", "63": "brand", "64": "leesin", "67": "vayne", "68": "rumble", "69": "cassiopeia", "72": "skarner", "74": "heimerdinger", "75": "nasus", "76": "nidalee", "77": "udyr", "78": "poppy", "79": "gragas", "80": "pantheon", "81": "ezreal", "82": "mordekaiser", "83": "yorick", "84": "akali", "85": "kennen", "86": "garen", "89": "leona", "90": "malzahar", "91": "talon", "92": "riven", "96": "kogmaw", "98": "shen", "99": "lux", "101": "xerath", "102": "shyvana", "103": "ahri", "104": "graves", "105": "fizz", "106": "volibear", "107": "rengar", "110": "varus", "111": "nautilus", "112": "viktor", "113": "sejuani", "114": "fiora", "115": "ziggs", "117": "lulu", "119": "draven", "120": "hecarim", "121": "khazix", "122": "darius", "126": "jayce", "127": "lissandra", "131": "diana", "133": "quinn", "134": "syndra", "136": "aurelionsol", "143": "zyra", "150": "gnar", "154": "zac", "157": "yasuo", "161": "velkoz", "163": "taliyah", "201": "braum", "202": "jhin", "203": "kindred", "222": "jinx", "223": "tahmkench", "236": "lucian", "238": "zed", "245": "ekko", "254": "vi", "266": "aatrox", "267": "nami", "268": "azir", "412": "thresh", "420": "illaoi", "421": "reksai", "429": "kalista", "432": "bard"}
var CHAMPION_KEYS_BY_KEY_PROPER_CASING = map[string]string{"35": "Shaco", "36": "DrMundo", "33": "Rammus", "34": "Anivia", "39": "Irelia", "157": "Yasuo", "37": "Sona", "38": "Kassadin", "154": "Zac", "150": "Gnar", "43": "Karma", "42": "Corki", "41": "Gangplank", "40": "Janna", "202": "Jhin", "203": "Kindred", "201": "Braum", "22": "Ashe", "23": "Tryndamere", "24": "Jax", "25": "Morgana", "26": "Zilean", "27": "Singed", "28": "Evelynn", "29": "Twitch", "3": "Galio", "161": "Velkoz", "2": "Olaf", "1": "Annie", "30": "Karthus", "7": "Leblanc", "6": "Urgot", "32": "Amumu", "5": "XinZhao", "31": "Chogath", "4": "TwistedFate", "9": "FiddleSticks", "8": "Vladimir", "19": "Warwick", "17": "Teemo", "18": "Tristana", "15": "Sivir", "16": "Soraka", "13": "Ryze", "14": "Sion", "11": "MasterYi", "12": "Alistar", "21": "MissFortune", "20": "Nunu", "107": "Rengar", "106": "Volibear", "105": "Fizz", "104": "Graves", "103": "Ahri", "102": "Shyvana", "99": "Lux", "101": "Xerath", "412": "Thresh", "98": "Shen", "96": "KogMaw", "222": "Jinx", "223": "TahmKench", "92": "Riven", "91": "Talon", "90": "Malzahar", "10": "Kayle", "429": "Kalista", "421": "RekSai", "420": "Illaoi", "89": "Leona", "117": "Lulu", "79": "Gragas", "78": "Poppy", "114": "Fiora", "115": "Ziggs", "77": "Udyr", "112": "Viktor", "113": "Sejuani", "110": "Varus", "111": "Nautilus", "119": "Draven", "432": "Bard", "82": "Mordekaiser", "245": "Ekko", "83": "Yorick", "80": "Pantheon", "81": "Ezreal", "86": "Garen", "84": "Akali", "85": "Kennen", "67": "Vayne", "126": "Jayce", "127": "Lissandra", "69": "Cassiopeia", "68": "Rumble", "121": "Khazix", "122": "Darius", "120": "Hecarim", "72": "Skarner", "236": "Lucian", "74": "Heimerdinger", "75": "Nasus", "238": "Zed", "76": "Nidalee", "134": "Syndra", "59": "JarvanIV", "133": "Quinn", "58": "Renekton", "57": "Maokai", "136": "AurelionSol", "56": "Nocturne", "55": "Katarina", "64": "LeeSin", "62": "MonkeyKing", "268": "Azir", "63": "Brand", "131": "Diana", "60": "Elise", "267": "Nami", "266": "Aatrox", "61": "Orianna", "143": "Zyra", "48": "Trundle", "45": "Veigar", "44": "Taric", "51": "Caitlyn", "53": "Blitzcrank", "54": "Malphite", "254": "Vi", "50": "Swain", "163": "Taliyah"}

var apiEndpointMap map[string]*lol.APIEndpoint

type Mastery struct {
	ChampionID                   int    `json:"championId", sql:",pk"`
	ChampionLevel                int    `json:"championLevel"`
	ChampionPoints               int    `json:"championPoints"`
	ChampionPointsSinceLastLevel int    `json:"championPointsSinceLastLevel"`
	ChampionPointsUntilNextLevel int    `json:"championPointsUntilNextLevel"`
	LastPlayTime                 uint64 `json:"lastPlayTime"`
	Region                       string `sql:",pk"`
	SummonerID                   uint64 `json:"playerId", sql:",pk"`
}

type RiotError struct {
	StatusCode int
}

func (err RiotError) Error() string {
	return fmt.Sprintf("Error: HTTP Status %d", err.StatusCode)
}

type MySummoner struct {
	Name               string `json:"name"`
	ProfileIconID      int    `json:"profileIconId"`
	MasteriesUpdatedAt time.Time
	MatchlistTimestamp uint64
	MatchlistUpdatedAt time.Time
	Region             string `sql:",pk"`
	RevisionDate       uint64 `json:"revisionDate"`
	SummonerID         uint64 `json:"id", sql:",pk"`
	SummonerLevel      uint32 `json:"summonerLevel"`
}

func (s *MySummoner) SetSummoner(name string) {
	s.Name = name
}

func (s *MySummoner) UpdateMatchlistUpdatedAt() {
	s.MatchlistUpdatedAt = time.Now().UTC()
}

func (s *MySummoner) UpdateMatchlistTimestamp(t uint64) {
	s.MatchlistTimestamp = t
	s.MatchlistUpdatedAt = time.Now().UTC()
}

type ChampionMatchupWinrate struct {
	Role     string `json:"role"`
	Matchups []Matchup
}

type Matchup struct {
	Games     int     `json:"games"`
	StatScore float32 `json:"statScore"`
	WinRate   float32 `json:"winRate"`
	Enemy     string  `json:"key"`
}

type ChampionMatchup struct {
	Champion        string `sql:",pk"`
	Enemy           string `sql:",pk"`
	Games           int
	Role            string `sql:",pk"`
	StatScore       float32
	WinRate         float32
	PersonalWinRate string
	PersonalGames   int
}

func (c *ChampionMatchup) SetChampion(name string) {
	c.Champion = name
}

func (m ChampionMatchup) String() string {
	return fmt.Sprintf("[%s] %s vs %s. %g %d games\n", m.Role, CHAMPION_KEYS_BY_KEY[m.Champion], CHAMPION_KEYS_BY_KEY[m.Enemy], m.WinRate, m.Games)
}

func (c *ChampionMatchup) UpdatePersonalData(p PersonalMatchup) {
	c.PersonalWinRate = strings.TrimSuffix(fmt.Sprintf("%0.2f", p.WinRate), ".00")
	c.PersonalGames = p.Games
}

type PersonalMatch struct {
	Lane             string
	LoserChampionID  string
	LoserSummonerID  uint64 `sql:",pk"`
	MatchId          uint64 `sql:",pk"`
	PlatformId       string
	QueueType        string
	Region           string `sql:",pk"`
	Season           string
	Timestamp        uint64
	WinnerChampionID string
	WinnerSummonerID uint64 `sql:",pk"`
}

type PersonalMatchup struct {
	Champion string
	WinRate  float64
	Games    int
}

func setupRiotApi() error {
	key := lol.APIKey(viper.GetString("riot.key"))
	apiEndpointMap = make(map[string]*lol.APIEndpoint)
	for r, _ := range RIOT_REGIONS {
		region, _ := lol.NewRegionByCode(r)
		api, err := lol.NewAPIEndpoint(region, key)
		if err != nil {
			return err
		}
		apiEndpointMap[r] = api
	}
	return nil
}

func createSchema(db *pg.DB) error {
	queries := []string{

		`CREATE TABLE IF NOT EXISTS champion_matchups (
      champion text,
      enemy text,
      games int,
      role text,
      stat_score decimal,
      win_rate decimal,
      PRIMARY KEY(champion, enemy, role))`,

		`CREATE TABLE IF NOT EXISTS masteries (
      champion_id bigint,
      champion_level int,
      champion_points int,
      champion_points_since_last_level bigint,
      champion_points_until_next_level bigint,
      last_play_time bigint,
			region text,
      summoner_id bigint,
      PRIMARY KEY (champion_id, region, summoner_id))`,

		`CREATE TABLE IF NOT EXISTS personal_matches (
      lane         text,
      loser_champion_id  text,
      loser_summoner_id  bigint,
      match_id           bigint,
      platform_id        text,
      queue_type         text,
      region             text,
      season             text,
      timestamp          bigint,
      winner_champion_id text,
      winner_summoner_id bigint,
      PRIMARY KEY (match_id, loser_summoner_id, winner_summoner_id, region))`,

		`CREATE TABLE IF NOT EXISTS my_summoners (
      name                 text,
      profile_icon_id      int,
      masteries_updated_at timestamp with time zone DEFAULT (now() at time zone 'utc'),
      matchlist_timestamp  bigint DEFAULT 0,
      matchlist_updated_at timestamp with time zone DEFAULT (now() at time zone 'utc'),
			region               text,
      revision_date        bigint,
			summoner_id          bigint,
      summoner_level       int,
			PRIMARY KEY (region, summoner_id))`,
	}
	queries = append(queries, createTableSql...)

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			fmt.Println("Failed to exec postgres table schemas", err)
			return err
		}
	}
	return nil
}

func getSummonerMasteriesAndSave(region string, summonerName string, db *pg.DB) (err error) {
	names := NormalizeSummonerName(summonerName)
	summoners, err := getSummonerIDByNameAndSave(region, names, db)
	if err != nil {
		return
	}
	_, err = getChampionMasteriesBySummonerIDAndSave(region, summoners[0], db)
	return
}

func getSummonerIDByNameAndSave(region string, names []string, db *pg.DB) (players []MySummoner, err error) {
	summoners, err := apiEndpointMap[region].GetSummonerByName(names)
	if err != nil {
		fmt.Println("Failed to get summoners:", names, "region:", region, err)
		err = errors.New("Summoner does not exist")
		return
	}
	for _, summoner := range summoners {
		if summoner.ID == 0 {
			fmt.Println("Summoner not found in returning json from riot, summoner name:", names, "region:", region, err)
			err = errors.New("Failed to get summoner")
			return
		} else if summoner.Level < 30 {
			fmt.Println("Summoner is not level 30 yet, summoner name:", names, "region:", region)
			err = errors.New("Summoner is not level 30")
			return
		}
		s := MySummoner{
			NormalizeSummonerName(summoner.Name)[0],
			summoner.ProfileIconID,
			time.Now().UTC(),
			0,
			time.Now().UTC(),
			region,
			uint64(summoner.RevisionDate),
			uint64(summoner.ID),
			summoner.Level}
		err = db.Create(&s)
		if err != nil {
			fmt.Println("Failed to save summoner:", summoner.Name, "region:", region, "in db.", err)
			err = errors.New("Failed to save summoner")
			return
		}
		players = append(players, s)
	}
	return
}

func getChampionMasteriesBySummonerIDAndSave(region string, summoner MySummoner, db *pg.DB) (masteries []lol.ChampionMastery, err error) {
	masteries, err = apiEndpointMap[region].GetChampionMasteries(lol.SummonerID(summoner.SummonerID))
	if err != nil {
		fmt.Println("Failed to get champion masteries for summoner:", summoner.SummonerID, "region:", region, err)
		err = errors.New("Failed to get champion masteries")
		return
	}

	for _, mastery := range masteries {
		_, err = db.Model(&Mastery{
			int(mastery.Champion),
			mastery.Level,
			mastery.Points,
			mastery.PointsSinceLastLevel,
			mastery.PointsUntilNextLevel,
			uint64(mastery.LastPlayTime),
			region,
			uint64(mastery.Player),
		}).OnConflict("(champion_id, region, summoner_id) DO UPDATE").Set(`
            champion_level = ?champion_level,
            champion_points = ?champion_points,
            champion_points_since_last_level = ?champion_points_since_last_level,
            champion_points_until_next_level = ?champion_points_until_next_level,
            last_play_time = ?last_play_time`).Create()
		if err != nil {
			fmt.Println("Failed to update masteries for summoner:", summoner.SummonerID, "region:", region, err)
			// err = errors.New("Failed to save champion masteries")
			err = nil
			return
		}
	}
	var s MySummoner
	_, err = db.Model(&s).Set("masteries_updated_at = ?", time.Now().UTC()).Where("summoner_id = ?", summoner.SummonerID).Update()
	if err != nil {
		fmt.Println("Failed to update summoner:", summoner.SummonerID, "region:", region, "masteries updated at time.", err)
		// err = errors.New("Failed to update summoner masteries updated time")
		// return
	}
	return
}

/* ======================== db methods ======================= */

func getOrCreateSummoner(region string, summonerName string, db *pg.DB) (summoner MySummoner, err error) {
	name := NormalizeSummonerName(summonerName)[0]
	err = db.Model(&summoner).Where("name = ?", name).Select()
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			fmt.Println("Summoner not found, lets add one with name:", summonerName, "region:", region)
			err = getSummonerMasteriesAndSave(region, name, db)
			if err != nil {
				return
			}
			err = db.Model(&summoner).Where("name = ?", name).Select()
			if err != nil {
				fmt.Println("Failed to find summoner with name:", summonerName, "region:", region, "in db", err)
				err = errors.New("Failed to find summoner in db")
				return
			}
			getMatchlistBySummonerIDAndSave(region, summoner, db)
		} else {
			fmt.Println("Failed to get summoner:", summonerName, "region:", region, err)
			err = errors.New("Failed to find summoner in db")
			return
		}
	} else {
		// fmt.Println("checking is masteries or personal matches need updating...")
		if summoner.MasteriesUpdatedAt.UTC().Add(time.Duration(60*60*12) * time.Second).Before(time.Now().UTC()) {
			fmt.Printf("Updating masteries for Summoner: %v, region: %s!\n", summoner.SummonerID, region)
			getChampionMasteriesBySummonerIDAndSave(region, summoner, db)
		}
		if summoner.MatchlistUpdatedAt.UTC().Add(time.Duration(60*60*3) * time.Second).Before(time.Now().UTC()) {
			fmt.Printf("Updating matchlist for Summoner: %v: region: %s!\n", summoner.SummonerID, region)
			getMatchlistBySummonerIDAndSave(region, summoner, db)
		}
	}

	log.Printf("getOrCreateSummoner: %+v\n", summoner)
	return
}

func getMatchups(region string, summoner MySummoner, enemyChampionID string, role string, db *pg.DB) (matchups []ChampionMatchup, pms []PersonalMatchup, err error) {
	//select * from champion_matchups where enemy = '57' and champion IN (select cast(champion_id as text)  from masteries where summoner_id = 26691960) order by win_rate desc;
	// defer un(trace("get matchup"))
	lane := strings.ToLower(role)
	err = db.Model(&matchups).Where("role = ? AND enemy = ? AND champion IN (SELECT CAST(champion_id AS text) FROM masteries WHERE summoner_id = ? and region = ?)", role, enemyChampionID, summoner.SummonerID, region).Order("win_rate desc").Select()
	if err != nil {
		fmt.Println("Failed to get matchups from database for summoner:", summoner.SummonerID, "region:", region, err)
		err = errors.New("Failed to get championgg matchups")
		return
	}

	sql := fmt.Sprintf(
		`
      SELECT champion, (COALESCE(0.0+loss_count, 0) / games)*100 as win_rate, games
      FROM
      (
          SELECT loss_count, COALESCE(wins.loser_champion_id, losses.winner_champion_id) as champion,
                (COALESCE(win_count, 0) + COALESCE(0.0+loss_count, 0))::bigint as games
          FROM (
              SELECT winner_champion_id, loser_champion_id, count(*) as win_count
              FROM personal_matches
              WHERE winner_champion_id = '%s' and lane = '%s' and loser_summoner_id = %v
              GROUP BY winner_champion_id, loser_champion_id
          )
          as wins

          FULL JOIN
          (
              SELECT winner_champion_id, loser_champion_id, count(*) as loss_count
              FROM personal_matches
              WHERE loser_champion_id = '%s' and lane = '%s' and winner_summoner_id = %v
              GROUP BY winner_champion_id, loser_champion_id
          ) as losses ON row(wins.winner_champion_id, wins.loser_champion_id) = row(losses.loser_champion_id, losses.winner_champion_id)
      ) as g`,
		enemyChampionID,
		lane,
		summoner.SummonerID,
		enemyChampionID,
		lane,
		summoner.SummonerID)
	log.Println(sql)
	_, err = db.Query(&pms, sql)
	if err != nil {
		fmt.Println("Failed to get personal winrates for summoner:", summoner.SummonerID, "region:", region)
		// err = errors.New("Failed to get personal winrates")
		// return
	}
	log.Println("personal matchups:", pms)
	return
}

func getMatchlistBySummonerIDAndSave(region string, summoner MySummoner, db *pg.DB) (err error) {
	matchList, err := apiEndpointMap[region].GetMatchlist(summoner.SummonerID, summoner.MatchlistTimestamp)
	if err != nil {
		fmt.Println("Failed to get matchlist for summoner:", summoner.SummonerID, "region:", region, err)
		return
	}

	if len(matchList.Matches) == 0 {
		summoner.UpdateMatchlistUpdatedAt()
		err = updateMatchlistUpdatedAt(summoner)
		if err != nil {
			fmt.Println("Failed to update summoner:", summoner.SummonerID, "region:", region, "matchlist updated at time.", err)
			err = errors.New("Failed to update matchlist updated at time")
		}
		return
	}

	fmt.Printf("Going to determine the winners and losers for 50/%d matches for summoner: %v, region: %s.\n", len(matchList.Matches), summoner.SummonerID, region)
	ids := ""
	for i, m := range matchList.Matches {
		ids += fmt.Sprintf("(%v)", m.MatchId)
		if i < len(matchList.Matches)-1 {
			ids += ", "
		}
	}
	log.Println(ids)
	//FROM unnest('{%s}'::bigint[]) ids(v)
	sql := fmt.Sprintf(
		`SELECT ids.v
     FROM (VALUES %s) ids(v)
     LEFT JOIN personal_matches ON ids.v = personal_matches.match_id
     WHERE personal_matches.match_id IS NULL`,
		ids)
	log.Println("sql:", sql)
	matchIDs := make(sortutil.Uint64Slice, 0)
	_, err = db.Query(&matchIDs, sql)
	if err != nil {
		fmt.Println("Failed to get list of un-generated personal matches from db for summoner:", summoner.SummonerID, "region:", region, err)
		err = errors.New("Failed to get un-generated personal matches")
		return
	}
	log.Println("matchIDs:", matchIDs)
	matchIDs.Sort()
	matchIDsCount := 0
	var matchlistUpdatedTimestamp uint64 = math.MaxUint64
	for _, m := range matchList.Matches {
		index := sortutil.SearchUint64s(matchIDs, m.MatchId)
		if index >= len(matchIDs) {
			continue
		} else if uint64(matchIDs[index]) == m.MatchId {
			log.Println("analyzing game:", m.MatchId)
			if matchIDsCount >= 50 {
				break
			}

			game, err := apiEndpointMap[region].GetMatch(m.MatchId, false)
			if err != nil {
				fmt.Println("Failed to get ranked game:", m.MatchId, "for summoner:", summoner.SummonerID, "region:", region, err)
				if err.Error() == "Too Many request to server" {
					if m.Timestamp < matchlistUpdatedTimestamp {
						matchlistUpdatedTimestamp = m.Timestamp
					}
					matchIDsCount++
				}
				continue
			}
			matchIDsCount++
			var pmToInsert []PersonalMatch
			winnerLosers := generateMatchups(game)
			for _, w := range winnerLosers {
				if w.LoserChampionID == "" || w.LoserSummonerID == 0 || w.WinnerChampionID == "" || w.WinnerSummonerID == 0 {
					continue
				}
				pmToInsert = append(pmToInsert, PersonalMatch{
					w.Lane,
					w.LoserChampionID,
					w.LoserSummonerID,
					m.MatchId,
					m.PlatformId,
					m.Queue,
					m.Region,
					m.Season,
					m.Timestamp,
					w.WinnerChampionID,
					w.WinnerSummonerID})
			}
			_, err = db.Model(&pmToInsert).OnConflict("DO NOTHING").Create()
			if err != nil {
				fmt.Println("failed to bulk insert personal matches for match:", m.MatchId, "summoner:", summoner.SummonerID, "region:", region, err)
				continue
			}
		} else {
			continue
		}
	}

	if matchlistUpdatedTimestamp == math.MaxUint64 {
		matchlistUpdatedTimestamp = matchList.Matches[0].Timestamp
	}

	//NOTE so if it failed, we'll still wont try again for a day
	summoner.UpdateMatchlistTimestamp(matchlistUpdatedTimestamp)
	err = updateSummonerMatchlistData(summoner)
	if err != nil {
		fmt.Println("Failed to update summoner:", summoner.SummonerID, "region:", region, "matchlist updated at and matchlist timestamp.", err)
		err = errors.New("Failed to update matchlist updated at time and matchlist timestamp")
		return
	}
	return nil
}

func updateMatchlistUpdatedAt(summoner MySummoner) (err error) {
	var s MySummoner
	_, err = db.Model(&s).Set("matchlist_updated_at = ?", s.MatchlistUpdatedAt).Where("summoner_id = ?", summoner.SummonerID).Update()
	return
}

func updateSummonerMatchlistData(summoner MySummoner) (err error) {
	var s MySummoner
	_, err = db.Model(&s).Set("matchlist_updated_at = ?, matchlist_timestamp = ?", summoner.MatchlistUpdatedAt, summoner.MatchlistTimestamp).Where("summoner_id = ?", summoner.SummonerID).Update()
	return
}

/* ======================== helper ========================== */

func createSummonerIDString(summonerID []int64) (summonerIDstr string, err error) {
	if len(summonerID) > 40 {
		return summonerIDstr, errors.New("A Maximum of 40 SummonerIDs are allowed")
	}
	for k, v := range summonerID {
		summonerIDstr += strconv.FormatInt(v, 10)
		if k != len(summonerID)-1 {
			summonerIDstr += ","
		}
	}
	return
}

func NormalizeSummonerName(summonerNames ...string) []string {
	for i, v := range summonerNames {
		summonerName := strings.ToLower(v)
		summonerName = strings.Replace(summonerName, " ", "", -1)
		summonerNames[i] = summonerName
	}
	return summonerNames
}

func NormalizeChampion(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "", -1)
	name = strings.Replace(name, "'", "", -1)
	name = strings.Replace(name, ".", "", -1)
	return name
}

func generateMatchups(game *lol.Match) (winnerLosers []WinnerLoser) {
	// defer un(trace("generate matchups"))
	team1Strings := []string{}
	team2Strings := []string{}
	for _, p := range game.Participants {
		t := fmt.Sprintf("%s:%s:%d",
			p.Timeline.Role,
			p.Timeline.Lane,
			p.ChampionId)
		if p.TeamId == 100 {
			team1Strings = append(team1Strings, t)
		} else {
			team2Strings = append(team2Strings, t)
		}
	}
	sort.Strings(team1Strings)
	sort.Strings(team2Strings)
	log.Println("Team strings:", team1Strings, team2Strings)
	lanes1 := checkTemplate(team1Strings, game.MatchId)
	lanes2 := checkTemplate(team2Strings, game.MatchId)
	sort.Strings(lanes1)
	sort.Strings(lanes2)
	log.Println("Lanes:", lanes1, lanes2)
	winnerLosers = make([]WinnerLoser, 5)
	for i, l := range append(lanes1, lanes2...) {
		colonIndex := strings.LastIndex(l, ":")
		champID, _ := strconv.Atoi(l[colonIndex+1:])
		for _, p := range game.Participants {
			if champID == p.ChampionId {
				for _, pl := range game.ParticipantIdentities {
					if pl.ParticipantId == p.ParticipantId {
						i = int(math.Mod(float64(i), 5))
						winnerLosers[i].Lane = l[:colonIndex]
						if p.Stat.Winner {
							winnerLosers[i].WinnerSummonerID = pl.Player.SummonerId
							winnerLosers[i].WinnerChampionID = l[colonIndex+1:]
						} else {
							winnerLosers[i].LoserSummonerID = pl.Player.SummonerId
							winnerLosers[i].LoserChampionID = l[colonIndex+1:]
						}
					}
				}
			}
		}
	}
	// fmt.Printf("WinnerLoser generated: %+v\n", winnerLosers)
	return
}

func checkTemplate(teamStrings []string, matchId uint64) []string {
	lanes := []string{}
	bestNumMatches := -1
	for i := 0; i < len(GAME_TEMPLATES); i += 2 {
		log.Println("checking GAME_TEMPLATES ", (i/2)+1)
		numMatches := 0
		tmpLanes := []string{}
		for j, comparer := range GAME_TEMPLATES[i] {
			currRoleLaneChampID := teamStrings[j]
			// fmt.Printf("currRoleLaneChampID: %s\n", currRoleLaneChampID)
			// fmt.Printf("comparer at %d: %s\n", j, comparer)
			champIDColonIndex := strings.LastIndex(currRoleLaneChampID, ":")
			conditionColonIndex := strings.LastIndex(comparer, ":")
			champID := currRoleLaneChampID[champIDColonIndex+1:]
			currRoleLane := currRoleLaneChampID[:champIDColonIndex+1]
			templateRoleLane := comparer[:conditionColonIndex+1]
			condition := strings.ToLower(comparer[conditionColonIndex+1:])

			log.Printf("cond: %s, champID: %s, lanes: %v\n", condition, champID, championToRoleMapping[champID])
			if condition == "*" {
				// eg. *
				if templateRoleLane == currRoleLane {
					tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":"+champID)
					log.Println("passed condition ", condition)
					numMatches++
				} else {
					tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":")
				}
			} else {
				if contains(championToRoleMapping[champID], condition) {
					// eg. MIDDLE
					tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":"+champID)
					log.Println("passed condition ", condition)
					numMatches++
					continue
				}
				if strings.HasPrefix(condition, "!") {
					// eg. !ADC
					if !contains(championToRoleMapping[champID], condition[1:]) {
						tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":"+champID)
						log.Println("passed condition ", condition)
						numMatches++
						continue
					}
				}

				log.Println("j+1:", j+1, GAME_TEMPLATES[i])

				if j+1 < len(GAME_TEMPLATES[i]) {
					log.Println(GAME_TEMPLATES[i][j+1], currRoleLane, conditionColonIndex, strings.LastIndex(GAME_TEMPLATES[i][j+1], ":"))
					if strings.Contains(GAME_TEMPLATES[i][j+1], currRoleLane) {
						// does the next condition have the same ROLE:LANE prefix?
						tmpConditionColonIndex := strings.LastIndex(GAME_TEMPLATES[i][j+1], ":")
						tmpCondition := GAME_TEMPLATES[i][j+1][tmpConditionColonIndex+1:]
						log.Println("tmp cond:", tmpCondition)
						if strings.HasPrefix(tmpCondition, "!") {
							// eg. !ADC
							if !contains(championToRoleMapping[champID], tmpCondition[1:]) {
								GAME_TEMPLATES[i][j], GAME_TEMPLATES[i][j+1] = GAME_TEMPLATES[i][j+1], GAME_TEMPLATES[i][j]
								tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":"+champID)
								log.Println("passed condition ", tmpCondition, "by swapping with", j+1)
								numMatches++
							} else {
								tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":")
							}
						} else if contains(championToRoleMapping[champID], strings.ToLower(tmpCondition)) {
							// check if the next template string can work and swap next with current and do the same check
							GAME_TEMPLATES[i][j], GAME_TEMPLATES[i][j+1] = GAME_TEMPLATES[i][j+1], GAME_TEMPLATES[i][j]
							tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":"+champID)
							log.Println("passed condition ", tmpCondition, "by swapping with", j+1)
							numMatches++
						} else {
							tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":")
						}
					} else {
						tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":")
					}
				} else {
					tmpLanes = append(tmpLanes, GAME_TEMPLATES[i+1][j]+":")
				}
			}
		}
		if numMatches == 5 {
			log.Println("matched GAME_TEMPLATES ", i/2)
			return tmpLanes
		}
		if numMatches > bestNumMatches {
			bestNumMatches = numMatches
			lanes = tmpLanes
		}
	}
	fmt.Println("ERROR, match_id:", matchId)
	for _, p := range teamStrings {
		lastColon := strings.LastIndex(p, ":")
		fmt.Println(p[:lastColon+1]+CHAMPION_KEYS_BY_KEY[p[lastColon+1:]], p[lastColon+1:], championToRoleMapping[p[lastColon+1:]])
	}
	return lanes
}

func trace(s string) (string, time.Time) {
	// fmt.Println("START:", s)
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	fmt.Println(s, endTime.Sub(startTime))
}
