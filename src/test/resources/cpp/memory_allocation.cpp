int main() {
    for (;;) {
        int *x = new int[100000000];
        x[0] = 0;
        x[100000000 - 1] = 123;
    }
}