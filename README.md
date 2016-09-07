(work in progress)

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
alias bullyme='./askme.sh' 
alias answerme='./answerme.sh $1'
alias runtest='./askme.sh'
$BULLYNG_PATH=/your/path/to/root/directory
```

* Run test
```bash
$ runtest
```

