# Mini-Bank

Mini-Bank is a simple banking application written in Go. The project is developed in three different ways using various architectural approaches: **Corporate**, **Speedster**, and **Minimalist**. Each version is optimized for different needs and priorities.

## Project Overview

All versions of the project are developed for educational purposes. They are not production-ready applications. They may contain bugs, security vulnerabilities, and performance issues.

All versions of the project have the following features:
- **User Management:** User creation, editing, and deletion.
- **Account Management:** Account creation with different currencies, balance checking, and transaction history.
- **Transaction Management:** Deposit, withdrawal, and transfer between accounts.
- **Admin Management:** Admin login and user management.
- **Email Management:** Sending emails for admin passwordless login and transaction notifications.
- **Fee Management:** Fee creation and deduction from transactions.

## Project Structures

### 1. Corporate

**Features:**
- **Framework:** Gin
- **Dependency Injection:** Uber/fx
- **ORM:** Gorm
- **Logger:** Zap
- **Architecture:** Domain-Driven Design (DDD)
- **ID Management:** UUID
- **Monetary Calculations**: shopspring/decimal
- **Email:** Third-party library
- **Admin Login:** Redis and passwordless login via email
- **Documentation:** Fully documented with Swagger

**Purpose:** To provide a business-oriented and extendable structure.

### 2. Speedster

**Features:**
- **Framework:** Fiber
- **Database Driver:** pgx
- **Logger:** Zerolog
- **Dependency Injection:** Manual
- **Configuration Management:** Envconfig
- **ID Management:** TSUID
- **Monetary Calculations:** cockroachdb/apd
- **Architecture:** Clean Architecture
- **Email:** Third-party library
- **Admin Login:** Redis and passwordless login via email
- **Documentation:** Swagger comments added, but no routes created

**Purpose:** To provide a speed-oriented and high-performance structure.

### 3. Minimalist

**Features:**
- **Standard Libraries:** Go standard library
- **Logger:** Go standard logger
- **Database:** SQLite
- **Dependency Injection:** Manual
- **ID Management:** int64
- **Monetary Calculations:** math/big.Float
- **Architecture:** Simple Layered Architecture
- **Email:** Go standard library
- **Admin Login:** Redis and passwordless login via email
- **Documentation:** No documentation

**Purpose:** To provide a simple and minimal structure.