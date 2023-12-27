#ifndef CYGNUDGE_JUDGE_TASK_HPP
#define CYGNUDGE_JUDGE_TASK_HPP

#include"../config.h"

#include"judge_pack.hpp"
#include"judge_problem.hpp"
#include"judge_env.hpp"
#include"compile_info.hpp"

#include<vector>
#include<string>
#include<chrono>
#include<ctime>
#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<iostream>
#include<iomanip>
#include<map>
#include<thread>
#include<concepts>
#include<thread>
#include<future>
#include<cstdlib>

#include<boost/functional/hash.hpp>
#include<boost/property_tree/ptree.hpp>
#include<boost/property_tree/json_parser.hpp>

#include<unistd.h>

namespace cyg{
	/*
	hash_code of task is also the temporary directory in /tmp
	hash_code is related to the time client sent the task, user id, and problem id
	*/
	class task{
	private:
		size_t hash_code;
		std::tm time;
		int uid;
		std::string pid;
		std::string code_path;
		std::string program_path;
		std::string compile_msg;
		std::string language;
		//std::vector<std::string> options;
		//not supported yet
		problem problem_info;
		std::string timer;
		std::string task_dir_path;
		bool is_compiled{false};
	public:
		/*
		task.zip name:
		{hash code of the time client sent the task}_{user id}_{problem id}
		e.g.: 2023-12-05_18:28:52_1_P1001.zip
		*/
		task(std::string task_zip_name){
			char pid_buf[32];
			std::sscanf(task_zip_name.c_str(),
				"%d-%d-%d_%d:%d:%d_%d_%s",
				&time.tm_year,&time.tm_mon,&time.tm_mday,&time.tm_hour,&time.tm_min,&time.tm_sec,&uid,pid_buf
			);
			//begin process time & pid
			time.tm_year-=1900;
			time.tm_mon-=1;
			time.tm_isdst=0;
			pid=pid_buf;
			pid.erase(pid.find("."));//requirement: problem id doesn't include '.'
			//end process time & pid
			timer=get_env("timer");//get timer shell from cygnudge_server.json
			std::time_t t=std::mktime(&time);
			hash_code=0;
			boost::hash_combine(hash_code,t);
			boost::hash_combine(hash_code,uid);
			boost::hash_combine(hash_code,pid);

			std::string temp_dir=get_env("temporary_directory");
			task_dir_path=std::format("{}/{}",temp_dir,hash_code);
			std::string zip_path=std::format("{}/{}",temp_dir,task_zip_name);

			unpack(task_dir_path/*dir*/,zip_path/*zip_path*/);
			//todo give a independent function to do clear work
			//unlink(zip_path.c_str());//remove task.zip
			
			//parse json
			using namespace boost::property_tree;
			ptree task_conf_json;
			read_json(std::format("{}/task.json",task_dir_path).c_str(),task_conf_json);
			language=task_conf_json.get<std::string>("language");
			code_path=std::format("{}/code.{}",task_dir_path,language);
			program_path=std::format("{}/program.o",task_dir_path);
			problem_info=problem(pid);
			/*
			full path of `task.zip` is /tmp/cygnudge/{task.zip}
			1. create new directory /tmp/cygnudge/{hash_code}
			2. unzip task.zip into /tmp/cygnudge/{hash_code}
			3. remove task.zip
			4. get task information from /tmp/cygnudge/{hash_code}/task.json (language, pid, uid, time)
			5. get problem information from /var/lib/{pid}/judge.json (problem_info)
			*/
		}

		void log(){
			std::cout
			<<std::format("task: {}\n",hash_code)
			<<std::put_time(&time,"%c %Z\n")
			<<std::format(
				"uid: {}\npid: {}\ntimer: {}\nlanguage: {}\n",
				uid,pid,timer,language
			)
			<<std::flush;
			problem_info.log();
		}

		int compile(){
			/*
			1. program_path: /tmp/cygnudge/{hash_code}/program.o
			2. code_path: /tmp/cygnudge/{hash_code}/code.{language}
			*/
			std::string compile_command=compile_info::gen_command(language,code_path,program_path);
#ifdef CYGNUDGE_DEBUG
			std::cerr<<compile_command<<std::endl;
#endif
			FILE* pipe=popen(compile_command.c_str(),"r");
			char buf;
			while((buf=std::fgetc(pipe))!=EOF)
				compile_msg.push_back(buf);
			is_compiled=true;
			return pclose(pipe);//todo judge CE
		}

		static bool is_correct(std::string out_data_path,std::string ans_path){
			FILE* out_ptr=std::fopen(out_data_path.c_str(),"r");
			FILE* ans_ptr=std::fopen(ans_path.c_str(),"r");
			if(out_ptr==nullptr)
				throw std::runtime_error(std::format("failed to open file: {}",out_data_path));
			if(ans_ptr==nullptr)
				throw std::runtime_error(std::format("failed to open file: {}",ans_path));
			char out_buf,ans_buf;
			std::string out_content,ans_content;
			while((out_buf=std::fgetc(out_ptr))!=EOF)
				out_content.push_back(out_buf);
			while((ans_buf=std::fgetc(ans_ptr))!=EOF)
				ans_content.push_back(ans_buf);
			std::fclose(out_ptr);
			std::fclose(ans_ptr);
			auto clear_end_blank=[](std::string& str){
				size_t lst_line_end_pos=0;
				while(true){
					size_t line_end_pos=str.find('\n',lst_line_end_pos+1);
					if(line_end_pos==std::string::npos)
						break;
					while(line_end_pos>0 and line_end_pos!=std::string::npos){
						if(std::isblank(str[line_end_pos-1]) or std::iscntrl(str[line_end_pos-1]))
							str.erase(line_end_pos-1,1);
						else
							break;
						line_end_pos=str.find('\n');
					}
					lst_line_end_pos=line_end_pos;
				};
				while(std::isblank(str.back()) or std::iscntrl(str.back()) or str.back()=='\n')
					str.pop_back();
			};
			clear_end_blank(out_content);
			clear_end_blank(ans_content);
#ifdef CYGNUDGE_DEBUG
			std::cerr<<out_data_path<<": "<<std::endl
			<<out_content<<std::endl
			<<ans_path<<": "<<std::endl
			<<ans_content<<std::endl;
#endif
			return out_content==ans_content;
		}

		/*static bool popen_with_timeout(const char* command,size_t timeout){
			pid_t pid=fork();
			if(pid==0){//subprocess
				std::system(command);
				exit(EXIT_SUCCESS);
			}
			//std::this_thread::sleep_for(std::chrono::milliseconds(timeout));
			if(waitpid(pid,NULL,WNOHANG)==0){//still running
				std::this_thread::sleep_for(std::chrono::milliseconds(timeout));
				kill(pid,SIGKILL);
				return true;
			}
			else
				return false;
		}*///popen_with_timeout(judge_command.c_str(),500);
		//return: true=>timeout false=>exit

		static bool popen_with_timeout(const char* command,size_t timeout/*ms*/){
			std::promise<void> p;
			std::future<void> f=p.get_future();
			std::thread t([&command](std::promise<void>&& _p){
#ifdef CYGNUDGE_DEBUG
				std::cerr<<command<<std::endl;
#endif
				std::system(command);
			},std::move(p));
			t.detach();
			if(f.wait_for(std::chrono::milliseconds(timeout))==std::future_status::timeout)
				return true;
			else
				return false;
		}

		void judge_single_point(size_t subtask,size_t subpoint){
			if(not is_compiled)
				throw std::runtime_error("the program hasn't been compiled");

			auto [in_data_path,out_data_path]=problem_info.get_data_point(subtask,subpoint);
			std::string ans_path{std::format("{}/{}:{}.ans",task_dir_path,subtask,subpoint)};
			std::string state_path{std::format("{}/{}:{}_state",task_dir_path,subtask,subpoint)};
			std::string judge_command{
				std::format(
					"cat {} | {} -f \"%e %M %x\" -o {} {} > {}",
					in_data_path,timer,state_path,program_path,ans_path
				)
			};
			point_result& judge_point=problem_info.result[subtask][subpoint];

			bool is_timeout=popen_with_timeout(judge_command.c_str(),
				problem_info.subtasks[subtask][subpoint].time+CYGNUDGE_JUDGE_DELAY);
			//running within timeout
			if(is_timeout){
				judge_point.status=4;//TLE
				return;
			}

			double judge_time/*second*/;
			int judge_memory/*KiB*/;
			int judge_return_code;
			FILE* state_ptr=std::fopen(state_path.c_str(),"r");
			if(state_ptr==nullptr)
				throw std::runtime_error(std::format("failed to open file: {}",state_path));
			std::fscanf(state_ptr,"%lf %d %d",&judge_time,&judge_memory,&judge_return_code);
			judge_point.time=judge_time*1000;
			judge_point.memory=judge_memory/1024;
			judge_point.return_code=judge_return_code;
			std::fclose(state_ptr);

			if(judge_return_code!=0){
				judge_point.status=3;//RE
			}
			else if(judge_point.memory>problem_info.subtasks[subtask][subpoint].memory){
				judge_point.status=5;//MLE
			}
			else{
				bool judge_state=is_correct(out_data_path,ans_path);
				if(judge_state)
					judge_point.status=0;//AC
				else
					judge_point.status=1;//WA
			}
		}

		void judge(){
			if(int ec=compile();ec!=0){
				for(size_t subtask=0;subtask<problem_info.result.size();++subtask)
					for(size_t subpoint=0;subpoint<problem_info.result[subtask].size();++subpoint)
						problem_info.result[subtask][subpoint].status=2;//CE
				problem_info.final_status=2;//CE
				return;
			}
#ifdef CYGNUDGE_DEBUG
			std::cerr<<"compilation has finished"<<std::endl;
#endif
			std::map<size_t,size_t> status_count;/*status_code=>count*/
			for(size_t subtask=0;subtask<problem_info.result.size();++subtask){
				for(size_t subpoint=0;subpoint<problem_info.result[subtask].size();++subpoint){
					judge_single_point(subtask,subpoint);
					status_count[problem_info.result[subtask][subpoint].status]++;
					if(problem_info.result[subtask][subpoint].status==0)
						problem_info.score_get+=problem_info.subtasks[subtask][subpoint].score;
				}
			}
			if(problem_info.score_get==0){
				for(const auto& count:status_count){
					if(count.second==problem_info.point_sum_n){
						problem_info.final_status=count.first;//completely a certain status
					}
				}
				problem_info.final_status=1;//WA
			}
			else if(problem_info.score_get==problem_info.full_score){
				problem_info.final_status=0;//AC
			}
			else{
				problem_info.final_status=6;//PAC
			}
			//cat {in_data_path} | {timer} -f "%e %M %x" {program_path} > {ans_file_path}
		}

		void log_judge(){
			problem_info.log_result();
		}

		void export_result_json(std::string result_json_path){
			problem_info.export_result_json(result_json_path);
		}
	};
}

#endif
