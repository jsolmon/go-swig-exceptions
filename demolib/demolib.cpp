#include "demolib.h"
#include <cstdlib>
#include <stdexcept>

DemoLib::DemoLib() {}

double DemoLib::DivideBy(int n) {
    if (n == 0) {
        throw std::invalid_argument("Cannot divide by zero");
    }
    return 1.0 / n;
}

int DemoLib::NegativeThrows(int in) {
    if (in < 0) {
        throw std::range_error("NegativeThrows threw exception");
    }
    return in;
}

int DemoLib::NeverThrows(int in) {
    return in;
}
