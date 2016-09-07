#!/bin/bash
# declare bullyme script 

#read from a file a specific file
#take a random number between 1 and the file's number of lines
random_line=$((1 + RANDOM % 10));
#show the file line
sed $random_line'q;d' useless_questions.txt;
#add further validations for file exists, real line randomic and so
