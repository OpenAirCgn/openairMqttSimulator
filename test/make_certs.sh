for CLIENT_ID in `cat sha_client_id.txt`; do
	openssl req -newkey rsa:2048 \
		-keyout ${CLIENT_ID}.pem \
		-out ${CLIENT_ID}.csr \
		-subj /C=DE/O=PEK/CN=${CLIENT_ID}


done
