# Rules for Large Language Models (LLMs)

## Error handling

Use go-faster/errors:

```
import "github.com/go-faster/errors"

func f() error {
    return errors.New("something went wrong")
}

func g() error {
    if err := f(); err != nil {
        return errors.Wrap(err, "g") // NB: Do not add "failed:" prefix.
    }
    return nil
}
```
