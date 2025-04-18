#!/bin/bash

# Reset the drawing state
curl -X POST -d "reset" http://localhost:17000/

# Create a black rectangle on a white background, then change background to green
curl -X POST -d "white
bgrect 0.25 0.25 0.75 0.75
green
figure 0.5 0.5
figure 0.6 0.6
update" http://localhost:17000/

echo "Sent commands to create green frame with figures."