# Weather Subscription Service

## Architecture

The system uses a microservices architecture with three core services:

- **`subscription-service`**: Handles user subscriptions, fetches weather data from `weather-service` using gRPC, and sends email notifications via RabbitMQ.
- **`weather-service`**: Retrieves weather data from an external API and serves it via gRPC.
- **`email-service`**: Consumes RabbitMQ messages to send email notifications.

## Microservices

1. **subscription-service**:
   - **Ports**: HTTP (`8080`), gRPC client.
   - **Dependencies**: `subscription-db`, `rabbitmq`, `weather-service`.

2. **weather-service**:
   - **Ports**: gRPC (`50052`).
   - **Dependencies**: `weather-db`.

3. **email-service**:
   - **Dependencies**: `rabbitmq`.

4. **frontend**:
   - **Ports**: HTTP (`80`).
   - **Dependencies**: `subscription-service`.

![microservices.png](microservices.png)

## How to Run

### Prerequisites
- Docker and Docker Compose.

### Steps
1. **Clone the repo**:
   ```bash
   git clone <your-repo-url>
   cd weather-subscriber
   ```

2. **Configure**:
   - Edit `email-service/config/config.yaml` with SMTP credentials.

3. **Run**:
   ```bash
   docker-compose up --build
   ```
   - Frontend: `http://localhost:80`.
   - API: `http://localhost:8080`.