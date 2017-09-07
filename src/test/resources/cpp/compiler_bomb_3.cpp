template<int A, int B>

struct a {
    static const int n;
};

template<int A, int B> const int a<A, B>::n = a<A - 1, a<A, B - 1>::n>::n;

template<int A>

struct a<A, 0> {
    static const int n = a<A - 1, 1>::n;
};

template<int B>

struct a<0, B> {
    static const int n = B + 1;
};

int h = a<4, 2>::n;