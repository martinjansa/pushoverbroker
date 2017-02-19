#!/bin/bash

# if the directory does not exist
if [ ! -d ../private ]; then
	mkdir ../private
fi

# construct the generation command
cmd='openssl req -new -nodes -x509 -out ../private/server.pem -keyout ../private/server.key -days 365'

# if the configutation file exist
if [ -f ../private/certificate.conf ]; then
	echo generating the server certificate using the private configuration...
	$cmd < ../private/certificate.conf
else
	echo generating the server certificate using the user input...
	$cmd
fi

