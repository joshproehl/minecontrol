# Minecontrol

Minecontrol is a tool for interacting with a Minecraft server via it's RCON connection. Provided your server has RCON enabled,
and that you know the RCON password, Minecontrol will allow you to do the following, without actually being in the game console:
* Execute an arbitrary command and see the response
* Open a Read-Evaluate-Print-Loop shell, allowing you to enter multiple commands. This is basically just like the game console
  except that you do not see updates for things such as "player was killed by zombies".
* Create a web server which will server HTML pages displaying status for the server, and which provides a RESTful JSON API for
  interacting with the game's console.


This project is what happens when a programmer wants to be able to see who's logged in to his minecraft server, and then goes
completely off the rails.
