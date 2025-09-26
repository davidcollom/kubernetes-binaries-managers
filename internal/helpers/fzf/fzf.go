package fzf

import (
    "errors"
    "fmt"
    "os"

    fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
    "golang.org/x/term"
)

// ErrNonInteractive indicates stdout isn't a TTY; items were printed instead.
var ErrNonInteractive = errors.New("non-interactive mode: items printed to stdout")

// Select displays a fuzzy finder over the provided items and returns the chosen item.
// When stdout is not a TTY (e.g., output is piped), it prints items to stdout and returns ErrNonInteractive.
func Select(items []string, prompt string) (string, error) {
    if !term.IsTerminal(int(os.Stdout.Fd())) {
        for _, it := range items {
            fmt.Fprintln(os.Stdout, it)
        }
        return "", ErrNonInteractive
    }

    idx, err := fuzzyfinder.Find(
        items,
        func(i int) string { return items[i] },
        fuzzyfinder.WithPromptString(prompt),
    )
    if err != nil {
        return "", err
    }
    return items[idx], nil
}
