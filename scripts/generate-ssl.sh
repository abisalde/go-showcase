#!/bin/bash

set -e

echo "🔐 Generating SSL certificates for gRPC..."

cd ssl

# Output files
# ca.key: Certificate Authority private key file
# ca.crt: Certificate Authority trust certificate
# server.key: Server private key
# server.csr: Server certificate signing request
# server.crt: Server certificate signed by the CA
# server.pem: Server key in PEM format for gRPC

SERVER_CN="localhost"
MY_SUBJECT="/CN=$SERVER_CN"
PASSWORD="1111"

# Step 1: Generate Certificate Authority + Trust Certificate
echo "📝 Generating CA..."
openssl genrsa -passout pass:$PASSWORD -des3 -out ca.key 4096
openssl req -passin pass:$PASSWORD -new -x509 -sha256 -days 365 -key ca.key -out ca.crt -subj "/CN=ca"

# Step 2: Generate the Server Private Key
echo "🔑 Generating server private key..."
openssl genrsa -passout pass:$PASSWORD -des3 -out server.key 4096

# Step 3: Create OpenSSL config with SANs
cat > server.ext <<EOF
subjectAltName=DNS:localhost,DNS:grpc-server,DNS:*.localhost,IP:127.0.0.1
EOF

# Step 4: Get a certificate signing request from the CA
echo "📋 Creating certificate signing request..."
openssl req -passin pass:$PASSWORD -new -key server.key -out server.csr -subj "$MY_SUBJECT"

# Step 5: Sign the certificate with the CA (self-signing) with SANs
echo "✍️  Signing certificate..."
openssl x509 -req -sha256 -days 365 -passin pass:$PASSWORD -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt -extfile server.ext 

# Step 6: Convert the server certificate to .pem format
echo "🔄 Converting to PEM format..."
openssl pkcs8 -topk8 -nocrypt -passin pass:$PASSWORD -in server.key -out server.pem

# Cleanup temporary files
rm -f server.ext

echo "✅ SSL certificates generated successfully!"
echo ""
echo "Files created:"
echo "  - ca.crt (share with clients)"
echo "  - ca.key (keep private)"
echo "  - server.crt (server certificate - valid for localhost, grpc-server)"
echo "  - server.key (server private key)"
echo "  - server.pem (server key for gRPC)"
echo "  - server.csr (certificate signing request)"
