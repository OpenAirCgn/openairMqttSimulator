for PEM in `cat sha_client_id.txt`; do
	echo $PEM
	openssl rsa -in ${PEM}.pem -passin pass:${PEM} -out tmp.pem
	mv tmp.pem ${PEM}.pem
done
