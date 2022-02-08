![Go](https://github.com/dbsystel/golang-runtime-di/workflows/Go/badge.svg) [![codecov](https://codecov.io/gh/dbsystel/golang-runtime-di/branch/main/graph/badge.svg?token=E123SJUGFD)](https://codecov.io/gh/dbsystel/golang-runtime-di) [![Go Reference](https://pkg.go.dev/badge/github.com/dbsystel/golang-runtime-di/.svg)](https://pkg.go.dev/github.com/dbsystel/golang-runtime-di/)
# golang-runtime-di

## description

golang-runtime-di is a framework for runtime dependency injection in go.

## usage

### quickstart

- add it to your `go.mod`: `go get github.com/dbsystel/golang-runtime-di`
- create your components (interfaces are encouraged):
  ```golang
  type Producer interface {}
  type Producer1 interface { A() }
  type Producer1Impl struct {}
  func (p *Producer1Impl) A() {}

  type Producer2 struct {}
  
  type Producer3 string
  
  type Consumer struct {
    Producer1 `inject:""`               // inject by interface 
    *Producer2 `inject:""`              // inject by ptr
    Producer3 `inject:""`               // inject by value
    AllProducers []Producer `inject:""` // inject all known coercible components
  }
  ```
- create a new scope, register your components:
  ```golang
  scope := &di.Scope{}
  scope.MustRegister(&ProducerImpl1{})                                          // Component instance
  scope.MustRegister(func () *Producer2 { return &Producer2{} })                // Component factory
  scope.MustRegister(func () (Producer3, error) { return Producer3("a"), nil }) // Component factory with error
  ``` 
- Wire the target components:
  ```golang
  consumer := &Consumer{} 
  // kindly note: 
  // - you must use ptr here
  // - wired instances are not automatically registered
  scope.MustWire(consumer)
  ```

### the `inject` tag

The `inject` tag marks a struct field to be injected by the DI.

There are following options on the tag:

- `optional`: marks the dependency to be optional. this will not create an error if the dependency is not found in
  scope.
- `qualifier=<xy>`: use a qualifier to resolve the dependency. by default the qualifier is empty, thus only unqualified
  instances are selected.
- `qualifier=*`: resolve the dependency from any qualifier

Options may be combined, e.g.: `optional,qualifier=squash`.

### dependency injection

#### scoping

Each scope is isolated, but be aware that if you wire a struct multiple times in different scopes, the dependencies
maybe replace by each wiring, depending on the scopes registrations.

#### qualifiers

Qualifiers can be used to use the same dependency type more than once, the default qualifier is empty (`""`).\
Example:

```golang
scope.MustRegister(&Component{}).WithQualifier("yay")
```

#### priorities

Priorities can be used to allow overriding of components, lower numbers denote higher priorities.\
Example:

```golang
// The following factory will not even be called, if the injection of Dependency is requested
scope.MustRegister(func () (*Dependency, error) { return nil, errors.New("meh") }).WithPriority(1)
// Note: priority is higher
scope.MustRegister(&Dependency{}).WithPriority(-1)
```

#### registrations and wiring

For the registrations in a scope the following rules apply:

- matching the dependency fields and registered components will be done using reflection (`Type.AssignableTo()`
  and `Type.Implements()`, resp.)
- if you are using factory functions, factories for registered components will only be called if necessary
- wiring of the components will only happen once when the component is to be injected the first time.

#### component resolution

The component resolution for injection sticks by the following rules (imperatively applied):

- select registrations which can be coerced to the fields type
- select the registrations matching the qualifier (skipped if `qualifier=*`)
- order by priority
- *if single dependency injection* (not slice):
    - select highest priority (= lowest number)
    - error if there's more than a single candidate with that priority
- *if dependency is required (not `optional`)*:
    - error if no candidates found

## examples

This project comes with tested examples:

- [Simple wiring](./examples/simple.go)
- [Optional wiring](./examples/optional.go)
- [Wiring using factory method](./examples/factory.go)
- [Qualifier example](./examples/qualifier.go)
- [Priority example](./examples/priority.go)
- [Parent scope example](./examples/parent.go)

## License

This project is licensed under Apache License v2.0, which is [included in the repository](./LICENSE.txt).

## Contributions

Contributions are very welcome, please refer to the [contribution guide](./CONTRIBUTING.md).

## Code of conduct
Our code of conduct can be found [here](./CODE_OF_CONDUCT.md).
