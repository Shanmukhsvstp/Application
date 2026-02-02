Perfect question 😄 — this is where Go + SQL really **clicks**.

I’ll show you **exactly how SQL scripts are used in Go**, in a way that maps to how you already think about it.

We’ll cover:
1️⃣ SELECT
2️⃣ INSERT
3️⃣ UPDATE
4️⃣ DELETE
5️⃣ Single row vs multiple rows
6️⃣ Where SQL *lives* (inline vs files)

All with **pgx**, since that’s what you’re using.

---

## 🧠 First: how Go talks to SQL (mental model)

In Go:

> SQL is just a **string**
> The DB driver executes it
> Go scans results into variables

There’s no ORM magic unless you add one.

---

# 1️⃣ INSERT (most common)

### Example: Register user

```go
_, err := h.DB.Exec(
	context.Background(),
	`
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	`,
	req.Username,
	req.Email,
	hashedPassword,
)
```

### What’s happening

* `$1, $2, $3` → placeholders
* Values passed **separately** (SQL injection safe)
* `Exec` → use when you **don’t expect rows back**

Java parallel:

```java
preparedStatement.executeUpdate();
```

---

## 🔑 Rule

👉 Use `Exec` for:

* INSERT
* UPDATE
* DELETE
  (no returned rows)

---

# 2️⃣ SELECT (single row)

### Example: Login

```go
var userID string
var passwordHash string

err := h.DB.QueryRow(
	context.Background(),
	`
	SELECT id, password
	FROM users
	WHERE email = $1
	`,
	req.Email,
).Scan(&userID, &passwordHash)
```

### Important behavior

* If no row → `err == pgx.ErrNoRows`
* If row exists → values filled

Always check error 👇

```go
if err != nil {
	if err == pgx.ErrNoRows {
		// user not found
	}
	// real DB error
}
```

---

## 🔑 Rule

👉 Use `QueryRow` when you expect **exactly one row**

---

# 3️⃣ SELECT (multiple rows)

### Example: Get all users

```go
rows, err := h.DB.Query(
	context.Background(),
	`SELECT id, username, email FROM users`,
)
if err != nil {
	return err
}
defer rows.Close()

for rows.Next() {
	var id, username, email string
	err := rows.Scan(&id, &username, &email)
	if err != nil {
		return err
	}

	fmt.Println(id, username, email)
}
```

### Why `defer rows.Close()`?

* Frees DB resources
* REQUIRED

Java parallel:

```java
ResultSet rs = stmt.executeQuery();
while (rs.next()) { ... }
```

---

## 🔑 Rule

👉 Use `Query` when you expect **many rows**

---

# 4️⃣ UPDATE

```go
_, err := h.DB.Exec(
	context.Background(),
	`
	UPDATE users
	SET username = $1
	WHERE id = $2
	`,
	newUsername,
	userID,
)
```

You can check affected rows:

```go
cmd, _ := h.DB.Exec(...)
if cmd.RowsAffected() == 0 {
	// nothing updated
}
```

---

# 5️⃣ DELETE

```go
_, err := h.DB.Exec(
	context.Background(),
	`DELETE FROM users WHERE id = $1`,
	userID,
)
```

---

# 6️⃣ INSERT + RETURNING (🔥 Postgres feature)

VERY common pattern.

```go
var userID string

err := h.DB.QueryRow(
	context.Background(),
	`
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
	`,
	req.Username,
	req.Email,
	hashedPassword,
).Scan(&userID)
```

This:

* Inserts row
* Returns generated ID
* All in one query

---

# 7️⃣ Where should SQL live?

### Option A: Inline SQL (what you’re doing now)

✅ Simple
✅ Fine for small projects

```go
"SELECT id FROM users WHERE email=$1"
```

---

### Option B: SQL as constants (cleaner)

```go
const getUserByEmail = `
	SELECT id, password
	FROM users
	WHERE email = $1
`
```

---

### Option C: `.sql` files (advanced)

```go
sqlBytes, _ := os.ReadFile("sql/get_user.sql")
query := string(sqlBytes)
```

Used in large codebases.

---

# 8️⃣ VERY IMPORTANT: placeholders difference

| DB         | Placeholder |
| ---------- | ----------- |
| PostgreSQL | `$1, $2`    |
| MySQL      | `?`         |
| SQLite     | `?`         |

You’re using Postgres → `$1`.

---

## 🔐 Security rule (burn this in)

❌ NEVER do this:

```go
"SELECT * FROM users WHERE email = '" + email + "'"
```

✅ ALWAYS do:

```go
"WHERE email = $1", email
```

---

## 🧠 Full login flow (realistic)

```go
err := h.DB.QueryRow(
	context.Background(),
	"SELECT password FROM users WHERE email=$1",
	req.Email,
).Scan(&hash)

if err == pgx.ErrNoRows {
	return unauthorized
}

if !PasswordMatches(req.Password, hash) {
	return unauthorized
}
```

---

## 🔑 Cheat sheet (save this)

| Operation   | Function    |
| ----------- | ----------- |
| INSERT      | `Exec`      |
| UPDATE      | `Exec`      |
| DELETE      | `Exec`      |
| SELECT 1    | `QueryRow`  |
| SELECT many | `Query`     |
| Return ID   | `RETURNING` |

---

## One-liner to remember

> **SQL is a string.
> pgx executes it.
> Go scans the result.**

You’re now officially past “tutorial backend” territory 😈
Next natural steps:

* transactions
* repositories
* pagination
* joins
* indexes

Say what you want next 👌
