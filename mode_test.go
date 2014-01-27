package goirc

import (
	"testing"
)

func TestMode(t *testing.T) {
	bot := create_bot_with_chan(4)
	names := []string{"", "", "#Test2", "manacit @manacit2 +manacit3 manacit4 manacit5"}
	bot.Names(names)
	mode_string := []string{"#Test2", "+v", "manacit"}
	bot.ParseMode(mode_string)
	test_chan, _ := bot.GetChan("#Test2")
	for _, v := range test_chan.Users {
		if v.Name == "manacit" && v.Voice == false {
			t.Error("manacit should be voiced")
		}
	}
	mode_string = []string{"#Test2", "-v", "manacit"}
	bot.ParseMode(mode_string)
	for _, v := range test_chan.Users {
		if v.Name == "manacit" && v.Voice == true {
			t.Error("manacit should not be voiced")
		}
	}
	mode_string = []string{"#Test2", "-o+v", "manacit2", "manacit"}
	bot.ParseMode(mode_string)
	for _, v := range test_chan.Users {
		if v.Name == "manacit" && v.Voice == false {
			t.Error("manacit should be voiced")
		}
		if v.Name == "manacit2" && v.Op == true {
			t.Error("manacit2 should not be voiced")
		}
	}
	mode_string = []string{"#Test2", "+vv", "manacit4", "manacit5"}
	bot.ParseMode(mode_string)
	for _, v := range test_chan.Users {
		if (v.Name == "manacit4" || v.Name == "manacit5") && v.Voice == false {
			t.Error("manacit should be voiced")
		}
	}
}
