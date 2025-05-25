#include <emscripten/emscripten.h>

extern "C" {

inline double f(double x) { 
    return x * x; 
}

EMSCRIPTEN_KEEPALIVE
double x2Integrate(double xmin, double xmax, int intervals_count) {
    double dx = (xmax - xmin) / static_cast<double>(intervals_count);
    double total = 0.0;
    double x = xmin;
    for (int i = 0; i < intervals_count; i++) {
        total += dx * (f(x) + f(x + dx)) / 2.0;
        x += dx;
    }
    return total;
}

EMSCRIPTEN_KEEPALIVE
double x2IntegrateMock(double, double, int) {
    return 0.0;
}

EMSCRIPTEN_KEEPALIVE
int fibonacciRecursive(int n) {
    if (n <= 1) return n;
    return fibonacciRecursive(n - 1) + fibonacciRecursive(n - 2);
}

EMSCRIPTEN_KEEPALIVE
int fibonacciRecursiveMock(int) {
    return 0;
}

EMSCRIPTEN_KEEPALIVE
int fibonacciIterative(int n) {
    if (n <= 1) return n;
    int a = 0, b = 1, tmp;
    for (int i = 2; i <= n; i++) {
        tmp = a + b;
        a = b;
        b = tmp;
    }
    return b;
}

EMSCRIPTEN_KEEPALIVE
int fibonacciIterativeMock(int) {
    return 0;
}

EMSCRIPTEN_KEEPALIVE
int multiply(int size) {
    int a = 33, b = 10, result = 0;
    for (int i = 0; i < size; i++) {
        result = a * b;
    }
    return result;
}

EMSCRIPTEN_KEEPALIVE
int multiplyMock(int) {
    return 0;
}

EMSCRIPTEN_KEEPALIVE
void multiplyVector(int size) {
    int a = 33, b = 10;
    int* a_vector = new int[size];
    int* b_vector = new int[size];
    int* result_vector = new int[size];

    for (int i = 0; i < size; i++) {
        a_vector[i] = a;
        b_vector[i] = b;
    }
    for (int i = 0; i < size; i++) {
        result_vector[i] = a_vector[i] * b_vector[i];
    }

    delete[] a_vector;
    delete[] b_vector;
    delete[] result_vector;
}

EMSCRIPTEN_KEEPALIVE
void multiplyVectorMock(int) {}

EMSCRIPTEN_KEEPALIVE
void factorize(int n) {
    int d = 2;
    while (n > 1) {
        while (n % d == 0) {
            n /= d;
        }
        d += 1;
    }
}

EMSCRIPTEN_KEEPALIVE
void factorizeMock(int) {}

}
