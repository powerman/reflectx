// Package reflectx provides generic helpers built on reflect and runtime,
// such as listing an interface's methods and resolving the caller's
// function/method name or package name from the call stack.
//
// The Caller* functions — CallerPkg, CallerPkgPath, CallerMethodName,
// CallerFuncName, CallerTypeMethodName — examine the program counter
// of the calling function by using [runtime.Caller]. Each accepts a
// skip parameter that controls how many stack frames to ascend:
//   - 0 = the immediate caller of the Caller* function
//   - 1 = the caller's caller, and so on.
//
// Panic behavior:
//   - CallerMethodName and CallerTypeMethodName panic when the caller
//     is a top-level function (no receiver type).
//     They are designed for use inside methods.
//   - CallerFuncName handles both functions and methods without panicking.
//   - CallerPkg and CallerPkgPath panic if the call stack entry is
//     malformed — this should never happen in practice.
//   - MethodsOf panics when v is not a pointer to an interface.
package reflectx

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// callDepthBase is the number of internal stack frames between
// [runtime.Caller] and an exported function's immediate caller that
// must always be skipped: callerName (1) + the exported function (1).
const callDepthBase = 2

// MethodsOf returns the names of all methods of the interface v points to.
// v must be a pointer to an interface, e.g. new(Service).
func MethodsOf(v any) []string {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Interface {
		panic("require pointer to interface")
	}
	typ = typ.Elem()
	names := make([]string, typ.NumMethod())
	for i := range typ.NumMethod() {
		names[i] = typ.Method(i).Name
	}
	return names
}

// CallerPkg returns the package name of the caller at the given
// stack depth. For a function in "a/b/pkg.Func" it returns "pkg".
func CallerPkg(skip int) string {
	return pkgName(callerName(callDepthBase + skip))
}

// CallerPkgPath returns the full package path of the caller at the
// given stack depth. For a function in "a/b/pkg.Func" it returns
// "a/b/pkg".
func CallerPkgPath(skip int) string {
	return pkgPath(callerName(callDepthBase + skip))
}

// CallerFuncName returns the function or method name of the caller at the given stack depth.
// Unlike CallerMethodName, it works for both top-level functions and methods without panicking.
func CallerFuncName(skip int) string {
	return funcName(callerName(callDepthBase + skip))
}

// CallerMethodName returns the method name of the caller at the given stack depth.
// For a method "pkg.Type.Method" it returns "Method".
// For a closure "pkg.Type.Method.func1" it returns "Method".
// It panics if the caller is a top-level function.
func CallerMethodName(skip int) string {
	return methodName(callerName(callDepthBase + skip))
}

// CallerTypeMethodName returns the "Type.Method" name of the caller at the given stack depth,
// with any pointer-receiver parentheses stripped (e.g. "(*T).M" becomes "T.M").
// For a closure it returns "EnclosingType.EnclosingMethod".
// It panics if the caller is a top-level function.
func CallerTypeMethodName(skip int) string {
	return stripTypeRef(typeMethodName(callerName(callDepthBase + skip)))
}

// callerName returns the fully-qualified runtime function name of the
// function skip frames up the call stack (0 = callerName itself).
func callerName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	return runtime.FuncForPC(pc).Name()
}

// pkgPath returns the full package path from a runtime function name.
func pkgPath(name string) string {
	start := strings.LastIndexByte(name, '/') + 1
	end := strings.IndexByte(name[start:], '.')
	if end == -1 {
		panic(fmt.Sprintf("bad runtime func name: %s", name))
	}
	return name[:start+end]
}

// pkgName returns the package name (last path segment) from a runtime function name.
func pkgName(name string) string {
	start := strings.LastIndexByte(name, '/') + 1
	end := strings.IndexByte(name[start:], '.')
	if end == -1 {
		panic(fmt.Sprintf("bad runtime func name: %s", name))
	}
	return name[start : start+end]
}

// typeMethodName returns the "type.method" (or "func.funcN") portion
// of a runtime function name, dropping the package path.
// It panics if the name refers to a top-level function.
func typeMethodName(name string) string {
	start := strings.LastIndexByte(name, '/') + 1
	pos := strings.IndexByte(name[start:], '.')
	if pos == -1 {
		panic(fmt.Sprintf("bad runtime func name: %s", name))
	}
	start += pos + 1
	pos = strings.IndexByte(name[start:], '.')
	if pos == -1 {
		panic(fmt.Sprintf("not a method name: %s", name))
	}
	end := strings.IndexByte(name[start+pos+1:], '.')
	if end == -1 {
		end = len(name)
	} else {
		end += start + pos + 1
	}
	return name[start:end]
}

// methodName extracts the method name from the "type.method" portion of
// a runtime function name.
// It panics if the name refers to a top-level function;
// for closures inside a top-level function it returns "funcN".
func methodName(name string) string {
	name = typeMethodName(name)
	pos := strings.IndexByte(name, '.')
	if pos == -1 {
		panic(fmt.Sprintf("not a method name: %s", name))
	}
	return name[pos+1:]
}

// funcName returns the function or method name from a runtime function name.
// Unlike methodName, it handles top-level functions without panicking
// (returns the function name directly).
func funcName(name string) string {
	start := strings.LastIndexByte(name, '/') + 1
	pos := strings.IndexByte(name[start:], '.')
	if pos == -1 {
		panic(fmt.Sprintf("bad runtime func name: %s", name))
	}
	start += pos + 1
	pos = strings.IndexByte(name[start:], '.')
	if pos == -1 {
		return name[start:]
	}
	return methodName(name)
}

// stripTypeRef removes the pointer-receiver parentheses from a "(*T).method" name,
// returning "T.method". Other names are returned unchanged.
func stripTypeRef(name string) string {
	if name != "" && name[0] == '(' {
		pos := strings.IndexByte(name, ')')
		name = name[2:pos] + name[pos+1:]
	}
	return name
}
