# Game Application


## Use Cases

### User
- User register with phone number
- Uer login with phone number and password

### Game
- Each gaem have a given number of questions
- difficulty levels are "Easy, Medium, Hard"
- Winner is determined by number of correct answers
- Game have categories: sport, history, etc


## Entity

- User:
    - ID
    - PhoneNumber
    - Avatar
    - Name

- Game:
    - ID
    - Category
    - Question List
    - Players

- Question:
    - ID
    - Question
    - Answer List
    - Correct Answer
    - Difficulty
    - Category