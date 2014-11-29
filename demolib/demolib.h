#ifndef DEMOLIB_H
#define DEMOLIB_H

class DemoLib {
    public:
        DemoLib();

        // throws invalid_argument exception if arg = 0
        double DivideBy(int); 
        // throws range_error exception if arg < 0
        int NegativeThrows(int);

        // does not throw any exceptions
        int NeverThrows(int);
};

#endif
