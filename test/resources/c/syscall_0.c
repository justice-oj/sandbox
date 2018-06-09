#include <stdio.h>
#include <sys/signal.h>

int main(void) {
    kill(1, SIGSEGV);
    return 0;
}