# Book Management System

## Description

The Book Management System streamlines the retrieval, creation, updating, and deletion of book records, providing an efficient and organized way to manage book data.

## Setup

### 1. Cloning

```bash
# Clone the repository
git clone <repository_url>
cd <repository_directory>
```

### 2. Installation

#### Kafka Configuration

📂 **Create Kafka Topic (****`book_events`****)**:

```bash
kafka-topics --create --topic book_events --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

🔍 **Verify Topic Creation:**

```bash
kafka-topics --list --bootstrap-server localhost:9092
```

If the topic `book_events` appears in the list, the creation was successful!

#### Redis Configuration

⚙️ **Verify Redis is running:**

```bash
redis-cli ping
```

If Redis is running, the response should be `PONG`.

#### PostgreSQL Setup

🐘 **Start PostgreSQL Server:**

```bash
pg_ctl -D /usr/local/var/postgres start
```

⚙️ **Connect to the Database:**

```bash
psql postgres://postgres:postgres@localhost:5432/books-management-system?sslmode=disable
```

📂 **Verify Connection:** You should be able to run SQL queries on the `books-management-system` database.

### 3. Run Migrations and Start Server

```bash
# Create books table using migrations
make migrateup

# Start the server
make start
```

