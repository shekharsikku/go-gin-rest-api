## **REST API - Using Go Gin**

Create `.env` file and add environment variables

```bash
JWT_SECRET=""
PORT=""
DB_PATH=""
```

### **Run database migration**

```bash
go run cmd/migrate/main.go <up|down> <db_path>
```

### **Run in debug mode**

This command will use air.toml configuration and run main.go file

```bash
air
```

### **Run using go run command**

```bash
go run cmd/api/main.go
```

---
