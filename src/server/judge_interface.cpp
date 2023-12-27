#include"judge_task.hpp"
#include"judge_interface.h"

const char* GetEnv(const char* var){
	std::string tmp=cyg::get_env(var);
	char* ret=static_cast<char*>(std::malloc((tmp.size()+1)*sizeof(char)));
	std::strcpy(ret,tmp.c_str());
	return ret;
}

Task NewTask(const char* task_zip_name){
	cyg::task* ret=new cyg::task(task_zip_name);
	return static_cast<Task>(ret);
}

void DeleteTask(Task t){
	cyg::task* task=static_cast<cyg::task*>(t);
	delete task;
}

void TaskLog(Task t){
	cyg::task* task=static_cast<cyg::task*>(t);
	task->log();
}

void TaskJudge(Task t){
	cyg::task* task=static_cast<cyg::task*>(t);
	task->judge();
}

void TaskLogJudge(Task t){
	cyg::task* task=static_cast<cyg::task*>(t);
	task->log_judge();
}

void TaskExportResultJson(Task t,const char* result_json_path){
	cyg::task* task=static_cast<cyg::task*>(t);
	task->export_result_json(result_json_path);
}

void ImportCompileCommands(){
	cyg::compile_info::import_compile_commands(cyg::get_env("compile_json"));
}

void PrintCompileCommands(){
	cyg::compile_info::print();
}
