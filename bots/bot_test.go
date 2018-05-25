package bots

import "testing"

func TestRunCmdBot(t *testing.T) {
	t.Log(RunCmdBot("","ls","-lh"))

	t.Log(RunCmdBot("/","ls","-lh"))
}

func TestRunTulingTalkBot(t *testing.T) {
	t.Log(RunTulingTalkBot("mm","你好啊"))
}

func TestRunMoliTalkBot(t *testing.T) {
	t.Log(RunMoliTalkBot("mm","你好啊"))
}

func TestRunYoudaoQueryBot(t *testing.T) {
	t.Log(RunYoudaoQueryBot("sex"))
	t.Log(RunYoudaoQueryBot("boy","is","good"))
}
