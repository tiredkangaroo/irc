cool to have besides just be an irc client:

- slack emojis/custom emojis
- different colors for different users
- notifications on mentions/pings
- games between two users (like it creates a seperate channel as a little pipe for the game and data is transmitted as messages, tic tac toe, hangman, etc)
- see multiple channels at once
- when disconnecting, have a server store messages for you and deliver them when you reconnect

how this is going to work:

- a web app connects to my backend server via websockets
- the backend server connects to irc servers via irc protocol
- the backend server relays messages between the web app and the irc server
- the backend server handles any extra features (like games, storing messages, etc)
