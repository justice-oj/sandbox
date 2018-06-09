#define Z struct X
#define T template<int N
T,int M=N>Z;struct Y{static int f(){return 0;}};T>Z<N,0>:Y{};T>Z<0,N>:Y{};T,int M>Z{static int f(){static int x[99999]={X<N-1,M>::f()+X<N,M-1>::f()};}};int x=X<80>::f();