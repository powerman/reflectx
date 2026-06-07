package reflectx_test

import (
	"context"
	"fmt"

	"github.com/powerman/reflectx"
)

// exampleService demonstrates how a production service would use CallerMethodName
// and CallerTypeMethodName to tag errors or logs with the method name
// without explicit configuration at each call site.
type exampleService struct{}

func (exampleService) Create(_ context.Context) string {
	return fmt.Sprintf("%s: not found", reflectx.CallerMethodName(0))
}

func (exampleService) Get(_ context.Context) string {
	return reflectx.CallerTypeMethodName(0)
}

// ExampleMethodsOf shows how to enumerate all methods of an interface.
//
// This is useful for pre-registering Prometheus metric label dimensions
// before any requests arrive, so every requested label combination is
// populated from the start.
func ExampleMethodsOf() {
	type Service interface {
		Create(ctx context.Context, _ any) error
		Get(ctx context.Context, _ int64) (any, error)
		Delete(ctx context.Context, _ int64) error
	}

	for _, name := range reflectx.MethodsOf(new(Service)) {
		fmt.Println(name)
	}
	// Unordered output:
	// Create
	// Delete
	// Get
}

// ExampleCallerMethodName shows how to capture the calling method's name
// without passing it explicitly through the call chain.
//
// This is used e.g. by transaction wrappers (NoTx/Tx) to tag every DAL
// error with the name of the repository method that triggered it.
func ExampleCallerMethodName() {
	s := exampleService{}
	out := s.Create(context.Background())
	fmt.Println(out)
	// Output: Create: not found
}

// ExampleCallerTypeMethodName shows how to capture the "Type.Method"
// string from the call stack, with any pointer-receiver parentheses stripped.
//
// This is used e.g. by middleware factories to identify exactly which
// service method is being invoked, without manual method-name wiring.
func ExampleCallerTypeMethodName() {
	s := exampleService{}
	method := s.Get(context.Background())
	fmt.Println(method)
	// Output: exampleService.Get
}

// ExampleCallerPkg shows how to detect the caller's package name from the call stack.
//
// This is used e.g. by logging middleware factories (MakeUnaryServerLogger,
// MakeStreamServerLogger) that need to know which package created them
// so each service's logs are attributed automatically.
func ExampleCallerPkg() {
	// In a real application the returned value would be something like
	// "repo", "auth", "grpcx" — the name of the calling package.
	fmt.Println(reflectx.CallerPkg(0))
	// Output: reflectx_test
}
