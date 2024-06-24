package main

import (
	"cygnudge/judge"
)

func main() {
	var t judge.Task
	t = judge.NewTask("2023-12-05_18:28:52_1_P1001.zip")
	t.Log()
	judge.ImportCompileCommands()
	judge.PrintCompileCommands()
	t.Judge()
	t.LogJudge()
	t.ExportResultJson("P1001_result_go.json")
	t.Free()
}
