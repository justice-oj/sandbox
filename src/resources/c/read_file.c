#include <stdio.h>

int main() {
    int ch;
    FILE *file;

    file = fopen("/etc/hosts", "r");
    if (file) {
        while ((ch = getc(file)) != EOF)
            putchar(ch);
        fclose(file);
    } else {
        printf("File not found");
    }

    return 0;
}