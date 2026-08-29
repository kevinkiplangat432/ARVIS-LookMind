#!/bin/sh
set -e

# Run database migrations
./arvis migrate up

# Start the server (replace shell with server process so signals work)
exec ./arvis server