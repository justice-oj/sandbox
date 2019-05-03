#include <cmath>
#include <iostream>
#include <iomanip>

using namespace std;

int main(){
    int h, m, s;
    char ch, aorp;

    cin >> h >> ch >> m >> ch >> s >> aorp >> ch;
    h = (aorp == 'A') ? (h==12 ? 0 : h) : (h==12 ? 12 : h+12);

    cout << setw(2) << setfill('0') << h << ":"
         << setw(2) << setfill('0') << m << ":"
         << setw(2) << setfill('0') << s << endl;

    return 0;
}