## Using SWIG to port a C++ library with exceptions to Go
This code is intended as an example of how to port a C++ library which throws exceptions to go. All methods which throw exceptions are wrapped in try/catch blocks which transform the exceptions into panics which are then recovered from to return go errors.

#### Install
From the demolib directory, run
```Bash
go install
```
#### Test
From the base directory, run
```Bash
go test
```

## Description

### DemoLib C++ Class
The DemoLib C++ class has two methods (`DivideBy` and `NegativeThrows`) which throw exceptions on some inputs and one method (`NeverThrows`) which will never throw an exception.
```C++
DemoLib::DemoLib() {}o

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
```

### .swigcxx File
The `%exception` directive will wrap all calls to the C++ library in a try/catch block and turn any caught exceptions into go panics containing the C++ `exception.what()` message. The `_swig_gopanic()` method is defined in `exception.i` in the swig library.
```
%include "exception.i"
%exception {
    try {
        $action;
    } catch (std::exception &e) {
        _swig_gopanic(e.what());
    }
}
```

The go interface and methods created by swig must be renamed so that they can be wrapped by the code which catches the panics and turns the result into errors. The `DemoLib` class is renamed to `Wrapped_DemoLib` and the `NegativeThrows` and `DivideBy` methods are wrapped as `Wrapped_NegativeThrows` and `Wrapped_DivideBy`. No wrapping is necessary for the `NeverThrows` method because it does not throw exceptions and will not be changed.
```
%rename(Wrapped_DemoLib) DemoLib;
%rename(Wrapped_NegativeThrows) DemoLib::NegativeThrows;
%rename(Wrapped_DivideBy) DemoLib::DivideBy;
```

The `go_wrapper` allows for writing go code in the swig file to wrap the interfaces and methods produced by swig. Prior to entering the `go_wrapper` you can import any needed go libraries with the `%go_import` directive.
```
%go_import("fmt")
%insert(go_wrapper) %{
```

Inside the `go_wrapper`, create a wrapper interface called DemoLib which includes functions for `NegativeThrows` and `DivideBy` which return tuples which include error types.
```
type DemoLib interface {
    Wrapped_DemoLib
    NegativeThrows(int) (int, error)
    DivideBy(int) (float64, error)
}

func NewDemoLib() DemoLib {
    return _swig_wrap_new_Wrapped_DemoLib()
}

```

Also in the `go_wrapper` are wrapper functions for `NegativeThrows` and `DivideBy` which recover from panics caused by exceptions and write the panic contents to the returned error. Unfortunately, this winds up with a lot of repeated code across each function, but the result is a clean go-like API with exceptions safely caught and transformed into errors.
```
func (e SwigcptrWrapped_DemoLib) NegativeThrows(n int) (i int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	i = e.Wrapped_NegativeThrows(n)
	return
}

func (e SwigcptrWrapped_DemoLib) DivideBy(n int) (f float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	f = e.Wrapped_DivideBy(n)
	return
}
```

### Go Tests
Golang tests which demonstrate the demolib API in action. When the functions are called with values that throw exceptions in the C++ library, those exceptions are translated into go errors and can be handled as such.
```Go
var (
	demo = demolib.NewDemoLib()
)

func TestThrowsNegativeThrows(t *testing.T) {
	expectedErr := "NegativeThrows threw exception"
	_, err := demo.NegativeThrows(-1)

	if err == nil {
		t.Fatal("Expected an error.")
	}
	if err.Error() != expectedErr {
		t.Errorf("Expected error message %v but got %v", expectedErr, err.Error())
	}
}

func TestDivideByZero(t *testing.T) {
	expectedErr := "Cannot divide by zero"
	_, err := demo.DivideBy(0)

	if err == nil {
		t.Fatal("Expected an error when dividing by zero.")
	}
	if err.Error() != expectedErr {
		t.Errorf("Expected error message %v but got %v", expectedErr, err.Error())
	}
}

func TestNeverThrowsReturnsInput(t *testing.T) {
	n := demo.NeverThrows(-1)

	if n != -1 {
		t.Errorf("Expected -1 but got %v", n)
	}
}
```
