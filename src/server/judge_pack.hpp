#ifndef CYGNUDGE_JUDGE_PACK_HPP
#define CYGNUDGE_JUDGE_PACK_HPP

#include"judge_env.hpp"
#include"../config.h"

#include<cstdlib>
#include<string>
#include<format>

#include<unistd.h>
#include<sys/stat.h>
#include<sys/types.h>

namespace cyg{
	int pack(std::string dir,std::string zip_name){
		return std::system(std::format("{} -z {} -d {}",get_env("cygpack"),zip_name,dir).c_str());
	}

	int unpack(std::string dir,std::string zip_name){
		if(access(dir.c_str(),F_OK)>=0);//the directory has existed
		else{
			mkdir(dir.c_str(),0777);
		}
		return std::system(std::format("{} -u {} -d {}",get_env("cygpack"),zip_name,dir).c_str());
	}
}

#endif
