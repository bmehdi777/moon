#!/bin/bash

echo "make server cert"
openssl req -new -nodes -x509 -out assets/server/certs/server.pem -keyout assets/server/certs/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=m00n/OU=IT/CN=m00n.fr/emailAddress=$1"
#echo "make client cert"
#openssl req -new -nodes -x509 -out assets/client/certs/client.pem -keyout assets/client/certs/client.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=$1"
