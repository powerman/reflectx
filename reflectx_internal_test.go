package reflectx

import (
	"testing"

	"github.com/powerman/check"
)

func TestPkgPath(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	tests := []struct {
		given     string
		want      string
		wantPanic bool
	}{
		{"", "", true},
		{"bad", "", true},
		{"main.main", "main", false},
		{"main.f", "main", false},
		{"main.f.func1", "main", false},
		{"main.f.func2", "main", false},
		{"main.f.func2.1", "main", false},
		{"main.f.func2.1.1", "main", false},
		{"main.f.func3", "main", false},
		{"main.T.m", "main", false},
		{"main.T.m.func1", "main", false},
		{"main.T.m.func2", "main", false},
		{"main.T.m.func2.1", "main", false},
		{"main.(*T).M", "main", false},
		{"github.com/powerman/whoami/subpkg.F", "github.com/powerman/whoami/subpkg", false},
		{"github.com/powerman/whoami/subpkg.F.func1", "github.com/powerman/whoami/subpkg", false},
		{"github.com/powerman/whoami/subpkg.F.func2", "github.com/powerman/whoami/subpkg", false},
		{"github.com/powerman/whoami/subpkg.F.func2.1", "github.com/powerman/whoami/subpkg", false},
		{"github.com/powerman/whoami/subpkg.(*T).M", "github.com/powerman/whoami/subpkg", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func1", "github.com/powerman/whoami/subpkg", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2", "github.com/powerman/whoami/subpkg", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2.1", "github.com/powerman/whoami/subpkg", false},
	}
	for _, tc := range tests {
		t.Run(tc.given, func(tt *testing.T) {
			tt.Parallel()
			t := check.T(tt).MustAll()
			if tc.wantPanic {
				t.Panic(func() { pkgPath(tc.given) })
			} else {
				t.Equal(pkgPath(tc.given), tc.want)
			}
		})
	}
}

func TestPkgName(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	tests := []struct {
		given     string
		want      string
		wantPanic bool
	}{
		{"", "", true},
		{"bad", "", true},
		{"main.main", "main", false},
		{"main.f", "main", false},
		{"main.f.func1", "main", false},
		{"main.f.func2", "main", false},
		{"main.f.func2.1", "main", false},
		{"main.f.func2.1.1", "main", false},
		{"main.f.func3", "main", false},
		{"main.T.m", "main", false},
		{"main.T.m.func1", "main", false},
		{"main.T.m.func2", "main", false},
		{"main.T.m.func2.1", "main", false},
		{"main.(*T).M", "main", false},
		{"github.com/powerman/whoami/subpkg.F", "subpkg", false},
		{"github.com/powerman/whoami/subpkg.F.func1", "subpkg", false},
		{"github.com/powerman/whoami/subpkg.F.func2", "subpkg", false},
		{"github.com/powerman/whoami/subpkg.F.func2.1", "subpkg", false},
		{"github.com/powerman/whoami/subpkg.(*T).M", "subpkg", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func1", "subpkg", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2", "subpkg", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2.1", "subpkg", false},
	}
	for _, tc := range tests {
		t.Run(tc.given, func(tt *testing.T) {
			tt.Parallel()
			t := check.T(tt).MustAll()
			if tc.wantPanic {
				t.Panic(func() { pkgName(tc.given) })
			} else {
				t.Equal(pkgName(tc.given), tc.want)
			}
		})
	}
}

func TestTypeMethodName(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	tests := []struct {
		given     string
		want      string
		wantPanic bool
	}{
		{"", "", true},
		{"bad", "", true},
		{"main.main", "", true},
		{"main.f", "", true},
		{"main.f.func1", "f.func1", false},
		{"main.f.func2", "f.func2", false},
		{"main.f.func2.1", "f.func2", false},
		{"main.f.func2.1.1", "f.func2", false},
		{"main.f.func3", "f.func3", false},
		{"main.T.m", "T.m", false},
		{"main.T.m.func1", "T.m", false},
		{"main.T.m.func2", "T.m", false},
		{"main.T.m.func2.1", "T.m", false},
		{"main.(*T).M", "(*T).M", false},
		{"github.com/powerman/whoami/subpkg.F", "", true},
		{"github.com/powerman/whoami/subpkg.F.func1", "F.func1", false},
		{"github.com/powerman/whoami/subpkg.F.func2", "F.func2", false},
		{"github.com/powerman/whoami/subpkg.F.func2.1", "F.func2", false},
		{"github.com/powerman/whoami/subpkg.(*T).M", "(*T).M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func1", "(*T).M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2", "(*T).M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2.1", "(*T).M", false},
	}
	for _, tc := range tests {
		t.Run(tc.given, func(tt *testing.T) {
			tt.Parallel()
			t := check.T(tt).MustAll()
			if tc.wantPanic {
				t.Panic(func() { typeMethodName(tc.given) })
			} else {
				t.Equal(typeMethodName(tc.given), tc.want)
			}
		})
	}
}

func TestMethodName(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	tests := []struct {
		given     string
		want      string
		wantPanic bool
	}{
		{"", "", true},
		{"bad", "", true},
		{"main.main", "", true},
		{"main.f", "", true},
		{"main.f.func1", "func1", false},
		{"main.f.func2", "func2", false},
		{"main.f.func2.1", "func2", false},
		{"main.f.func2.1.1", "func2", false},
		{"main.f.func3", "func3", false},
		{"main.T.m", "m", false},
		{"main.T.m.func1", "m", false},
		{"main.T.m.func2", "m", false},
		{"main.T.m.func2.1", "m", false},
		{"main.(*T).M", "M", false},
		{"github.com/powerman/whoami/subpkg.F", "", true},
		{"github.com/powerman/whoami/subpkg.F.func1", "func1", false},
		{"github.com/powerman/whoami/subpkg.F.func2", "func2", false},
		{"github.com/powerman/whoami/subpkg.F.func2.1", "func2", false},
		{"github.com/powerman/whoami/subpkg.(*T).M", "M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func1", "M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2", "M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2.1", "M", false},
	}
	for _, tc := range tests {
		t.Run(tc.given, func(tt *testing.T) {
			tt.Parallel()
			t := check.T(tt).MustAll()
			if tc.wantPanic {
				t.Panic(func() { methodName(tc.given) })
			} else {
				t.Equal(methodName(tc.given), tc.want)
			}
		})
	}
}

func TestFuncName(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	tests := []struct {
		given     string
		want      string
		wantPanic bool
	}{
		{"", "", true},
		{"bad", "", true},
		{"main.main", "main", false},
		{"main.f", "f", false},
		{"main.f.func1", "func1", false},
		{"main.f.func2", "func2", false},
		{"main.f.func2.1", "func2", false},
		{"main.f.func2.1.1", "func2", false},
		{"main.f.func3", "func3", false},
		{"main.T.m", "m", false},
		{"main.T.m.func1", "m", false},
		{"main.T.m.func2", "m", false},
		{"main.T.m.func2.1", "m", false},
		{"main.(*T).M", "M", false},
		{"github.com/powerman/whoami/subpkg.F", "F", false},
		{"github.com/powerman/whoami/subpkg.F.func1", "func1", false},
		{"github.com/powerman/whoami/subpkg.F.func2", "func2", false},
		{"github.com/powerman/whoami/subpkg.F.func2.1", "func2", false},
		{"github.com/powerman/whoami/subpkg.(*T).M", "M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func1", "M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2", "M", false},
		{"github.com/powerman/whoami/subpkg.(*T).M.func2.1", "M", false},
	}
	for _, tc := range tests {
		t.Run(tc.given, func(tt *testing.T) {
			tt.Parallel()
			t := check.T(tt).MustAll()
			if tc.wantPanic {
				t.Panic(func() { funcName(tc.given) })
			} else {
				t.Equal(funcName(tc.given), tc.want)
			}
		})
	}
}

func TestStripTypeRef(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt).MustAll()

	tests := []struct {
		given string
		want  string
	}{
		{"T.M", "T.M"},
		{"(*T).M", "T.M"},
		{"(*PackageType).Method", "PackageType.Method"},
		{"Func.func1", "Func.func1"},
		{"T.m.func1", "T.m.func1"},
		{"", ""},
	}
	for _, tc := range tests {
		t.Run(tc.given, func(tt *testing.T) {
			tt.Parallel()
			t := check.T(tt).MustAll()
			t.Equal(stripTypeRef(tc.given), tc.want)
		})
	}
}
