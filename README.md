# Locator
Locator determines the number of airplanes in a given area.  
It consists of a server and a client part.  
Clients connect to server with websocket protocol. Each client sends its area parameters: longitude, latitude and radius.  
Server caches all airplanes from opensky-network.org and sends to clients number of planes in the area they specified. The value is sent after each cache update.
 
## To run an example:
```
docker-compose build
docker-compose up
```