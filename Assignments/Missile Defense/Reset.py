"""
    Name: Byron Dowling
    Class: 5443 Graduate Spatial Databases
    
    File Description:
      This file resets the game server. Is necessary for testing
      and when all the regions have been occuipied.
"""

import requests

if __name__ == "__main__":

    # "RESET THE GAME!!!
    # Send a request to reset the game.
    requests.get("http://missilecommand.live:8080/RESET")
    print("Reset the Game ...")
