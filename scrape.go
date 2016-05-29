package main

import (
  "strings"
)

var lastMatchID uint64

func init() {
  //TODO get last matchid processed
  for ; ; lastMatchID++ {
    for _, r := range RIOT_REGIONS {

      game, err := apiEndpointMap[r].GetMatch(lastMatchID, false)
      if err != nil {
        errStr := strings.TrimSpace(err.Error())
        if strings.HasSuffix(errStr, "404") {
          // game doesnt exist
        } else if strings.HasSuffix(errStr, "500") {
          // try again
        } else {

        }
        continue
      }
      db.Model
    }
  }
}
fasdf
