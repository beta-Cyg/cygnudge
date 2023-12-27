package judge

import "unsafe"

// #cgo CFLAGS: -DCYGNUDGE_DEBUG
// #cgo LDFLAGS: -L/home/beta-cyg/cygnudge/lib -ljudge_interface -lstdc++ -lm
// #include"/home/beta-cyg/cygnudge/src/server/judge_interface.h"
// #include"stdlib.h"
import "C"

type Task struct{
	task C.Task
}

func NewTask (task_zip_name string) Task{
	var ret Task
	tmp:=C.CString(task_zip_name)
	ret.task=C.NewTask(tmp)
	C.free(unsafe.Pointer(tmp))
	return ret
}

func (t *Task) Log(){
	C.TaskLog(t.task)
}

func (t *Task) Judge(){
	C.TaskJudge(t.task)
}

func (t *Task) LogJudge(){
	C.TaskLogJudge(t.task)
}

func (t *Task) ExportResultJson(result_json_name string){
	tmp:=C.CString(result_json_name)
	C.TaskExportResultJson(t.task,tmp)
	C.free(unsafe.Pointer(tmp))
}

func (t *Task) Free(){
	C.DeleteTask(t.task)
}

func ImportCompileCommands(){
	C.ImportCompileCommands()
}

func PrintCompileCommands(){
	C.PrintCompileCommands()
}
