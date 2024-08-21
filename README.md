# Building a Online Payment Platform

## Background
With the rapid expansion of e-commerce, there is a pressing need for an efficient payment gateway. This project aims to
develop an online payment platform, which will be an API-based application enabling e-commerce businesses to securely
and seamlessly process transactions.

## Entities Involved:
1. Customer: Individuals who make online purchases and complete payments through the platform.
2. Merchant: The seller who utilizes the payment platform to receive payments from customers.
3. Online Payment Platform: An application that validates requests, stores card information, and manages payment
   requests and responses to and from the acquiring bank.
4. Acquiring Bank: Facilitates the actual retrieval of funds from the customer's card and transfers them to the merchant.
   Additionally, it validates card information and sends payment details to the relevant processing organization.

## Requirements:
The requirements for this initial phase are as follows:
1. Payment Processing:
   The online payment platform should provide merchants with the ability to process a payment and receive either a
   successful or unsuccessful response.
2. Querying Details of Previous Payments:
   Merchants should be able to retrieve details of previously made payments using a unique payment identifier.
3. Bank Simulation:
   Utilize a bank simulator to simulate the interaction with the acquiring bank.

## Deliverables:
1. Build an API that allows merchants:
   a. To process a payment through the online payment platform.
   b. To retrieve details of a previously made payment.
   c. To process refunds for specific transactions.
2. Integrate a Bank Simulator:
   Use a bank simulator to test and simulate responses from the acquiring bank.

## Considerations:
- Execution of the Solution:
    - Provide clear instructions for setting up and running the online payment platform API.
    - Specify any dependencies or prerequisites for the solution.
- Assumptions:
  - Clarify any assumptions made during the design and implementation.
  - Areas for Improvement:
  - Identify potential areas for improvement in the online payment platform.
  - Discuss any design decisions or trade-offs made during development.
- Cloud Technologies:
  - Specify the cloud technologies used and justify the choice.

## Extra:
- Authentication and Security:
  - Implement measures for authentication and security to ensure secure transactions.
- Audit Trail:
    - Include an audit trail feature to track activities such as payment processing, queries for payment details, and
  refunds.

---

# Challenge Documentation

## First steps
There's a docker-compose file. You can run it with `docker-compose up --build`.
Once every service is up and running, you can check the logs with `docker-compose logs -f`.

## Setting up different scenarios
If you want to modify your experience like if there was any issue in the app or the bank application you can try editing the `dockerenv.json` or `env.json` file.

- `AUTH_TOKEN_IS_VALID`: Here's the happy path, the used is "logged in" 
- `CLIENT_HAS_ENOUGH_BALANCE`: The bank will return an error because the client doesn't has enough balance if set to false
- `CARD_HASH_IS_VALID`: We won't be sending card information, the bank should be able to validate if the CC is valid or not with a hash
- `CLIENT_HAS_EXCEEDED_LIMIT`: The client has exceeded the limit and the bank should return an error
- `BANK_TX_FAILED`: The bank request failed

## Testing the application
I have created a swagger file that you can read it through the swagger UI in `http://localhost:3000` if the docker container is running.

## IMPORTANT NOTES
The application has 3 banks loaded, 10 merchants and 10 customers. The 3 banks are "hardcoded", the other 20 rows are generated with random data.
The customer id will be random generated in the auth middleware and is a number between 1 and 10.

## Project structure
There are 2 folders in the root:
- application: this is where the application code is located
- swagger: this is where the swagger file is located

### Application folder:
This is where the application code split into 3 different type of folders:
- Shared code: everything that's not into bank-app or payments-app then it is code that lives there so it can be reused easily
- Bank APP: Mock for bank simulator
- Payments APP: The actual application

### Payments APP structure
- cmd: this is where the main.go file lives
- database: handles the database connection and some helpers
- domain: this is where the domain-specific files are stored
- payment: this is the where the business logic is stored, you can find the handler, service and repository there
- bank: bank repository, it is used to interact with the bank simulator
- idempotency: helper for idempotency
- http: all http server related

### Bank APP structure
- cmd: this is where the main.go file lives
- http: all http server related
- bank: this is where the business logic is stored

### Areas for improvement
- Add more unit test and some integration tests, I wouldn't deploy an application without AT LEAST 90% coverage
- Add security. There's no actual security in the application
- *Wouldn't push the `.env` file*, just an `.env.example`
- Maybe I would have considering not using an ORM
- Add more documentation in general
- Observability
- Retry if the bank request fails, probably implementing an exponential backoff and a circuit breaker
- There's -of course- a broken access control. I could check for payments for a given customer even if I'm not and admin or that customer


### Things that I consider that are interesting
- I have added a Idempotency Key* check
- The actual codes that I'm using are the real ones

#### * Idempotency Key
Services like Paypal, Stripe or MercadoPago uses an idempotency key to avoid duplicate transactions.
Most of them store a unique key in a (usually) fast NoSQL database and with a short TTL. If the client sends 2 requests with the same idempotency key, then probably it's a duplicated one.

### Architecture
This is probably the most debatable and interesting part of the project.
What I would do if this were a cloud service?

**In AWS I would need:**
- API Gateway, so we could handle requests and scale per needed
- (depending on the size of the application) The classic ECS - EKS service for the application
- Managed Cache for the idempotency key
- RDS for the database
- Kinesis for logging aggregation
- S3 for storing log files
- ELB

**On the CI/CD side:**
- Gitub actions that would check
  - Security vulnerabilities
  - Run the tests
  - Build the code
  - General code quality and static analysis tools
  - Deploy the application to AWS (at least to development - staging envs)

**Observability stack:**
- DataDog
- NewRelic
- ELK

### Assumptions
- I'm always assuming that the bank always will be able to refund a payment and to do a reversal
- I'm assuming that the bank will deposit the money into the merchant account and withdraw it from the customer account
- I'm assuming that we're communicating through a secure network + authenticated users + encryption

### Data models
I'm using Bun, which is great but the tables might not be as obvious as expected. You can find them in `./domain/database`
Every struct will be turned into a table (if it not exists).

### cURLs

#### Bank

###### GET - /ping 
```curl
curl --location 'localhost:8888/ping'
```

###### POST - /pay 
```curl
curl --location 'localhost:8888/pay' \
--header 'Content-Type: application/json' \
--data '{
    "amount": 100,
    "merchant_id": 1,
    "bank_id": 1,
    "card_hash": "test"
}'
```

#### Payments

###### GET - /ping
```curl
curl --location 'localhost:8080/ping'
```

###### POST - /pay 
```curl
curl --location 'localhost:8080/pay' \
--header 'Content-Type: application/json' \
--data '{
    "amount": 100,
    "merchant_id": 1,
    "bank_id": 1,
    "card_hash": "test"
}'
```

##### PUT - /payments/{payment_id}/refund
```curl
curl --location --request PUT 'localhost:8080/payments/1/refund' \
--header 'Authorization: user-token' \
```

##### GET - /payments (fetch every payment)
```curl
curl --location 'localhost:8080/payments' \
--header 'Authorization: user-test'
```

##### GET - /payments/{payment_id} (fetch a given payment id)
```curl
curl --location 'localhost:8080/payments/66' \
--header 'Authorization: user-test'
```

##### GET - /customers/{customer_id}/payments (fetch every payment for a given customer)
```curl
curl --location 'localhost:8080/customers/4/payments' \
--header 'Authorization: user-test'
```