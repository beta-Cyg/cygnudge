#ifndef CYGNUDGE_JUDGE_INTERFACE_H
#define CYGNUDGE_JUDGE_INTERFACE_H

#ifdef __cplusplus
extern "C"{
#endif
	typedef void* Task;

	const char* GetEnv(const char*);

	Task NewTask(const char*);

	void DeleteTask(Task);

	void TaskLog(Task);

	void TaskJudge(Task);

	void TaskLogJudge(Task);

	void TaskExportResultJson(Task,const char*);

	void ImportCompileCommands();

	void PrintCompileCommands();
#ifdef __cplusplus
}
#endif

#endif
