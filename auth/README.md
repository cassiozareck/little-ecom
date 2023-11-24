### JWT Token Validation

The Auth Service uses JSON Web Tokens (JWT) to manage user authentication. When a user signs in, their credentials are verified, and a JWT is generated and signed with a secret key. This token contains encoded user information and an expiration time. For each subsequent request to a protected endpoint, the client must include this token in the `Authorization` header.

The `ValidateToken` function is responsible for parsing and validating the JWT. It checks the token's signature to ensure it was signed with the correct key and verifies that it has not expired. If the token is valid, the user's email and token expiration time are returned in the response, allowing the client to confirm the user's identity and session validity.

### Secure Password Storage with bcrypt

To ensure the security of user passwords, the Auth Service uses bcrypt to hash passwords before storing them in the database. When a user registers, their plain-text password is hashed using bcrypt's `GenerateFromPassword` function, which incorporates a salt to protect against rainbow table attacks. The hashed password is then stored in the database.

During the sign-in process, the provided password is hashed and compared to the stored hash using bcrypt's `CompareHashAndPassword` function. This method securely verifies the user's password without ever storing or transmitting the plain-text version.

### Kubernetes Manifests

The Auth Service is deployed on Kubernetes using a set of manifest files. These manifests define the necessary Kubernetes resources, such as Deployments, Services, and Secrets. The `auth-service.yaml` file specifies the service configuration, exposing the Auth Service within the Kubernetes cluster. The `auth.yaml` file defines the deployment, including the container image to use, environment variables from secrets, and the number of replicas.

### Kubernetes Integration

Kubernetes integration is achieved through the use of environment variables and secrets, which are defined in the `secret.yaml` file. These secrets provide sensitive configuration options, such as the JWT secret key and database credentials, without exposing them in the source code or Docker image.

The Auth Service's Docker image is built and pushed to a container registry using the `deploy.sh` script. Kubernetes then pulls the image from the registry and deploys it according to the specifications in the manifest files. This allows for a seamless CI/CD pipeline and easy updates to the service.

### Security

Security is a key aspect of the backend API. All sensitive endpoints are protected with JWT-based authentication, ensuring that only authorized users can access them. The communication between services is also secured to prevent unauthorized access.

For more detailed information about each service, refer to their respective README files and documentation.
