#ifndef CYGNUDGE_JUDGE_ENV_HPP
#define CYGNUDGE_JUDGE_ENV_HPP

#include<boost/property_tree/ptree.hpp>
#include<boost/property_tree/json_parser.hpp>

#include"../config.h"

#include<string>
#include<fstream>

namespace cyg{
	std::string get_env(std::string var){
		using namespace boost::property_tree;
		ptree server_conf_json;
		read_json(CYGNUDGE_SERVER_JSON,server_conf_json);
		return server_conf_json.get<std::string>(var);
	}
}

#endif
