# Panda Vault

Secure storage of credit card PANs using encipherment provided by the related shannon-engine project. 

Panda Vault has a REST front end for adding PANs and getting back a new token and pad while storing the token and padded pan (OTP) in the datatabase.

Panda Vault uses SQLLite to store the token and padded pan via GORM. The whole compiles to a single executable appropriate for stand alone use or for use in a Docker container. Later associated project may create deployable Docker images. 

