#!/bin/bash

echo "Starting blinking rectangle animation..."

while true; do
  # Reset and draw white background with black rectangle
  curl -X POST -d "reset
white
bgrect 0.3 0.3 0.7 0.7
update" http://localhost:17000/
  sleep 0.5 # Час відображення прямокутника

  # Reset to default black background
   curl -X POST -d "reset
update" http://localhost:17000/
  sleep 0.5

done