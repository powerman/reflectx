package reflectx_test

import (
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/reflectx"
)

type valRecv struct{}

func (valRecv) typeName() string { return reflectx.CallerTypeMethodName(0) }
func (valRecv) method() string   { return reflectx.CallerMethodName(0) }
func (valRecv) pkg() string      { return reflectx.CallerPkg(0) }

type ptrRecv struct{}

func (*ptrRecv) typeName() string { return reflectx.CallerTypeMethodName(0) }
func (*ptrRecv) method() string   { return reflectx.CallerMethodName(0) }
func (*ptrRecv) pkg() string      { return reflectx.CallerPkg(0) }

type pkgPathTest struct{}

func (pkgPathTest) callerPkgPath() string { return reflectx.CallerPkgPath(0) }

// Calls to reflectx from a top-level function have no receiver type,
// so TypeMethodName and MethodName must panic.
func topTypeName() string { return reflectx.CallerTypeMethodName(0) }
func topMethod() string   { return reflectx.CallerMethodName(0) }
func topFuncName() string { return reflectx.CallerFuncName(0) }

func TestMethodsOf(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	type testInterface interface {
		A()
		B()
	}
	methods := reflectx.MethodsOf(new(testInterface))
	t.Len(methods, 2)
	t.Equal(methods[0], "A")
	t.Equal(methods[1], "B")
}

func TestMethodsOfPanics(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	t.Panic(func() { reflectx.MethodsOf(nil) })
	t.Panic(func() { reflectx.MethodsOf("not a pointer to interface") })
	t.Panic(func() { reflectx.MethodsOf(struct{}{}) })
}

func TestCallerPkg(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	t.Equal(valRecv{}.pkg(), "reflectx_test")
	t.Equal((&ptrRecv{}).pkg(), "reflectx_test")
}

func TestCallerPkgPath(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	t.HasSuffix(pkgPathTest{}.callerPkgPath(), "/reflectx_test")
}

func TestCallerMethodName(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	t.Equal(valRecv{}.method(), "method")
	t.Equal((&ptrRecv{}).method(), "method")

	// From a closure inside a method — returns the enclosing method name.
	closure := func() string { return valRecv{}.method() }
	t.Equal(closure(), "method")
}

func TestCallerMethodNamePanics(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	t.Panic(func() { topMethod() })
}

func TestCallerFuncName(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	t.Equal(reflectx.CallerFuncName(0), "TestCallerFuncName")

	// From a top-level function — should NOT panic.
	t.Equal(topFuncName(), "topFuncName")
}

func TestCallerTypeMethodName(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	t.Equal(valRecv{}.typeName(), "valRecv.typeName")
	t.Equal((&ptrRecv{}).typeName(), "ptrRecv.typeName")

	// A closure inside a function reports its enclosing context.
	closure := func() string { return valRecv{}.typeName() }
	t.Equal(closure(), "valRecv.typeName")
}

func TestCallerTypeMethodNamePanics(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	t.Panic(func() { topTypeName() })
}

func TestCallerPkg_skip(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	// skip=0 from test helper should return our test package.
	t.Equal(reflectx.CallerPkg(0), "reflectx_test")

	// From a helper function one frame up, skip=1 should also
	// still return the test package (the helper is in this package).
	checkPkg := func(skip int) string { return reflectx.CallerPkg(skip + 1) }
	t.Equal(checkPkg(0), "reflectx_test")
}
