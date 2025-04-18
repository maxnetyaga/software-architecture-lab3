#!/bin/bash

# Reset the drawing state
curl -X POST -d "reset" http://localhost:17000/

# Draw initial figure (Cross) on a white background
curl -X POST -d "white
figure 0.5 0.5
update" http://localhost:17000/

echo "Drawing initial figure. Starting diagonal movement..."

# Move the figure diagonally
x=0.5
y=0.5
dx=0.01
dy=0.01

while true; do
  # Calculate new coordinates (ensure within 0..1 bounds)
  # Using awk for floating point arithmetic
  x=$(awk "BEGIN { \$x = $x + $dx; if (\$x > 1 || \$x < 0) \$x = $x - $dx; print \$x }")
  y=$(awk "BEGIN { \$y = $y + $dy; if (\$y > 1 || \$y < 0) \$y = $y - $dy; print \$y }")


  # Simple boundary check and reverse direction if needed (considering figure size)
  # Approximate figure half size is around 0.05 in relative coordinates
  if (( $(awk "BEGIN {print ($x + 0.05) > 1}") )) || (( $(awk "BEGIN {print ($x - 0.05) < 0}") )); then
      dx=$(awk "BEGIN {print -$dx}")
  fi
   if (( $(awk "BEGIN {print ($y + 0.05) > 1}") )) || (( $(awk "BEGIN {print ($y - 0.05) < 0}") )); then
      dy=$(awk "BEGIN {print -$dy}")
  fi

  # Recalculate with potentially reversed direction
   x=$(awk "BEGIN { \$x = $x + $dx; print \$x }")
   y=$(awk "BEGIN { \$y = $y + $dy; print \$y }")


  # Send move and update commands
  curl -X POST -d "move ${dx} ${dy}
update" http://localhost:17000/

  sleep 1 # Wait
done