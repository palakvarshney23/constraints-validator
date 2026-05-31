@echo off
cd /d "%~dp0"
where go >nul 2>&1
if errorlevel 1 (
  echo Go is not installed. Install from https://go.dev/dl/ then run again.
  pause
  exit /b 1
)
go run .
pause
