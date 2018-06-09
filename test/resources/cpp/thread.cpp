#include <string>
#include <iostream>
#include <thread>

using namespace std;

void task() {
    cout << "invoke in threads.\n";
}

int main() {
    thread t(task);
    t.join();
}