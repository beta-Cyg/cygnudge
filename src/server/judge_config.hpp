#ifndef CYGNUDGE_JUDGE_CONFIG_HPP
#define CYGNUDGE_JUDGE_CONFIG_HPP

#include<vector>
#include<string>
#include<format>

/*
 * client data:
 * user id
 * problem id
 * task code
 * compilation config
 * config includes:
 * language id
 * optimization option
*/

namespace cyg{
	struct compilation_config{
		std::string command_format;
		std::vector<std::string> option_supported;

		compilation_config(const std::string& cf,std::initializer_list<std::string> os):command_format{cf},option_supported{os}{}

		std::string gen(std::string code_file,std::string program_file,std::vector<std::string> options);//todo
	};

	//command_format: {0} -> source code {1} -> program {2} -> options

	struct execution_config{
		std::string command_format;

		execution_config(const std::string& cf):command_format{cf}{}

		std::string gen(std::string program_file){
			return std::vformat(command_format,program_file);
		}
	};

	//command_format: {0} -> program

	compilation_config get_compilation_config(std::string language);//todo

	execution_config get_execution_config(std::string language);//todo
}

#endif
