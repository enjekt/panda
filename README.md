# Panda Vault
Important Note: Panda Vault will be adding Go-Kit with both REST and gRPC implementations soon.

Secure storage of credit card PANs using encipherment provided by the related shannon-engine project. 

Panda Vault has a REST front end for adding PANs and getting back a new token and pad while storing the token and padded pan (OTP) in the datatabase.

Panda Vault uses SQLLite to store the token and padded pan via GORM. The whole compiles to a single executable appropriate for stand alone use or for use in a Docker container. Later associated project may create deployable Docker images. 

## Basic Round Trip

The REST service is fairly straightforward and recevies a credit card pan and returns the generated Token and One-Time Pad. 

![image](https://github.com/enjekt/panda/assets/3209869/2f507a8c-4290-4893-8d29-6df76b3eb7ce)

To restore the credit card PAN for use, one sends the Token and associated Pad. 

![image](https://github.com/enjekt/panda/assets/3209869/83eb17c4-c8a5-4406-87ae-51f84063aa33)

## How it Works
The initial store of the PAN results a Token and One-Time Pad being generated for the PAN (more on that in a bit). The One-Time Pad is XOR'd to the Pan resulting in the PaddedPan (see the shannon engine dependency for more details.) The Token and PaddedPan are stored via GORM in the SQLite database. The Token and Pad are then returned to the caller who must store them for future use. When the Token and the Pad are sent back to the service, the PaddedPan is retrieved from the database and XOR'd with the One-Time Pad which restores the original PAN. The PAN is then returned to the caller who can use it along with the expiration date for transactions.

## Important!
It's important to note here that the credit card number, the PAN, is never stored as data at rest. The Token, PaddedPan and Pad are all useless numbers. The Token uses the BIN and the last digits of the PAN for identification purposes. But all the middle digits are randomly generated. Furthermore, the Token is verified to be a **non-valid** PAN (it fails the Luhn check). It is, therefore, always possible to determine if the number is a valid PAN or a Token. 

It is anticipated that an application that works with the Panda Vault will store the Token, Pad and the credit card expiration date. To make a charge, it will send the Token and Pad to the vault to reconstruct the PAN and then use it, along with the expiration date, for a transaction. After that, the system will once again return to the state where no valid credit card data is stored in a datbase nor are the constituents necessary to reconstruct the credit card stored in any single database. 

In our example, the original credit card PAN is 5513746525703556. The token is 5513740800553556. At first glance they appear the same but there's a critical difference. Note that the BIN, the first set of digits are the same and the last 4 digits are the same, but the middle digits are different. 

 
PAN                 | Token
------------------- | -------------
55137**4652570**3556|55137**4080055**3556


Further, if you run the first number through a Luhn check, it will pass. The token, by contrast, will fail.

## Basic Wiring

While more functionality will be added as the project moves forward, there are certain challenges which may not have canned solutions. For example, Visa and MasterCard have moved from 6 digit to 8 digit BIN numbers while retaining a 16 digit length. To compensate, they've reduced the number of storable last digits from 4 to 2. If you look at the baseline code in the Panda Vault, you'll see we construct a 6+4 pipeline based on the older style of credit card number.

		pipelines.NewPipeline().Add(pipelines.CompactAndStripPanFunc).Add(pipelines.CreatePadFunc).Add(pipelines.EncipherFunc).Add(pipelines.TokenFunc(6,4)))
  
Note the TokenFunc is passed in the parameters of 6 and 4 teling it how to parse and construct the token. If we changed that to 8,4 it would work with the newer credit card and we could even have pools of pipelines for both types. However, we'd have to have a mechanism for determining when to use the one or the other.

It is anticipated that the Panda Vault will run in Docker container and that orchestration services will handle the external security. However, it can run standalone in its own VM or on hardware. In addition, it will likely live in a subnet and there is rarely, if ever, a need for external services to connect to the Panda Vault. That is why the current implemetation fo the vault doesn't use SSH, certificates or even basic user/passwords as those will likely be tailored based on an installation by installation basis to meet the needs of individual clients. 
