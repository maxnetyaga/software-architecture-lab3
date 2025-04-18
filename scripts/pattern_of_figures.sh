#!/bin/bash

# Reset the drawing state
curl -X POST -d "reset" http://localhost:17000/

# Set background color (e.g., white)
curl -X POST -d "white" http://localhost:17000/

# Draw figures in a pattern (e.g., corners and center)
curl -X POST -d "figure 0.1 0.1
figure 0.1 0.9
figure 0.9 0.1
figure 0.9 0.9
figure 0.5 0.5
update" http://localhost:17000/

echo "Sent commands to draw a pattern of figures."