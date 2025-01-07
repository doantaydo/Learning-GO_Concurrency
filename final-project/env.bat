@echo off
set DSN=host=localhost port=5432 user=postgres password=24072001do dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5
set REDIS=127.0.0.1:6379