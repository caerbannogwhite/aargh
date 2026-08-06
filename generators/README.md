## Enchanter Gen

Code generator for the `series` package: it produces the `*_base.go` files
from templates and rewrites the type-specific operation code in the
`*_ops.go` files.

It is a standalone Go module. Run it from the repository root with:

```sh
go generate ./series/
```

or directly:

```sh
go -C generators run .
```

Generation is idempotent; CI fails if the checked-in generated files drift
from the generator output.
