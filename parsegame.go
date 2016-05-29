package main

import (
	"errors"
	"fmt"

	"github.com/yawhide/go-lol"
)

type WinnerLoser struct {
	Lane             string
	LoserChampionID  string
	LoserSummonerID  uint64
	WinnerChampionID string
	WinnerSummonerID uint64
}

var GAME_TEMPLATES = [][]string{}

func init() {
	// ROLE:LANE:GUESS
	t1 := []string{
		`DUO_CARRY:BOTTOM:*`,
		`DUO_SUPPORT:BOTTOM:*`,
		`NONE:JUNGLE:*`,
		`SOLO:MIDDLE:*`,
		`SOLO:TOP:*`}
	a1 := []string{
		"adc",
		"support",
		"jungle",
		"middle",
		"top"}
	t2 := []string{
		`DUO_CARRY:BOTTOM:*`,
		`DUO_CARRY:MIDDLE:MIDDLE`,
		`DUO_SUPPORT:BOTTOM:*`,
		`DUO_SUPPORT:MIDDLE:JUNGLE`,
		`SOLO:TOP:*`}
	a2 := []string{
		"adc",
		"middle",
		"support",
		"jungle",
		"top"}
	t3 := []string{
		`NONE:JUNGLE:JUNGLE`,
		`NONE:JUNGLE:SUPPORT`,
		`SOLO:BOTTOM:*`,
		`SOLO:MIDDLE:*`,
		`SOLO:TOP:*`}
	a3 := []string{
		"jungle",
		"support",
		"adc",
		"middle",
		"top"}
	t4 := []string{
		`DUO_CARRY:MIDDLE:MIDDLE`,
		`DUO_SUPPORT:MIDDLE:*`,
		`NONE:JUNGLE:*`,
		`SOLO:BOTTOM:*`,
		`SOLO:TOP:*`}
	a4 := []string{
		"middle",
		"support",
		"jungle",
		"adc",
		"top"}
	t5 := []string{
		"DUO:BOTTOM:!ADC",
		"DUO:BOTTOM:ADC",
		"NONE:JUNGLE:*",
		"SOLO:MIDDLE:*",
		"SOLO:TOP:*"}
	a5 := []string{
		"support",
		"adc",
		"jungle",
		"middle",
		"top"}
	t6 := []string{
		"DUO_CARRY:BOTTOM:*",
		"DUO_SUPPORT:BOTTOM:*",
		"NONE:JUNGLE:JUNGLE",
		"NONE:JUNGLE:MIDDLE",
		"SOLO:TOP:*"}
	a6 := []string{
		"adc",
		"support",
		"jungle",
		"middle",
		"top"}
	t7 := []string{ // tough one...
		"DUO_CARRY:MIDDLE:TOP",
		"DUO_SUPPORT:MIDDLE:MIDDLE",
		"NONE:JUNGLE:!ADC",
		"NONE:JUNGLE:ADC",
		"SOLO:BOTTOM:SUPPORT"}
	a7 := []string{
		"top",
		"middle",
		"jungle",
		"adc",
		"support"}
	/* ERROR.
	   DUO_CARRY:MIDDLE:riven 92 [top]
	   DUO_SUPPORT:MIDDLE:karma 43 [middle support]
	   NONE:JUNGLE:nautilus 111 [support top]
	   NONE:JUNGLE:caitlyn 51 [adc]
	   SOLO:BOTTOM:sona 37 [support]
	*/
	t8 := []string{
		"DUO:BOTTOM:!ADC",
		"DUO:BOTTOM:ADC",
		"DUO_CARRY:MIDDLE:JUNGLE",
		"DUO_SUPPORT:MIDDLE:MIDDLE",
		"SOLO:TOP:*"}
	a8 := []string{
		"support",
		"adc",
		"jungle",
		"middle",
		"top"}
	t9 := []string{
		"DUO:BOTTOM:!ADC",
		"DUO:BOTTOM:ADC",
		"DUO_CARRY:MIDDLE:MIDDLE",
		"DUO_SUPPORT:MIDDLE:JUNGLE",
		"SOLO:TOP:*"}
	a9 := []string{
		"support",
		"adc",
		"middle",
		"jungle",
		"top"}
	t10 := []string{
		"DUO:BOTTOM:!ADC",
		"DUO:BOTTOM:ADC",
		"NONE:JUNGLE:JUNGLE",
		"NONE:JUNGLE:MIDDLE",
		"SOLO:TOP:*"}
	a10 := []string{
		"support",
		"adc",
		"jungle",
		"middle",
		"top"}
	t11 := []string{
		`NONE:JUNGLE:JUNGLE`,
		`NONE:JUNGLE:ADC`,
		`SOLO:BOTTOM:*`,
		`SOLO:MIDDLE:*`,
		`SOLO:TOP:*`}
	a11 := []string{
		"jungle",
		"adc",
		"support",
		"middle",
		"top"}
	t12 := []string{ // tough one...
		"DUO_CARRY:BOTTOM:ADC",
		"DUO_SUPPORT:BOTTOM:SUPPORT",
		"NONE:JUNGLE:TOP",
		"NONE:JUNGLE:JUNGLE",
		"SOLO:TOP:MIDDLE"}
	a12 := []string{
		"adc",
		"support",
		"top",
		"jungle",
		"middle"}
	/*
	   DUO_CARRY:BOTTOM:lucian 236 [adc]
	   DUO_SUPPORT:BOTTOM:lulu 117 [top middle support]
	   NONE:JUNGLE:nasus 75 [top]
	   NONE:JUNGLE:nidalee 76 [jungle]
	   SOLO:TOP:brand 63 [support middle]
	*/

	GAME_TEMPLATES = append(GAME_TEMPLATES, t1)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a1)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t2)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a2)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t3)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a3)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t4)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a4)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t5)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a5)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t6)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a6)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t7)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a7)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t8)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a8)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t9)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a9)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t10)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a10)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t11)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a11)
	GAME_TEMPLATES = append(GAME_TEMPLATES, t12)
	GAME_TEMPLATES = append(GAME_TEMPLATES, a12)

}

func contains(arr []string, comparer string) bool {
	for _, s := range arr {
		if s == comparer {
			return true
		}
	}
	return false
}

func chooseMostLikelyForLane(game *lol.Match, lane string, formattedLane string) (err error, champID string) {
	for _, p := range game.Participants {
		if p.Timeline.Role == lane {
			formattedChampID := fmt.Sprintf("%d", p.ChampionID)
			if champID == "" {
				// nothing found yet, set it (this supports unconventional picks)
				champID = formattedChampID
			} else {
				// there are multiple in this lane... lets guess!
				savedChampIDFits := contains(championToRoleMapping[champID], formattedLane)
				currChampIDFits := contains(championToRoleMapping[formattedChampID], formattedLane)
				if savedChampIDFits && currChampIDFits {
					return errors.New("Multiple champs can be " + formattedLane), ""
				} else if !currChampIDFits {
					// keep the saved champ id
					continue
				} else if !savedChampIDFits {
					// make curr champ the saved champ
					champID = formattedChampID
				}
			}
		}
	}
	return nil, champID
}
