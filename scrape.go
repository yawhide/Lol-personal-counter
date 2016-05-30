package main

import (
  "bytes"
  "encoding/gob"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  // "runtime"
  "strings"
  "sync"
  "sync/atomic"
  "time"

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

/*  for region, _ := range RIOT_REGIONS {
    regionLockMaps[region] = &sync.Mutex{}
  }*/

  // read in matchIDDump and load values in memory

  matchIDDump, err := ioutil.ReadFile("matchIDDump")
  if err != nil {
    panic(err)
  }
  matchIDBuffer := bytes.NewBuffer(matchIDDump)
  decoder := gob.NewDecoder(matchIDBuffer)
  decoder.Decode(&matchIDs)

  log.Printf("Decoded: ")
  log.Println(matchIDs)

  brMatchID = matchIDs["br"]
  euneMatchID = matchIDs["eune"]
  euwMatchID = matchIDs["euw"]
  jpMatchID = matchIDs["jp"]
  krMatchID = matchIDs["kr"]
  lanMatchID = matchIDs["lan"]
  lasMatchID = matchIDs["las"]
  naMatchID = matchIDs["na"]
  oceMatchID = matchIDs["oce"]
  trMatchID = matchIDs["tr"]
  ruMatchID = matchIDs["ru"]

  waitTime, _ := time.ParseDuration("1500ms")
  log.Println("wait time between api calls is:", waitTime)
  // waitTime := time.Second * 2

  naTime := time.Now()
  var naOps int64 = 0
  for i := 0; i < 140; i++ {
    go func() {
      for {
        // log.Println(region)
        now := time.Now()
        naLock.Lock()
        lastMatchID := naMatchID
        naMatchID += 1
        naLock.Unlock()

        err = scrape("na", lastMatchID)
        atomic.AddInt64(&naOps, 1)
        if err != nil {
          errStr := strings.TrimSpace(err.Error())
          if strings.HasSuffix(errStr, "Too Many request to server") {
            log.Println("hit a 429, naOps:", naOps, "naTime:", time.Now().Sub(naTime))
            panic(err)
          }
        }
        // log.Println(region, err)
        timeElasped := time.Now().Sub(now)
        if timeElasped < waitTime {
          // log.Println("sleep for:", time.Second - timeElasped, waitTime - timeElasped)
          time.Sleep(waitTime - timeElasped)
        }
      }
    }()
  }

  for {
  // update the regionCurrMatchIDs every
  // time.Duration(60*5) * time.Second
  // make sure to acquire all locks
    time.Sleep(60*1 * time.Second)
    brLock.Lock()
    euneLock.Lock()
    euwLock.Lock()
    jpLock.Lock()
    krLock.Lock()
    lanLock.Lock()
    lasLock.Lock()
    naLock.Lock()
    oceLock.Lock()
    trLock.Lock()
    ruLock.Lock()
    matchIDs["br"] = brMatchID
    matchIDs["eune"] = euneMatchID
    matchIDs["euw"] = euwMatchID
    matchIDs["jp"] = jpMatchID
    matchIDs["kr"] = krMatchID
    matchIDs["lan"] = lanMatchID
    matchIDs["las"] = lasMatchID
    matchIDs["na"] = naMatchID
    matchIDs["oce"] = oceMatchID
    matchIDs["tr"] = trMatchID
    matchIDs["ru"] = ruMatchID
    log.Printf("Encoded: ")
    log.Println(matchIDs)
    e := new(bytes.Buffer)
    encoder := gob.NewEncoder(e)
    encoder.Encode(matchIDs)
    err = ioutil.WriteFile("matchIDDump", e.Bytes(), 0644)
    if err != nil {
      panic(err)
    }
    brLock.Unlock()
    euneLock.Unlock()
    euwLock.Unlock()
    jpLock.Unlock()
    krLock.Unlock()
    lanLock.Unlock()
    lasLock.Unlock()
    naLock.Unlock()
    oceLock.Unlock()
    trLock.Unlock()
    ruLock.Unlock()
    log.Println("Finished updating the match id file")
  }
  select {}

}


func scrape(region string, matchID uint64) error {
  game, err := apiEndpointMap[region].GetMatch(matchID, false)
  if err != nil {
    errStr := strings.TrimSpace(err.Error())
    if strings.HasSuffix(errStr, "404") {
      // game doesnt exist
    } else if strings.HasSuffix(errStr, "500") {
      // try again
      log.Println("MatchID:", matchID, "region:", region, errStr)
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
  _, err = db.Model(&MatchStore{ game.MatchID, game.Region, string(w) }).OnConflict("DO NOTHING").Create()

  if err != nil {
    log.Println("Failed to insert the game... MatchID:", game.MatchID, "region:", region, err)
    return err
  }
  return nil
}
