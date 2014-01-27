package goirc

import (
	"testing"
)

func TestJoinUser(t *testing.T) {
	bot := create_bot_with_chan(3)
	bot.JoinUser("#Test1", "manacit")
	test_chan, _ := bot.GetChan("#Test1")
	if test_chan.Users[0].Name != "manacit" && len(test_chan.Users) != 1 {
		t.Error("There should be one user named manacit")
	}
}

func TestPartUser(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	bot.JoinUser("#Test1", "manacit")
	if test_chan.Users[0].Name != "manacit" && len(test_chan.Users) != 1 {
		t.Error("There should be one user named manacit")
	}
	bot.PartUser("#Test1", "manacit")
	if len(test_chan.Users) > 0 {
		t.Error("There should be 0 users")
	}
}

func TestQuitUser(t *testing.T) {
	bot := create_bot_with_chan(3)
	bot.JoinUser("#Test1", "manacit")
	bot.JoinUser("#Test0", "manacit")
	bot.JoinUser("#Test2", "manacit")
	for _, v := range bot.Channels {
		if v.Users[0].Name != "manacit" {
			t.Error("manacit should be the first person in the channel")
		}
	}
	bot.QuitUser("manacit")
	for _, v := range bot.Channels {
		if len(v.Users) > 0 {
			t.Error("There should be no users in the channel")
		}
	}
}

func TestUserJoin(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.UserJoin("manacit")
	if test_chan.Users[0].Name != "manacit" || len(test_chan.Users) > 1 {
		t.Error("The only user in the channel should be manacit")
	}
	if test_chan.Users[0].Op == true || test_chan.Users[0].Voice == true {
		t.Error("manacit should be an op only")
	}
}

func TestUserOpJoin(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.UserJoin("@manacit")
	if test_chan.Users[0].Name != "manacit" || len(test_chan.Users) > 1 {
		t.Error("The only user in the channel should be manacit")
	}
	if test_chan.Users[0].Op != true || test_chan.Users[0].Voice == true {
		t.Error("manacit should be an op only")
	}
}

func TestVoiceUserJoin(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.UserJoin("+manacit")
	if test_chan.Users[0].Name != "manacit" || len(test_chan.Users) > 1 {
		t.Error("The only user in the channel should be manacit")
	}
	if test_chan.Users[0].Voice != true || test_chan.Users[0].Op == true {
		t.Error("manacit should be an op only")
	}
}

func TestUserDel(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.UserJoin("manacit")
	if test_chan.Users[0].Name != "manacit" || len(test_chan.Users) > 1 {
		t.Error("The only user in the channel should be manacit")
	}
	test_chan.UserDel("manacit")
	if len(test_chan.Users) > 0 {
		t.Error("manacit shouldn't be in the channel")
	}
	test_chan.UserJoin("manacit")
	test_chan.UserJoin("manacit2")
	test_chan.UserJoin("manacit3")
	test_chan.UserJoin("manacit4")
	test_chan.UserJoin("manacit5")
	test_chan.UserJoin("manacit6")
	test_chan.UserDel("manacit")
	if len(test_chan.Users) > 5 {
		t.Error("manacit shouldn't be in the channel")
	}
}

func TestGetUser(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.UserJoin("manacit")
	us := test_chan.GetUser("manacit")
	if us.Name != "manacit" {
		t.Error("manacit should have been found")
	}

	us = test_chan.GetUser("nope")
	if us.Name != "" {
		t.Error("This user does not exist!")
	}
}

func TestBotGetUser(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.UserJoin("manacit")
	us, er := bot.GetUser("#Test1", "manacit")
	if er != nil {
		t.Error("Looking for manacit should not raise an error")
	}
	if us.Name != "manacit" {
		t.Error("manacit should have been found")
	}

	us, er = bot.GetUser("#Test1", "nope")
	if er == nil {
		t.Error("This user does not exist!")
	}
	us, er = bot.GetUser("#Blah", "manacit")
	if er == nil {
		t.Error("That user does not exist")
	}
}

func TestRenameUser(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.UserJoin("manacit")
	bot.RenameUser("manacit", "manacit2")
	if test_chan.Users[0].Name != "manacit2" {
		t.Error("manacit2 sould be the only user in the channel")
	}
}

func TestInChan(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.UserJoin("manacit")
	if test_chan.InChan("manacit") == false {
		t.Error("manacit is in the channel")
	}
	if test_chan.InChan("notinchan") == true {
		t.Error("InChan should not return true when someone is not in the channel")
	}
}

func TestAddUser(t *testing.T) {
	bot := create_bot_with_chan(3)
	test_chan, _ := bot.GetChan("#Test1")
	test_chan.AddUser("manacit")
	if test_chan.Users[0].Name != "manacit" || len(test_chan.Users) != 1 {
		t.Error("manacit sould be the only user in the channel")
	}
}
