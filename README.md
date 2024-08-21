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