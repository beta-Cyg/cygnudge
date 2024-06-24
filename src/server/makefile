CXX = g++
OPT = -std=c++20

cygnudge_judge: compile_info.hpp judge_config.hpp judge_env.hpp judge_pack.hpp judge_problem.hpp judge_task.hpp judge.hpp thread_pool.hpp judge_interface.cpp judge_interface.h
	mkdir -p ../../lib
	$(CXX) $(OPT) -fPIC -c judge_interface.cpp -o ../../lib/cygnudge_judge.o
	$(CXX) -shared -o ../../lib/libcygnudge_judge.so ../../lib/cygnudge_judge.o