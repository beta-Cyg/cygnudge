{
    files = {
        "src/test/popen_timeout_test.cpp"
    },
    depfiles_gcc = "popen_timeout_test.o: src/test/popen_timeout_test.cpp\
",
    values = {
        "/usr/bin/gcc",
        {
            "-m64",
            "-std=c++20",
            "-DCYGNUDGE_DEBUG"
        }
    }
}