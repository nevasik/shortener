# URL Shortener
# URL Shortener is a simple Golang project designed to shorten long URLs.

# Installation
Make sure you have Golang installed. If not, you can download it here.

Clone the repository:
bash
Copy code
```bash
git clone https://github.com/nevasik/shortener.git
cd url-shortener
```
# Install dependencies:

bash
Copy code
```bash
go get -v
```
# Database Setup
The project uses migrations to manage the database structure. You can run migrations with the command:

go
Copy code
```bash
migrate create -ext sql -dir database/migration/ -seq init_mg
migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose up
```
Running the Application
Execute the following command to run the application:

bash
Copy code
```bash
go run main.go
```
The application will be accessible at http://localhost:8080.

Testing
To run tests, use the following command:

bash
Copy code
```bash
go test -v ./...
```
# Middleware Usage
The project uses middleware for various tasks, such as logging and error handling. You can configure new middleware in the middleware.go file.

# Logging
The project utilizes the slogger library for logging. Logs will be written to the logs/app.log file. You can configure the logging level in the config/config.go file.

# Contribution
If you have suggestions for improving the project, please create an Issue or Pull Request.

# License
This project is licensed under the terms of the MIT License.
