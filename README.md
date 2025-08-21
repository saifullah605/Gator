# Gator - RSS Feed Blog Aggregator (CLI in Go)

Gator is a simple **RSS Feed Blog Aggregator** built in Go.
It allows you to follow feeds, aggregate blog posts, and manage them directly from the command line.

---

## Requirements

Before getting started, make sure you have the following installed:

* **Go** (>= 1.21 recommended)
  [Download](https://go.dev/dl/)

  ```bash
  # Verify installation
  go version
  ```

* **PostgreSQL** (>= 14 recommended)
  [Download](https://www.postgresql.org/download/)

  ```bash
  # Example installation on Ubuntu/Debian
  sudo apt update
  sudo apt install postgresql postgresql-contrib

  # Start PostgreSQL
  sudo service postgresql start

  # Verify installation
  psql --version
  ```

---

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/gator.git
cd gator
```

### 2. Create a Config File

Create the config file in your home directory:

```json
~/.gatorconfig.json
{
  "db_url": "postgres://example"
}
```

This will store your database connection string.

---

## PostgreSQL Setup

1. Enter PostgreSQL shell:

```bash
sudo -u postgres psql
```

2. Create the database:

```sql
CREATE DATABASE gator;
```

3. Create or alter the user (example credentials):

```sql
ALTER USER postgres PASSWORD 'postgres';
```

**Default example setup:**

* Database: `gator`
* Username: `postgres`
* Password: `postgres`

---

## Install Goose (for Migrations)

We use **Goose** to manage database migrations: [Goose GitHub](https://github.com/pressly/goose)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Make sure `go/bin` is in your `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Run Database Migrations

Navigate to the schema directory:

```bash
cd sql/schema
```

Connection string examples:

* With username and password:
  `postgres://postgres:postgres@localhost:5432/gator?sslmode=disable`

* With username but no password:
  `postgres://postgres@localhost:5432/gator?sslmode=disable`

* With no username/password (trust mode):
  `postgres://localhost:5432/gator?sslmode=disable`

Run migrations:

```bash
goose postgres "<connection_string>" up
```

Example:

```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

---

## Build the Project

From the project root:

```bash
go build
```

This will create an executable named `gator`.

---

## Usage

Run the CLI tool:

```bash
./gator <command>
```

Example:

```bash
./gator login
```

---

## Available Commands

* `login`       – Log in as a user
* `register`    – Register a new user
* `users`       – List all users
* `agg`         – Run the feed aggregator
* `addfeed`     – Add a new RSS feed (requires login)
* `feeds`       – List all feeds
* `follow`      – Follow a feed (requires login)
* `following`   – Show feeds you are following (requires login)

---

## Notes

* Make sure PostgreSQL is running before using the CLI.
* Update your `~/.gatorconfig.json` with the correct database URL.
* If you run into connection issues, double-check your Postgres user, password, and host.
