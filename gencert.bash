echo "creating a server.key"

openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
echo "creating server.crt"
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 -subj "/C=IN/ST=Karnataka/L=Bangalore/O=MyCompany/OU=IT/CN=arpitsrivastava529@gmail.com"
