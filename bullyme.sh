#!/bin/bash
# declare bullyme script 

#read from a file a specific line
#a random line by the way
random_line = $((1 + RANDOM % 10))
#and simply show it
sed '$random_lineq;d' file

#add further validations for file exists, real line randomic and so
