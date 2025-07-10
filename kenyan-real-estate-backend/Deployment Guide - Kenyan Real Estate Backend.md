# Deployment Guide - Kenyan Real Estate Backend

This guide provides step-by-step instructions for deploying the Kenyan Real Estate Backend in different environments.

## Prerequisites

- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **PostgreSQL 12+**: [Download PostgreSQL](https://www.postgresql.org/download/)
- **Git**: For cloning the repository
- **M-Pesa Developer Account**: [Safaricom Developer Portal](https://developer.safaricom.co.ke/)

## Local Development Setup

### 1. Clone and Setup Project

```bash
# Clone the repository
git clone <your-repository-url>
cd kenyan-real-estate-backend

# Install Go dependencies
go mod download

# Make scripts executable
chmod +x scripts/test.sh
```

### 2. Database Setup

```bash
# Create PostgreSQL database
createdb kenyan_real_estate

# Run database migrations
psql -d kenyan_real_estate -f migrations/001_initial_schema.sql
psql -d kenyan_real_estate -f migrations/002_kenyan_counties_data.sql
```

### 3. Environment Configuration

```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your configuration
nano .env
```

Required environment variables:
```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=kenyan_real_estate
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24

# M-Pesa (Get from Safaricom Developer Portal)
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_ENVIRONMENT=sandbox
MPESA_PASS_KEY=your_pass_key
MPESA_SHORT_CODE=your_short_code
```

### 4. Run the Application

```bash
# Run directly
go run cmd/server/main.go

# Or build and run
go build -o kenyan-real-estate-backend cmd/server/main.go
./kenyan-real-estate-backend
```

### 5. Test the Setup

```bash
# Run test script
./scripts/test.sh

# Or test manually
curl http://localhost:8080/health
```

## Production Deployment

### Option 1: Traditional Server Deployment

#### 1. Prepare the Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. Setup Database

```bash
# Switch to postgres user
sudo -u postgres psql

# Create database and user
CREATE DATABASE kenyan_real_estate;
CREATE USER real_estate_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE kenyan_real_estate TO real_estate_user;
\q

# Run migrations
psql -U real_estate_user -d kenyan_real_estate -f migrations/001_initial_schema.sql
psql -U real_estate_user -d kenyan_real_estate -f migrations/002_kenyan_counties_data.sql
```

#### 3. Deploy Application

```bash
# Clone repository
git clone <your-repository-url>
cd kenyan-real-estate-backend

# Build application
go build -o kenyan-real-estate-backend cmd/server/main.go

# Create production environment file
cp .env.example .env.production
# Edit with production values

# Create systemd service
sudo tee /etc/systemd/system/kenyan-real-estate.service > /dev/null <<EOF
[Unit]
Description=Kenyan Real Estate Backend
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/kenyan-real-estate-backend
ExecStart=/home/ubuntu/kenyan-real-estate-backend/kenyan-real-estate-backend
EnvironmentFile=/home/ubuntu/kenyan-real-estate-backend/.env.production
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Start and enable service
sudo systemctl daemon-reload
sudo systemctl enable kenyan-real-estate
sudo systemctl start kenyan-real-estate
```

#### 4. Setup Reverse Proxy (Nginx)

```bash
# Install Nginx
sudo apt install nginx -y

# Create Nginx configuration
sudo tee /etc/nginx/sites-available/kenyan-real-estate > /dev/null <<EOF
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

# Enable site
sudo ln -s /etc/nginx/sites-available/kenyan-real-estate /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Option 2: Docker Deployment

#### 1. Create Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/server/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
```

#### 2. Create Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
      - APP_ENV=production
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=kenyan_real_estate
      - DB_SSL_MODE=disable
      - JWT_SECRET=your-production-jwt-secret
      - MPESA_CONSUMER_KEY=your_consumer_key
      - MPESA_CONSUMER_SECRET=your_consumer_secret
      - MPESA_ENVIRONMENT=production
      - MPESA_PASS_KEY=your_pass_key
      - MPESA_SHORT_CODE=your_short_code
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=kenyan_real_estate
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    restart: unless-stopped

volumes:
  postgres_data:
```

#### 3. Deploy with Docker

```bash
# Build and run
docker-compose up -d

# Check logs
docker-compose logs -f app

# Stop
docker-compose down
```

### Option 3: Cloud Deployment (AWS/GCP/Azure)

#### AWS Deployment with ECS

1. **Create ECR Repository**
```bash
aws ecr create-repository --repository-name kenyan-real-estate-backend
```

2. **Build and Push Docker Image**
```bash
# Get login token
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

# Build and tag
docker build -t kenyan-real-estate-backend .
docker tag kenyan-real-estate-backend:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/kenyan-real-estate-backend:latest

# Push
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/kenyan-real-estate-backend:latest
```

3. **Create ECS Task Definition and Service**
```json
{
  "family": "kenyan-real-estate-backend",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::<account-id>:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "kenyan-real-estate-backend",
      "image": "<account-id>.dkr.ecr.us-east-1.amazonaws.com/kenyan-real-estate-backend:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {"name": "SERVER_HOST", "value": "0.0.0.0"},
        {"name": "SERVER_PORT", "value": "8080"},
        {"name": "APP_ENV", "value": "production"}
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/kenyan-real-estate-backend",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

## M-Pesa Configuration

### 1. Safaricom Developer Portal Setup

1. Visit [Safaricom Developer Portal](https://developer.safaricom.co.ke/)
2. Create an account and verify your email
3. Create a new app
4. Get your Consumer Key and Consumer Secret
5. Configure your app for STK Push

### 2. Environment Configuration

```env
# Sandbox (for testing)
MPESA_ENVIRONMENT=sandbox
MPESA_CONSUMER_KEY=your_sandbox_consumer_key
MPESA_CONSUMER_SECRET=your_sandbox_consumer_secret
MPESA_PASS_KEY=bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919
MPESA_SHORT_CODE=174379

# Production (for live)
MPESA_ENVIRONMENT=api
MPESA_CONSUMER_KEY=your_production_consumer_key
MPESA_CONSUMER_SECRET=your_production_consumer_secret
MPESA_PASS_KEY=your_production_pass_key
MPESA_SHORT_CODE=your_production_short_code
```

### 3. Callback URL Configuration

Set your callback URL in the Safaricom portal:
```
https://your-domain.com/api/v1/payments/mpesa/callback
```

## SSL/TLS Configuration

### Using Let's Encrypt with Certbot

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Get SSL certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

## Monitoring and Logging

### 1. Application Logs

```bash
# View systemd logs
sudo journalctl -u kenyan-real-estate -f

# View Docker logs
docker-compose logs -f app
```

### 2. Database Monitoring

```bash
# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-*.log

# Database connections
sudo -u postgres psql -c "SELECT * FROM pg_stat_activity;"
```

### 3. Health Checks

```bash
# Application health
curl https://your-domain.com/health

# Database health
sudo -u postgres psql -c "SELECT 1;"
```

## Backup and Recovery

### 1. Database Backup

```bash
# Create backup
pg_dump -U real_estate_user -h localhost kenyan_real_estate > backup_$(date +%Y%m%d_%H%M%S).sql

# Automated backup script
#!/bin/bash
BACKUP_DIR="/var/backups/kenyan-real-estate"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR
pg_dump -U real_estate_user -h localhost kenyan_real_estate > $BACKUP_DIR/backup_$DATE.sql
gzip $BACKUP_DIR/backup_$DATE.sql

# Keep only last 7 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
```

### 2. Application Backup

```bash
# Backup application files
tar -czf kenyan-real-estate-backup-$(date +%Y%m%d).tar.gz \
  /home/ubuntu/kenyan-real-estate-backend \
  /etc/systemd/system/kenyan-real-estate.service \
  /etc/nginx/sites-available/kenyan-real-estate
```

## Security Considerations

### 1. Environment Variables

- Never commit `.env` files to version control
- Use strong, unique passwords
- Rotate JWT secrets regularly
- Use environment-specific M-Pesa credentials

### 2. Database Security

```sql
-- Create read-only user for monitoring
CREATE USER monitoring WITH PASSWORD 'monitoring_password';
GRANT CONNECT ON DATABASE kenyan_real_estate TO monitoring;
GRANT USAGE ON SCHEMA public TO monitoring;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO monitoring;
```

### 3. Network Security

```bash
# Configure firewall
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Restrict PostgreSQL access
sudo nano /etc/postgresql/*/main/pg_hba.conf
# Change: host all all 0.0.0.0/0 md5
# To: host all all 127.0.0.1/32 md5
```

## Performance Optimization

### 1. Database Optimization

```sql
-- Add indexes for common queries
CREATE INDEX idx_properties_county_id ON properties(county_id);
CREATE INDEX idx_properties_rent_amount ON properties(rent_amount);
CREATE INDEX idx_properties_property_type ON properties(property_type);
CREATE INDEX idx_properties_is_available ON properties(is_available);
CREATE INDEX idx_payments_lease_id ON payments(lease_id);
CREATE INDEX idx_payments_status ON payments(status);
```

### 2. Application Optimization

```bash
# Build with optimizations
go build -ldflags="-s -w" -o kenyan-real-estate-backend cmd/server/main.go

# Use production Gin mode
export GIN_MODE=release
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check PostgreSQL status
   sudo systemctl status postgresql
   
   # Check database exists
   sudo -u postgres psql -l
   ```

2. **M-Pesa Integration Issues**
   ```bash
   # Check M-Pesa credentials
   curl -X GET "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials" \
     -H "Authorization: Basic $(echo -n 'consumer_key:consumer_secret' | base64)"
   ```

3. **Application Won't Start**
   ```bash
   # Check logs
   sudo journalctl -u kenyan-real-estate -n 50
   
   # Check port availability
   sudo netstat -tlnp | grep :8080
   ```

### Log Analysis

```bash
# Application errors
grep -i error /var/log/kenyan-real-estate.log

# M-Pesa transaction logs
grep -i mpesa /var/log/kenyan-real-estate.log

# Database query logs
sudo tail -f /var/log/postgresql/postgresql-*.log | grep ERROR
```

## Scaling Considerations

### Horizontal Scaling

1. **Load Balancer Configuration**
2. **Database Read Replicas**
3. **Redis for Session Management**
4. **CDN for Static Assets**

### Vertical Scaling

1. **Increase server resources**
2. **Database connection pooling**
3. **Application caching**

This deployment guide covers various scenarios from local development to production deployment. Choose the option that best fits your infrastructure and requirements.

