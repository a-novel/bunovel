# Bunovel

A-Novel utils library for the bun ORM, by Uptrace.

```bash
go get github.com/a-novel/bunovel
```

## Initialization

Bunovel implements an object-based initialization (in opposition to bun's imperative-based initialization), which 
automates some of its aspects.

The simplest initialization of bun with bunovel is close to the original package. However, it automates the creation
of the intermediate sql.DB object.

```go
import (
    "github.com/a-novel/bunovel"
    "os"
)

func main() {
    dsn := os.Getenv("POSTGRES_URL")
    bunDB, sqlDB, err := bunovel.NewClientWithDriver(bunovel.NewPGDriverWithDSN(dsn))
    
    defer bunDB.Close()
    defer sqlDB.Close()
}
```

For more advanced configuration, bunovel implements an object to automate some startup tasks, such as running migrations.

```go
import (
    "github.com/a-novel/bunovel"
    "github.com/my-package/migrations"
    "os"
)

func main() {
    ctx := context.Background()
    dsn := os.Getenv("POSTGRES_URL")
    bunDB, sqlDB, err := bunovel.NewClient(ctx, bunovel.Config{
        Driver: bunovel.NewPGDriverWithDSN(dsn),
        Migrations: bunovel.MigrationsConfig{
            Files: []fs.FS{migrations},
        },
        DiscardUnknownColumns: true,
    })
}
```

Every option is documented directly in Go. Driver can also be passed as an object configuration, rather than providing
each option as a function argument.

```go
import (
    "github.com/a-novel/bunovel"
    "github.com/my-package/migrations"
    "os"
)

func main() {
    ctx := context.Background()
    dsn := os.Getenv("POSTGRES_URL")
    bunDB, sqlDB, err := bunovel.NewClient(ctx, bunovel.Config{
        Driver: bunovel.PGDriver{
            DSN: dsn,
        },
        Migrations: bunovel.MigrationsConfig{
            Files: []fs.FS{migrations},
        },
        DiscardUnknownColumns: true,
    })
}
```

## Errors

Bunovel extends PostgreSQL error handling capabilities.

### HandlePGError

```go
import (
    "github.com/a-novel/bunovel"
    "github.com/uptrace/bun"
)

func main() {
    // ...
    
    res, err := bun.NewInsert().Model(myModel).Exec(ctx)
    err = bunovel.HandlePGError(err)
}
```

HandlePGError will parse the SQL error message, to return more precise errors, directly checkable in Go.
The following, additional errors can be returned:

- `bunovel.ErrNotFound`: alias for `sql.ErrNoRows`.
- `bunovel.ErrConstraintViolation`: the insertion/update/deletion failed, because it violates one or more constraints.
- `bunovel.ErrUniqConstraintViolation`: the insertion/update failed, because it violates a unique constraint (duplicate record).

### ForceRowsUpdate

As PostgreSQL does not throw an error when an update has no effect, this function checks the result of a mutation, and
returns an error if no row was updated.

```go
import (
    "github.com/a-novel/bunovel"
    "github.com/uptrace/bun"
)

func main() {
    // ...

    res, err := bun.NewInsert().Model(myModel).Exec(ctx)
    if err != nil {
        return bunovel.HandlePGError(err)        
    }
    if err := bunovel.ForceRowsUpdate(res); err != nil {
        return err
    }
}
```

## Models

Bunovel adds some convenient model, used across the a-novel apps.

### Metadata

Metadata implements a common metadata object for models, with the basic id/created_at/updated_at fields.

```go
import (
    "github.com/a-novel/bunovel"
    "github.com/uptrace/bun"
)

type MyModel struct {
    bun.BaseModel `bun:"table:my_model"`
    bunovel.Metadata
}
```

The model above matches the following PostgreSQL migration.

```sql
CREATE TABLE my_model (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ
);
```

## Types

Bunovel automatically extends bun with extra types. For now, the `time.Duration` type is natively supported.
