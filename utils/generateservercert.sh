#!/bin/bash

# if the directory does not exist
if [ ! -d ../private ]; then
	mkdir ../private
fi

# file paths
config=../private/certificate.config
server_cert=../private/server.cert.pem
server_key=../private/server.key.pem

#if the configutation file exist
if [ -f $config ]; then
	echo generating the server certificate using the private configuration...
	openssl req -new -nodes -x509 -out $server_cert -keyout $server_key -days 365' < $config
else
	echo generating the server certificate using the user input...
	openssl req -new -nodes -x509 -out $server_cert -keyout $server_key -days 365'
fi
