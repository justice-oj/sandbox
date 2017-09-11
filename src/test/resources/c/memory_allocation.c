#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

int main(void) {
    int chunk_size = 1024 * 1024;
    void *p = NULL;

    while(1) {
        if ((p = malloc(chunk_size*sizeof(float))) == NULL) {
            printf("Out of memory");
            break;
        }
        memset(p, 1, chunk_size);
    }
    return 0;
}