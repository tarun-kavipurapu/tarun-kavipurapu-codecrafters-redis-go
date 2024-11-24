# Set key with expiry of 5 seconds
(printf '*5\r\n$3\r\nSET\r\n$7\r\ntestKey\r\n$9\r\ntestValue\r\n$2\r\nPX\r\n$1\r\n5\r\n';) | nc localhost 6379

# Get the value before expiry
(printf '*2\r\n$3\r\nGET\r\n$7\r\ntestKey\r\n';) | nc localhost 6379

# Wait for expiry
sleep 6

# Try to get the value after expiry
(printf '*2\r\n$3\r\nGET\r\n$7\r\ntestKey\r\n';) | nc localhost 6379






