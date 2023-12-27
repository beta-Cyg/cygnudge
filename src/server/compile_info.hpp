#ifndef CYGNUDGE_COMPILE_INFO_HPP
#define CYGNUDGE_COMPILE_INFO_HPP

#include<stdexcept>
#include<string>
#include<map>
#include<format>
#include<iostream>

#include<boost/property_tree/ptree.hpp>
#include<boost/property_tree/json_parser.hpp>

namespace cyg{
	class compile_info{
	private:
		static std::map<std::string,std::string> compile_commands;/*language=>command format*/
		static bool is_imported;
	public:
		static void import_compile_commands(std::string compile_json_path){//compile_json is in get_env("")
			using namespace boost::property_tree;
			ptree compile_json;
			read_json(compile_json_path,compile_json);
			for(const auto& pair:compile_json){
				compile_commands.emplace(pair.first,pair.second.get<std::string>(""));
			}
			is_imported=true;
		}

		static void print(){
			if(not is_imported)
				throw std::runtime_error("compile commands hasn't been imported yet");
			for(const auto& i:compile_commands)
				std::cout<<i.first<<" : "<<i.second<<std::endl;
		}

		static std::string/*compile_command*/ gen_command(std::string language,std::string code_file_name,std::string program_name){
			if(not is_imported)
				throw std::runtime_error("compile commands hasn't been imported yet");
			if(compile_commands.count(language)==0)
				throw std::invalid_argument(std::format("undefined language: {}",language));
			return std::vformat(compile_commands[language],std::make_format_args(code_file_name,program_name));
		}
	};

	std::map<std::string,std::string> compile_info::compile_commands{};
	bool compile_info::is_imported{false};
}

#endif
