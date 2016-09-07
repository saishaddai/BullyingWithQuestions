#!/bin/bash
# declare bullyme script 

#read from a file a specific file
#take a random number between 1 and the file's number of lines
number_of_lines=`cat useless_questions.txt | wc -l`;
random_line=$((1 + RANDOM % $number_of_lines));
#show the file line
sed $random_line'q;d' useless_questions.txt;
#add further validations for file exists, real line randomic and so
