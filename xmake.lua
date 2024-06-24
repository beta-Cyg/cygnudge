target("problem_test")
	add_files("src/test/problem_test.cpp")
	set_kind("binary")
	set_languages("c++20")
	add_cxxflags("-DCYGNUDGE_DEBUG")

target("task_test")
	add_files("src/test/task_test.cpp")
	set_kind("binary")
	set_languages("c++20")
	add_cxxflags("-DCYGNUDGE_DEBUG")

target("thread_pool_test")
	add_files("src/test/thread_pool_test.cpp")
	set_kind("binary")
	set_languages("c++20")
	add_cxxflags("-DCYGNUDGE_DEBUG")

target("gen_data_P1001")
	add_files("src/test/gen_data_P1001.cpp")
	set_kind("binary")
