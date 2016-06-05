# instructions
--
    import "github.com/MJKWoolnough/instructions"


## Usage

```go
var (
	ErrInvalidEscape = errors.New("invalid escape character")
	ErrInvalidNumber = errors.New("invalid number format")
)
```
Errors

#### func  New

```go
func New(functionObj interface{}, data io.Reader) ([]Function, error)
```
New creates a new instruction parser from the given value - exported methods on
which will be turned into instructions

#### type Function

```go
type Function interface {
	Call() error
	Name() string
}
```

Function is used to interface the instructions to the methods
