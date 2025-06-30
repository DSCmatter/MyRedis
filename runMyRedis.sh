# Stop existing Redis if running via snap
sudo snap stop redis

# Navigate into src directory
cd src/

# Run MyRedis server in background
go run *.go &

# Wait a bit for the server to start
sleep 2

# Connect to MyRedis server using redis-cli
redis-cli 