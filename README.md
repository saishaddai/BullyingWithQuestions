(work in progress, not final version)

BullyingWithQuestions
======
this is a basic bash script to retrieve a random question from a file to help you study from a text plain file full of questions

## Gitflow

You must create a fork in order to contribute.

## Setup

You may want to set a couple of commands as alias in order to run each time you open a terminal

### Requirements

* bash

### Steps

* Navigate to root directory
* Run the following command:

```bash
$ ./install.sh
```
* If something fails, then do the setup manually. Open your .bashrc file and type the following
```bash
alias bullyme='./bullyme.sh' 
alias answerme='./answerme.sh $1'
alias runtest='./runtest.sh'
$BULLYNG_PATH=/your/path/to/root/directory
bullyme
```

* Run test
```bash
$ runtest
```
This will send a message of all the possible errors. Please refer to me in case troubleshot 

### Add new Questions file
* Open the setup.sh file
* Add the path to the questions file (it should be a txt file)
* Add the path to the answers file (it should be a txt file)
* Save and close the editor
* Run the run test command and see whether the new configurations are set rightly



