#include"../server/judge_task.hpp"

int main(){
	cyg::task t("2023-12-14_21:33:00_1_P1001.zip");
	t.log();
	cyg::compile_info::import_compile_commands(cyg::get_env("compile_json"));
	t.judge();
	t.log_judge();
	t.export_result_json("P1001_result.json");

	return 0;
}
