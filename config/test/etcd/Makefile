ca_key = ca.key
ca_crt = ca.crt
host_key = server.key
host_csr = server.csr
host_crt = server.crt
conf = openssl.conf

default: $(host_crt)

$(ca_key):
	openssl genrsa -out $(ca_key) 2048

$(ca_crt): $(ca_key)
	openssl req -x509 -new -nodes -key $(ca_key) \
		-subj "/CN=CA" \
		-sha256 -days 365 -out $(ca_crt)

$(host_key):
	openssl genrsa -out $(host_key) 2048

$(host_csr): $(host_key) $(conf)
	openssl req -new -key $(host_key) -out $(host_csr) -config $(conf)

$(host_crt): $(host_csr) $(ca_crt) $(ca_key)
	openssl x509 -req -in $(host_csr) -CA $(ca_crt) -CAkey $(ca_key) \
		-CAcreateserial -out $(host_crt) -days 365 -extensions v3_req \
		-extfile $(conf)
clean:
	rm -f $(ca_key) $(ca_crt) ca.srl
	rm -f $(host_key) $(host_csr) $(host_crt)
