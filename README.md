# easycert
easycert is tool help you create certificate easily.
# How to use
## Create a certificate authority
`
./easycert ca -n root --cn root-ca
`
## Create a certificate
`
./easycert cert -n cert --cn mysite.com --ca root --dns mysite1.com --dns mysite2.com
`
