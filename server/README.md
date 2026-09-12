# EcomSite Server — Architecture & Contribution Guide

This document explains how the codebase is put together, the recipe for adding a
new service/feature, and walks through the Cart/Order modeling question in
detail since that's the part that's easiest to get wrong.

## 1. Layered architecture

```
cmd/api/            entrypoint + wiring (main.go, api.go, middleware.go)
config/             env parsing, DB config struct
internal/db/        raw driver setup (postgres pool, redis client)
internal/models/    plain domain structs — no behavior, no DB tags beyond json
internal/handlers/  HTTP layer: gin handlers, grouped by feature (auth/, ...)
store/              persistence layer: one Repository interface + impl per aggregate
lib/                stateless helpers: hashing, jwt, validation
```

The dependency direction is one-way:

```
handlers  →  store (interfaces)  →  internal/db (driver)
   ↓
models (everyone depends on models; models depends on nothing)
```

Handlers never import `database/sql` or `redis` directly — they only ever see
a `store.XRepository` interface. That's what makes handlers testable with a
fake store and lets `store/` swap Postgres for something else later.

## 2. Request lifecycle

1. `cmd/api/main.go` calls `Setup()` (in `api.go`).
2. `Setup()` loads env vars (`config.AllEnvs`), opens the Postgres pool and
   Redis client (`internal/db`), then builds a `store.Storage` — a struct of
   repository interfaces (`store/store.go`).
3. `(a *application) routes()` builds the `gin.Engine`, attaches global
   middleware (`gin.Logger`, `gin.Recovery`, `CORS()`), and registers route
   groups. Handlers are constructed by passing in the repository they need,
   e.g. `auth.Login(app.store.Customer)`.
4. `run()` starts `http.Server` with the built engine.

So a single request flows: **router → handler closure → repository interface
→ concrete `*Store` → SQL/Redis**.

## 3. How to add a new service (step-by-step recipe)

Say you want to add, e.g., a `Wishlist` feature. Follow the same shape as
`Customer` (`store/client.go` + `internal/handlers/auth/`):

**Step 1 — Model** (`internal/models/models.go`)
Add the plain struct. No DB-specific tags, no methods that talk to the DB.

```go
type Wishlist struct {
    ID       string    `json:"id"`
    ClientID string    `json:"client_id"`
    Products []Product `json:"products"`
}
```

**Step 2 — Repository interface + store impl** (`store/wishlist.go`)
One file per aggregate. The interface documents the contract; the struct
implements it and holds the driver handle.

```go
package store

type WishlistRepository interface {
    GetByClientID(ctx context.Context, clientID string) (*models.Wishlist, error)
    AddProduct(ctx context.Context, clientID string, productID int64) error
}

type WishlistStore struct {
    db *sql.DB
}

func (s *WishlistStore) GetByClientID(ctx context.Context, clientID string) (*models.Wishlist, error) {
    // SQL here
}
```

**Step 3 — Register it in `store/store.go`**

```go
type Storage struct {
    Product   ProductRepository
    Order     OrderRepository
    Customer  CustomerRepository
    Cache     CacheRepository
    Wishlist  WishlistRepository // new
}

func NewStorage(db *sql.DB, rclient *redis.Client) Storage {
    return Storage{
        Product:  &ProductStore{db},
        Order:    &OrderStore{db},
        Customer: &CustomerStore{db},
        Cache:    &CacheStore{rclient},
        Wishlist: &WishlistStore{db}, // new
    }
}
```

**Step 4 — Handler package** (`internal/handlers/wishlist/`)
Mirror `internal/handlers/auth/login.go`: the handler is a function that
takes the *interface* (not the struct) and returns a `gin.HandlerFunc`. This
is what keeps handlers unit-testable.

```go
package wishlist

func Get(r store.WishlistRepository) func(c *gin.Context) {
    return func(c *gin.Context) {
        list, err := r.GetByClientID(c.Request.Context(), c.Param("clientId"))
        if err != nil {
            c.AbortWithError(http.StatusNotFound, err)
            return
        }
        c.JSON(http.StatusOK, list)
    }
}
```

**Step 5 — Wire the route** (`cmd/api/api.go`, inside `routes()`)

```go
protected := r.Group("/api")
protected.GET("/wishlist/:clientId", wishlist.Get(a.store.Wishlist))
```

That's the whole loop: **model → repository interface/impl → registered in
Storage → handler takes the interface → route wired in `api.go`**. Every
existing feature (`Customer`, `Product`, `Order`, `Cache`) follows this same
shape, even where the implementation is still a stub returning `nil, nil`
(see `store/order.go`, `store/product.go` — those are scaffolded but not
implemented yet).

## 4. Deep dive: `Cart`/`CartItem` vs `Order`/`LineOrder`

This is the actual conceptual split you're stuck on, so here's the reasoning
plus what's currently inconsistent in `internal/models/models.go`.

### The core distinction

| | Cart | Order |
|---|---|---|
| Represents | **intent** — what a customer is currently considering buying | **fact** — a completed transaction, a legal/financial record |
| Mutability | mutable — items added/removed/quantity changed constantly | immutable once placed — never edited after checkout |
| Price source | **live** — always reflects current `Product.Price` | **frozen** — must snapshot the price at the moment of purchase |
| Lifetime | short-lived, often per-session, can be abandoned | permanent, kept for accounting/returns/support |
| Relationship to `Product` | reference by ID + quantity | reference by ID + quantity **+ a frozen copy of price/name at purchase time** |

Because of the mutability/pricing difference, a cart line and an order line
are *not the same shape*, even though they both conceptually say "N of
product X":

- **`CartItem`** = a *pointer* into the product catalog + a quantity. It has
  no price of its own — it always asks `Product` for the current price,
  because if the product's price changes tomorrow, the cart should reflect
  that (nothing has been paid for yet).
- **`LineOrder`** = a *frozen receipt line*. It must carry its own copy of
  unit price (and ideally product name) at the time of purchase, because if
  the product's price changes next month, past orders must still show what
  the customer actually paid. This is why order lines are traditionally
  called "snapshots."

### What's inconsistent right now

Look at the current models:

```go
type LineOrder struct {
    ID        int64
    Quantity  int
    UnitPrice float64
    Total     int
}

type Cart struct {
    ID       string
    ClientID string
    Products []Product // full Product embedded, no quantity per item
}
```

Two problems fall directly out of the table above:

1. **`LineOrder` has no reference to which `Product` was bought.** It has a
   quantity and a price snapshot, but no `ProductID` (or embedded product
   name). You can total an order, but you can't render "2x Wireless Mouse"
   on an invoice, and you can't look up what was actually purchased.
2. **`Cart` embeds `[]Product` directly instead of `[]CartItem`.** This
   means a cart has no concept of quantity — adding "2 of the same product"
   either requires duplicating the full `Product` struct twice in the slice
   (wasteful, and `Total` becomes ambiguous) or the cart can only ever hold
   quantity 1 of each item. It also freezes catalog data (price, description,
   stock) into the cart, which contradicts the cart's "always live" nature.

### Recommended shapes

```go
// Cart: a reference + quantity, resolved against the live catalog at read time
type CartItem struct {
    ProductID int64 `json:"product_id"`
    Quantity  int   `json:"quantity"`
}

type Cart struct {
    ID       string     `json:"id"`
    ClientID string     `json:"client_id"`
    Items    []CartItem `json:"items"`
}
// Total is *computed* on read (Item.Quantity * live Product.Price), not stored.

// Order: a frozen snapshot, never recomputed from the catalog again
type LineOrder struct {
    ID        int64   `json:"id"`
    ProductID int64   `json:"product_id"`     // <- what was actually bought
    ProductName string `json:"product_name"`  // optional but nice for invoices
    Quantity  int     `json:"quantity"`
    UnitPrice float64 `json:"price"`           // price AT PURCHASE TIME, frozen
    Total     float64 `json:"total"`
}

type Order struct {
    ID              int64       `json:"id"`
    ClientID        string      `json:"client_id"`
    Status          string      `json:"status"`
    PaymentIntentID string      `json:"payment_intent_id"`
    Total           float64     `json:"total"`
    ShippingAddress string      `json:"shipping_address"`
    ProductsPlaced  []LineOrder `json:"products_placed"`
    PlacedAt        time.Time   `json:"placed_at"`
    PaidAt          *time.Time  `json:"paid_at,omitempty"`
}
```

### The checkout boundary

The reason these two types must diverge is the **checkout step itself**:
converting a `Cart` into an `Order` is exactly the operation that takes live
`CartItem` + current `Product.Price` and freezes them into `LineOrder`s.

```
Cart (live, mutable)  --[checkout]-->  Order (frozen, immutable)
CartItem{ProductID, Qty}  --resolve Product.Price-->  LineOrder{ProductID, Qty, UnitPrice(frozen), Total}
```

If `Cart` and `Order` shared the exact same item shape, you'd lose the
ability to express "this is still editable and price-live" vs "this is
locked and historical" — which is the whole point of having two types
instead of one. A useful rule of thumb when you're unsure which one to
reach for:

- Need something a user can still add/remove/change quantity on, before any
  money moves? → `Cart` / `CartItem`.
- Need something that must stay accurate forever regardless of future
  catalog changes, tied to a payment? → `Order` / `LineOrder`.

## 5. Current gaps worth knowing about

- `store/order.go`, `store/product.go`, `store/cache.go` are scaffolded
  (interface defined, method returns `nil, nil`) — no SQL implemented yet.
- `Order.Total` and `LineOrder.Total`/`UnitPrice` mix `int` and `float64` in
  the current struct — worth picking one (float64, or integer cents to avoid
  floating point rounding on money) before wiring up checkout.
- No `Cart` repository exists yet in `store/` — needed once `CartItem` lands.
- Routes in `cmd/api/api.go` (`auth`, `protected` groups) are still commented
  out even though the handlers in `internal/handlers/auth/` exist.
