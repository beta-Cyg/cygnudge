#ifndef CYGNUDGE_JUDGE_PROBLEM_HPP
#define CYGNUDGE_JUDGE_PROBLEM_HPP

#include<boost/property_tree/ptree.hpp>
#include<boost/property_tree/json_parser.hpp>

#include<string>
#include<vector>
#include<format>
#include<fstream>
#include<algorithm>

#include"judge_env.hpp"

#include<iostream>

namespace cyg{
	struct point{
		size_t time/*ms*/,memory/*MiB*/,score/*pts*/;
	};

	constexpr static std::array<std::string,7> status_name{"AC"/*0*/,"WA"/*1*/,"CE"/*2*/,"RE"/*3*/,"TLE"/*4*/,"MLE"/*5*/,"PAC"/*6*/};

	struct point_result{
		size_t time,memory,return_code,status;

		point_result(size_t _t=0,size_t _m=0,size_t _rc=0,size_t _s=0):time{_t},memory{_m},return_code{_rc},status{_s}{}
	};
 
	struct problem{
		std::string pid{};
		std::vector<std::vector<point>> subtasks{};
		int subtask_n{},point_sum_n{},full_score{};

		std::vector<std::vector<point_result>> result;
		size_t final_status{};
		size_t score_get{};

		problem()=default;

		problem(std::string _pid):pid{_pid},point_sum_n{0},full_score{0},score_get{0}{
			using namespace boost::property_tree;
			ptree task_json;
			read_json(
				std::format("{}/{}/judge.json"/*problem_dir, pid*/,get_env("problem_directory"),_pid),
				task_json
			);
			subtask_n=task_json.get<int>("subtask");
			for(int i=0;i<subtask_n;++i){
				std::vector<point> subtask;
				int point_n=task_json.get<int>(std::format("s{}.point",i));
				point_sum_n+=point_n;
				for(int j=0;j<point_n;++j){
					point tmp;
					tmp.time=task_json.get<int>(std::format("s{}.p{}.time",i,j));
					tmp.memory=task_json.get<int>(std::format("s{}.p{}.memory",i,j));
					tmp.score=task_json.get<int>(std::format("s{}.p{}.score",i,j));
					full_score+=tmp.score;
					subtask.push_back(tmp);
				}
				subtasks.push_back(subtask);
			}
			result.resize(subtask_n);
			for(size_t i=0;i<result.size();++i)
				result[i].resize(subtasks[i].size());
		}

		void export_result_json(std::string result_json_path){
			using namespace boost::property_tree;
			ptree result_json;
			if(final_status==2){//CE
				result_json.put("score",0);
				result_json.put("status","CE");
			}
			else{//non-CE
				result_json.put("score",score_get);
				result_json.put("status",status_name[final_status]);
				for(size_t subtask=0;subtask<result.size();++subtask){
					ptree subtask_json;
					subtask_json.put("point",result[subtask].size());
					for(size_t subpoint=0;subpoint<result[subtask].size();++subpoint){
						ptree subpoint_json;
						subpoint_json.put("time",result[subtask][subpoint].time);
						subpoint_json.put("memory",result[subtask][subpoint].memory);
						subpoint_json.put("return_code",result[subtask][subpoint].return_code);
						subpoint_json.put("status",status_name[result[subtask][subpoint].status]);
						subtask_json.add_child(std::format("p{}",subpoint),subpoint_json);
					}
					result_json.add_child(std::format("s{}",subtask),subtask_json);
				}
			}
			write_json(result_json_path,result_json);
		}

		void log(){
			std::cout<<std::format(
				"problem: {}\ntotal subtask: {}\ntotal point: {}\ntotal score: {}\n",
				pid,subtask_n,point_sum_n,full_score
			);
			for(size_t i=0;i<subtask_n;++i){
				std::cout<<std::format("subtask {}:\n",i);
				for(size_t j=0;j<subtasks[i].size();++j){
					std::cout<<std::format(
						"\tpoint {}:\n\t\ttime: {}ms\n\t\tmemory: {}MiB\n\t\tscore: {}pts\n",
						j,subtasks[i][j].time,subtasks[i][j].memory,subtasks[i][j].score
					);
				}
			}
		}

		void log_result(){
			std::cout<<std::format(
				"judge result:\nproblem: {}\nstatus: {}\nscore: {}\n",
				pid,status_name[final_status],score_get
			)<<std::flush;
			for(size_t i=0;i<subtask_n;++i){
				std::cout<<std::format("subtask {}:\n",i)<<std::flush;
				for(size_t j=0;j<result[i].size();++j){
					std::cout<<std::format(
						"\tpoint {}:\n\t\ttime {}ms\n\t\tmemory: {}MiB\n\t\tstatus: {}\n\t\tscore: {}pts\n",
						j,result[i][j].time,result[i][j].memory,status_name[result[i][j].status],result[i][j].status==0?subtasks[i][j].score:0
					);
				}
			}
		}

		std::pair<std::string,std::string> get_data_point(size_t subtask,size_t subpoint){
			return {
				std::format("{}/{}/data/{}:{}.in",get_env("problem_directory"),pid,subtask,subpoint),
				std::format("{}/{}/data/{}:{}.out",get_env("problem_directory"),pid,subtask,subpoint)
			};
		}
	};
}

#endif
