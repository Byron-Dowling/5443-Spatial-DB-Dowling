"""
    Name: Byron Dowling
    Class: 5443 Graduate Spatial Databases
    
    File Description:
      This file stops the missile defense game and calls the missile
      server with our ID removing us from the list of targets and 
      returning the region for another team to use.
"""

import json
import requests

if __name__ == "__main__":

    # Quit !!!
    # get the id from myregion.json
    with open("myregion.json") as f:
        data = json.load(f)
        id = data["id"]

    requests.get("http://missilecommand.live:8080/QUIT/" + str(id))
    print("Wave the white flag, region " + str(id) + " surrenders!")
