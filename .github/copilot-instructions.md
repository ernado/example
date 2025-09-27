# Rules for Large Language Models (LLMs)

[github](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions/add-repository-instructions?tool=jetbrains#about-repository-custom-instructions-for-copilot-3)

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
