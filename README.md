# enumcase

> `enumcase` checks every `switch` statement handles all const values of the type

## Install

```console
$ go get -u github.com/MakeNowJust/enumcase/cmd/enumcase
```

## Usage

```console
$ go vet -vettool=$(which enumcase) pkgname
```

Or

```console
$ enumcase pkgname
```

## Example

For example you have this type and consts:

```go
type FileMode int

const (
    Read FileMode = iota
    Write
    Append
)
```

Then, `enumcase` reports a warning for such a switch because the case to `Append` is missing.

```go
switch mode {
case Read:
    // ...
case Write:
    // ...
}
```

```console
$ go vet -vettool=$(which enumcase) .
/.../main.go:10:9: missing case(s) to FileMode value(s): Append
```

## Notice

`enumcase` reports many false-positives because this tool checks all switch statement whose tag type has `const` value.
It is hard to distinguish whether the type is enum-like or not.

I **don't recommend** to use `enumcase` on CI or everyday code check.
However, I **recommend** to use `enumcase` when a new `const` value is added.
It may help you.
