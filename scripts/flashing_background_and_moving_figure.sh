#!/bin/bash

# Reset the drawing state
curl -X POST -d "reset" http://localhost:17000/

# Draw initial figure (Cross) at a starting position
curl -X POST -d "white
figure 0.1 0.1
update" http://localhost:17000/

echo "Drawing initial figure. Starting flashing background and movement..."

x=0.1
y=0.1
dx=0.015 # Зміщення
dy=0.015

while true; do
  # Calculate new displacement for the next step
  # Simple boundary check and reverse direction if needed
   if (( $(awk "BEGIN {print ($x + $dx + 0.05) > 1}") )) || (( $(awk "BEGIN {print ($x + $dx - 0.05) < 0}") )); then
       dx=$(awk "BEGIN {print -$dx}")
   fi
    if (( $(awk "BEGIN {print ($y + $dy + 0.05) > 1}") )) || (( $(awk "BEGIN {print ($y + $dy - 0.05) < 0}") )); then
       dy=$(awk "BEGIN {print -$dy}")
   fi

   # Update current position for boundary check in next iteration (optional, move takes displacement)
   x=$(awk "BEGIN { print $x + $dx }")
   y=$(awk "BEGIN { print $y + $dy }")


  # Post move command (updates state, doesn't render)
  curl -X POST -d "move ${dx} ${dy}" http://localhost:17000/

  # Post white background and update (renders state with white background)
  curl -X POST -d "white
update" http://localhost:17000/
  sleep 0.1 # Короткий інтервал для миготіння

  # Post green background and update (renders state with green background)
  curl -X POST -d "green
update" http://localhost:17000/
  sleep 0.1 # Короткий інтервал для миготіння

done