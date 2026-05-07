# API Service Documentation

## Getting Started

- Prerequisites for running the API
- Installation steps specific to the API service
- Environment variable setup for the API

## Running the API

- Starting the development server
- Running in production mode

## Database Migration Guide

This document outlines the process for managing database migrations in this project. Follow these steps to ensure your database schema stays up-to-date and consistent across environments.

### Migration Files

- All migration files are located in the `migrations/` directory.
- Each migration consists of an `.up.sql` file (for applying changes) and a corresponding `.down.sql` file (for rolling back).

### Running Migrations

1. **Ensure your database is running and accessible.**
2. **Apply migrations:**
   - Use your migration tool (e.g., `golang-migrate`, `goose`, or another) to apply all pending migrations.
   - Example with `golang-migrate`:

     ```
     cd /Users/theppitak/Coding/Playground/web-playground/lotto-journal/apps/api
     migrate -path migrations -database <DB_URL> up
     ```

   - Make sure you run this command from the directory where the `migrations/` folder is located (usually the project root or `apps/api`). Adjust the path if running from elsewhere.

3. **Rollback migrations (if needed):**
   - To undo the last migration:

     ```
     migrate -path migrations -database <DB_URL> down 1
     ```

   - Again, ensure your working directory is correct so the `migrations/` path is valid.

### Creating New Migrations

1. Create new `.up.sql` and `.down.sql` files in the `migrations/` directory.
2. Name them with an incrementing prefix (e.g., `000002_add_table.up.sql` and `000002_add_table.down.sql`).
3. Write the SQL statements for upgrading and downgrading the schema.

### Tips

- Always test your migrations on a development database before applying to production.
- Keep migration files under version control.

## Testing

- How to run tests for the API

## API Endpoints

- Overview or link to API documentation (e.g., Swagger/OpenAPI)

## Troubleshooting

- Common issues and solutions for the API service

## Contributing

- Contribution guidelines for the API repo

## License

- License information for the API service
