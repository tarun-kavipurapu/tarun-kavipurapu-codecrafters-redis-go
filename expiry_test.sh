#!/bin/bash

# Send SET command with expiry
{
  printf '*5\r\n$3\r\nSET\r\n$7\r\ntestKey\r\n$9\r\ntestValue\r\n$2\r\nPX\r\n$1\r\n5000\r\n'
} | nc localhost 6379 &
NC_PID=$!  # Get the process ID of the `nc` process
sleep 1    # Wait for 1 second
kill $NC_PID  # Terminate the `nc` process

# Send GET command before expiry
{
  printf '*3\r\n$3\r\nGET\r\n$7\r\ntestKey\r\n'
} | nc localhost 6379 &
NC_PID=$!  # Get the process ID of the `nc` process
sleep 1    # Wait for 1 second
kill $NC_PID  # Terminate the `nc` process

# Wait for expiry (6 seconds)
sleep 6  # Simulate key expiration

# Send GET command after expiry
{
  printf '*2\r\n$3\r\nGET\r\n$7\r\ntestKey\r\n'
} | nc localhost 6379 &
NC_PID=$!  # Get the process ID of the `nc` process
sleep 1    # Wait for 1 second
kill $NC_PID  # Terminate the `nc` process

