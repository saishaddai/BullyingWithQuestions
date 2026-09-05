BullyingWithQuestions
======
Originally this was a basic bash script i wrote 10 years ago to retrieve a random question from a file to help you study from a text plain file full of questions. Now I decided to move to Go + Bubble Tea. I will keep updated this readme in order to improve step by step this old project I started with only good will but no time. 

## Gitflow

You must create a fork in order to contribute.

## Setup

You may want to set a couple of commands as alias in order to run each time you open a terminal

### Requirements

* Go
* Bubble tea
* SQLite (to store the information)

### Steps

* Navigate to root directory
* Run the following command:

```bash
$ go main.go
```

This will send a message of all the possible errors. Please refer to me in case troubleshot 

### Screenshots
This is a WIP and since it is a Text based user interface, it would be easier to share the screenshot as I'm building the app by panels

### Glossary 
- Deck: topics, basically topics. Since this is a set of flashcards to study, it means a deck is a category to group the flashcards
- Flashcard: a combination of question and a answer. So far, the information is stored in a local database. I think it should have a configuration file so it can connect to different sources to get the flashcards. 
- Session: A study session. It consist of a random set of 20 flashcards to study
- Navigation: The User can navigate through all the flashcards back and forth until they finishes the whole 20 cards

### AI involved in this project
I will use Supacode, Open Code and free LLMs as possible. Since Im kind new n Go, some issues may happen. Even so, this wil be my first formal TUI after a quick learning path with bubble tea




