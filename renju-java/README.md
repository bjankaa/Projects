# Renju Board Game

A Java desktop implementation of Renju, a two-player strategy board game
played on a 15×15 board.

The application provides a graphical interface for local two-player gameplay
and implements game logic, move validation, Renju-specific rule checking,
score tracking, and menu navigation.

## Features

- Local two-player gameplay
- 15×15 game board
- Turn-based stone placement
- Graphical user interface
- Automatic win detection
- Renju-specific forbidden-move detection
- Pause and resume functionality
- Rules screen
- Scoreboard for previous games
- File-based storage of game results
- New game and menu navigation

## Game Rules

Black makes the first move by placing a stone in the center of the board.
Players then alternate turns by placing stones on the board.

A player wins by creating an uninterrupted line of five stones.

The implementation also handles Renju-specific restrictions for the black
player, including:

- Overline
- Inline double
- Double three
- Double four

These conditions are checked automatically after moves.

## Game Logic

The game logic analyzes the board around each placed stone across multiple
directions.

Pattern matching is used to detect:

- Valid five-stone sequences
- Overlines
- Double threes
- Double fours
- Other restricted patterns

The application evaluates the current move after each turn and determines
whether the game should continue or end.

## Application Structure

The application separates gameplay, menu navigation, and score management
into different components.

### Game

Contains the main gameplay logic, including:

- Board state management
- Stone placement
- Move validation
- Pattern matching
- Win-condition detection
- Forbidden-move detection

### Menu

Handles application navigation, including:

- Starting a new game
- Opening the rules screen
- Returning to the main menu
- Exiting the application

### Score Management

Handles previous game results, including:

- Reading stored results
- Managing score data
- Displaying the scoreboard

## Technologies

- Java
- Object-Oriented Programming
- Desktop GUI
- File-based persistence