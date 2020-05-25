# astab
Very simply formats a go slice of structs as a table, in ~100 lines of gofmt'ed code. No dependencies.

## Why?
I write a lot of little utility functions in Go, and over time I'm realizing my eyes are tired of looking at output from Printf's `%+v` format.

If you have a slice of any arbitrary struct, you can dump its exported fields to os.Stdout in a nice table now.

Here's hoping this never comes up in a job interview.
