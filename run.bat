@echo off

echo Running tests...
go test ./...
if %ERRORLEVEL% neq 0 (
    echo Tests failed. Exiting.
    exit /b %ERRORLEVEL%
)

echo All tests passed. Running the project...
go run main.go