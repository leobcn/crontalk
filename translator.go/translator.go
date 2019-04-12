package translator

import (
	"crontalk/helper"

	"github.com/spf13/viper"
)

// TODO: Add validation for ranged values like the from range cannot be greater than the to range
// TODO: refacrot the huge validation code
// TODO: refactor timeoccurence more if possible
// TODO: reduce complexity
// TODO: add step values

const (
	english = "english"
	bangla  = "bangla"
)

var (
	weeks     = map[int]string{}
	months    = map[int]string{}
	baseIndex int
	configStr = ""
	language  = english
)

// Init initializes the translator
func Init() {
	if viper.GetBool(bangla) {
		language = bangla
		moments[dayIndex] = "দিন:" //not taking from config for the sake of simplicity
		moments[minuteIndex] = "মিনিট:"
		moments[hourIndex] = "ঘন্টা:"
	}
	configStr = "language." + language + "." //the config index to parse the config yaml file from viper

	weeks = map[int]string{
		0: viper.GetString(configStr + "sunday"),
		1: viper.GetString(configStr + "monday"),
		2: viper.GetString(configStr + "tuesday"),
		3: viper.GetString(configStr + "wednesday"),
		4: viper.GetString(configStr + "thursday"),
		5: viper.GetString(configStr + "friday"),
		6: viper.GetString(configStr + "saturday"),
	}
	months = map[int]string{
		1:  viper.GetString(configStr + "january"),
		2:  viper.GetString(configStr + "february"),
		3:  viper.GetString(configStr + "march"),
		4:  viper.GetString(configStr + "april"),
		5:  viper.GetString(configStr + "may"),
		6:  viper.GetString(configStr + "june"),
		7:  viper.GetString(configStr + "july"),
		8:  viper.GetString(configStr + "august"),
		9:  viper.GetString(configStr + "september"),
		10: viper.GetString(configStr + "october"),
		11: viper.GetString(configStr + "november"),
		12: viper.GetString(configStr + "december"),
	}
}

func translateBaseOccurence() error {
	var i int
	for i = weekIndex; i > hourIndex; i-- { //start iterating from the last sub-expressions to determine the starting string
		if cronSlice[i] != anyValue {
			cc, listed := helper.GetList(cronSlice[i], ",")
			for j, c := range cc { //iterating because values can be listed
				rr, ranged := helper.GetList(c, "-")
				t := translator{
					cron:          c,
					moment:        moments[i],
					cronRange:     rr,
					ranged:        ranged,
					listed:        listed,
					base:          true,
					cronListedLen: len(cc),
					index:         j,
				}
				if found, err := t.translateWeekMonth(); err != nil {
					return err
				} else if found {
					continue
				}
				t.translateDay()
			}
			break //once the base value is found no need for further iterations
		}
	}
	if i == hourIndex { // checking if every sub-expression contains asteriks apart from the time part
		translatedString += viper.GetString(configStr + "every_day")
	}
	baseIndex = i //storing the base index so that when checking every other than time , the base is also omitted because its
	//already checked
	return nil

}

func translateAllButBaseTimeOccurence() error {

	for i := dayIndex; i <= weekIndex; i++ { //checking every other sub-expressions apart from the base and time, no need for reverse travel
		if cronSlice[i] != anyValue && i != baseIndex { //not gonna check the base
			cc, listed := helper.GetList(cronSlice[i], ",")
			for j, c := range cc { //iterating the single sub-expressions
				rr, ranged := helper.GetList(c, "-")
				t := translator{
					cron:          c,
					moment:        moments[i],
					cronRange:     rr,
					ranged:        ranged,
					listed:        listed,
					base:          false,
					cronListedLen: len(cc),
					index:         j,
				}
				if found, err := t.translateWeekMonth(); err != nil {
					return err
				} else if !found {
					t.translateDay()
				}
			}
		}
	}
	return nil
}

func translateTimeOccurence() error {
	if cronSlice[minuteIndex] == anyValue && cronSlice[hourIndex] == anyValue { // checking if both hour and minute are defaults
		translatedString += viper.GetString(configStr + "at_every_minute")
	} else if cronSlice[minuteIndex] != anyValue && cronSlice[hourIndex] != anyValue { //checking if non of them are
		if err := translateMinuteAndHour(); err != nil {
			return err
		}
	} else { // checking if  just one of them is default
		translateMinuteOrHour()
	}
	return nil
}
