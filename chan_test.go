package goirc

import (
	"fmt"
	"testing"
)

func TestJoin(t *testing.T) {
	bot := Bot{}
	bot.Join("#Test")
	if len(bot.Channels) != 1 {
		t.Error("There should be one channel")
	}
	if bot.Channels[0].Name != "#Test" {
		t.Error("The first channel should be #Test")
	}
	bot.Join("#Test2")
	if len(bot.Channels) != 2 {
		t.Error("There should be two channels")
	}
}

func TestPart(t *testing.T) {
	bot := create_bot_with_chan(3)
	bot.Part("#Test1")
	if len(bot.Channels) != 2 {
		t.Error("There should be two channels")
	}
	for _, v := range bot.Channels {
		if v.Name != "#Test0" && v.Name != "#Test2" {
			t.Error("The wrong channel is there!")
		}
	}
}

func TestNames(t *testing.T) {
	names := []string{"", "", "#Test", "manacit @manacit2 +manacit3"}
	bot := Bot{}
	bot.Join("#Test")
	bot.Names(names)
	if len(bot.Channels[0].Users) != 3 {
		t.Error("There should be three users in the channel")
	}
	for _, v := range bot.Channels[0].Users {
		if v.Name != "manacit" && v.Name != "manacit2" && v.Name != "manacit3" {
			t.Error("A user is there that should not be there!")
		}
		if v.Name == "manacit2" && v.Op == false {
			t.Error("manacit2 should be an op!")
		}
		if v.Name == "manacit3" && v.Voice == false {
			t.Error("manacit4 should be voiced!")
		}
	}
}

func TestGetChan(t *testing.T) {
	bot := create_bot_with_chan(5)
	_, e := bot.GetChan("#Test3")
	if e != nil {
		t.Error("The channel should be there!")
	}
	_, ne := bot.GetChan("#YOLO")
	if ne == nil {
		t.Error("There should be an error - this channel doesn't exist")
	}
}

func TestChanList(t *testing.T) {
	bot := create_bot_with_chan(3)
	if bot.GetChanList() != "#Test0,#Test1,#Test2" {
		t.Error("The chan list is incorrect")
	}
}

func create_bot_with_chan(n int) Bot {
	bot := Bot{}
	for i := 0; i < n; i++ {
		to_join := fmt.Sprintf("#Test%d", i)
		bot.Join(to_join)
	}
	return bot
}
