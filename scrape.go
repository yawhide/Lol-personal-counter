package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	// "runtime"
	"strings"
	"sync"
	"time"

	"github.com/juju/ratelimit"
	"github.com/spf13/viper"
	"gopkg.in/pg.v4"
)

var championToRoleMapping = map[string][]string{}
var matchIDs map[string]uint64
var db *pg.DB

var brLock = &sync.Mutex{}
var euneLock = &sync.Mutex{}
var euwLock = &sync.Mutex{}
var jpLock = &sync.Mutex{}
var krLock = &sync.Mutex{}
var lanLock = &sync.Mutex{}
var lasLock = &sync.Mutex{}
var naLock = &sync.Mutex{}
var oceLock = &sync.Mutex{}
var trLock = &sync.Mutex{}
var ruLock = &sync.Mutex{}

var brMatchID uint64
var euneMatchID uint64
var euwMatchID uint64
var jpMatchID uint64
var krMatchID uint64
var lanMatchID uint64
var lasMatchID uint64
var naMatchID uint64
var oceMatchID uint64
var trMatchID uint64
var ruMatchID uint64

type ChampionRoleMapping struct {
	Champion string
	Role     string
}

type MatchStore struct {
	MatchID uint64
	Region  string
	Game    string
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	username := viper.GetString("postgres.username")
	password := viper.GetString("postgres.password")

	db = pg.Connect(&pg.Options{
		User:     username,
		Password: password,
	})

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

	createTableSql = append(createTableSql, `
    CREATE TABLE IF NOT EXISTS match_stores (
    match_id BIGINT,
    region TEXT,
    game JSON,
    PRIMARY KEY (match_id, region))`)

	err = createSchema(db)
	if err != nil {
		panic(err)
	}

	err = setupRiotApi()
	if err != nil {
		panic(err)
	}

	// brMatchID = matchIDs["br"]
	// euneMatchID = matchIDs["eune"]
	// euwMatchID = matchIDs["euw"]
	// jpMatchID = matchIDs["jp"]
	// krMatchID = matchIDs["kr"]
	// lanMatchID = matchIDs["lan"]
	// lasMatchID = matchIDs["las"]
	naMatchID = readMatchID("na")
	log.Println("Read na match id:", naMatchID)
	// oceMatchID = matchIDs["oce"]
	// trMatchID = matchIDs["tr"]
	// ruMatchID = matchIDs["ru"]

	waitTime, _ := time.ParseDuration("1050ms")
	log.Println("wait time between api calls is:", waitTime)

	bucket := ratelimit.NewBucketWithRate(100, 100)
	counter := 0
	var failedAPICalls = make([]uint64, 0)
	for {
		if counter == 60*5 {
			saveMatchID(naMatchID, "na")
			counter = 0
			panic(errors.New("WTF"))
		}
		now := time.Now()
		naLock.Lock()
		var failedIDs []uint64
		if len(failedAPICalls) < 100 {
			failedIDs = failedAPICalls[:len(failedAPICalls)]
			failedAPICalls = nil
		} else {
			failedIDs = failedAPICalls[:100]
			failedAPICalls = append(failedAPICalls[100:], failedAPICalls[len(failedAPICalls):]...)
		}
		naLock.Unlock()
		log.Println("retrying", len(failedIDs), "match IDs")
		for i := bucket.Available() - 1; i >= 0; i-- {
			var lastMatchID uint64
			if i < int64(len(failedIDs)) {
				lastMatchID = failedIDs[i]
			} else {
				lastMatchID = naMatchID
				naMatchID++
			}
			bucket.Wait(1)
			go func(ID uint64) {
				err = scrape("na", ID)
				if err != nil {
					errStr := strings.TrimSpace(err.Error())
					if strings.HasSuffix(errStr, "Too Many request to server") {
						log.Println("hit a real 429, naOps:", i, "naTime:", time.Now().Sub(now))
						panic(err)
					}
					if !strings.HasSuffix(errStr, "404") {
						naLock.Lock()
						failedAPICalls = append(failedAPICalls, ID)
						naLock.Unlock()
					}
				}
			}(lastMatchID)
		}
		// fmt.Println("done invoking 140 api calls...current MatchID:", naMatchID, "sleeping for:", waitTime-time.Now().Sub(now))
		time.Sleep(waitTime - time.Now().Sub(now))
		counter++
	}
}

func scrape(region string, matchID uint64) error {
	game, err := apiEndpointMap[region].GetMatch(matchID, false)
	if err != nil {
		errStr := strings.TrimSpace(err.Error())
		if strings.HasSuffix(errStr, "404") {
			// game doesnt exist
		} else {
			log.Println("MatchID:", matchID, "region:", region, errStr)
		}
		return err
	}
	w, err := json.Marshal(game)
	if err != nil {
		log.Println("Failed to marshal the game... MatchID:", game.MatchID, "region:", region, err)
		return err
	}
	_, err = db.Model(&MatchStore{game.MatchID, game.Region, string(w)}).OnConflict("DO NOTHING").Create()

	if err != nil {
		log.Println("Failed to insert the game... MatchID:", game.MatchID, "region:", region, err)
		return err
	}
	return nil
}

func saveMatchID(ID uint64, region string) {
	log.Printf("Encoded: ")
	log.Println(ID)
	e := new(bytes.Buffer)
	encoder := gob.NewEncoder(e)
	encoder.Encode(ID)
	err := ioutil.WriteFile(fmt.Sprintf("data/%s_id", region), e.Bytes(), 0644)
	if err != nil {
		log.Println("Failed to save match id for region:", region, err)
	} else {
		log.Println("Finished updating the match id file with match id:", ID, "region:", region)
	}
}

func readMatchID(region string) (ID uint64) {
	data, err := ioutil.ReadFile(fmt.Sprintf("data/%s_id", region))
	if err != nil {
		panic(err)
	}
	// read in matchIDDump and load values in memory

	matchIDBuffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(matchIDBuffer)
	decoder.Decode(&ID)

	log.Printf("Decoded: ")
	log.Println(ID)
	if ID != 0 {
		return ID
	}
	log.Println("ID read for region:", region, "is 0.... we are starting from scratch")
	return 2006095010
}
