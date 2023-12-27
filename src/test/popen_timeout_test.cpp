#include <unistd.h>
#include <signal.h>
#include <sys/wait.h>
#include <iostream>
#include <cstdio>

pid_t popen2(const char *command) {
    pid_t pid = fork();

    if (pid == 0) { // 子进程
        system(command);
        exit(EXIT_SUCCESS);
    }

    return pid; // 父进程返回子进程的 PID
}

int main() {
    const char *command = "python3 /home/beta-cyg/cygnudge/timeout_test.py"; // 你的命令
    int timeout = 2; // 超时时间（秒）

    pid_t pid = popen2(command);

    sleep(timeout);

    if (waitpid(pid, NULL, WNOHANG) == 0) { // 如果子进程还在运行
        std::cout << "超时，终止进程\n";
        kill(pid, SIGKILL); // 终止子进程
    } else {
        std::cout << "进程已完成\n";
    }

    return 0;
}

