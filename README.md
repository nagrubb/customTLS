# Custom TLS
This project is an example on how you can customize the authentication/authorization
parts of TLS in two different ways (v1 and v2) using Go.

# Custom Verification Handler (v1)
The example contained in v1 shows you how to completely bypass Go's TLS
authentication and authorization verification. I would view this example as a last
resort as it puts a lot of responsibility on your code to ensure the certificate
from the peer is actually valid. I would advice this method only if you are very
comfortable and knowledge about security and also have a thorough security review
especially on this particular function.

# Appending a Root Certification Authority (v2)
The example contained in v2 shows you how you can extend the system's "Root of Trust"
with a custom certificate that you've received out of band. The example's method
for "out of band" is just for demonstration purposes only and obviously is not secure,
but some true out of band ways to share such a certificate could be:
- If your server application source is encrypted, embedding a private key and public certificate
within it. Then embedded the public certificate in the client.
- If your server application has a key vault or equivalent, then storing/using the private key
from that key vault with the certificate already contained in the client.
- Sharing the certificate over another transport that's already been authorized and
authenticating.
- Retrieving the certificate from a different entity that's already been authenticated
and authorized. For instance a remote server.

# How To Run

Note the **sudo** when running test.sh as both v1 and v2 use 443 and that's a protected port
requiring elevated permission to bind to on localhost (at least on macOS). If either v1 or v2 succeeds,
it will output just "Success!". Otherwise, it will display the error encountered.

## v1
```
cd v1
chmod +x test.sh
sudo ./test.sh
```

## v2
```
cd v2
chmod +x test.sh
sudo ./test.sh
```
