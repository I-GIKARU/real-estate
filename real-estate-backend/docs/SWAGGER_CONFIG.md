# Swagger Configuration

## Environment Variables

### SWAGGER_HOST
Configures the host URL that Swagger UI will use to make API requests.

**Local Development:**
```env
SWAGGER_HOST=localhost:8080
```

**Production:**
```env
SWAGGER_HOST=real-estate-backend-840370620772.us-central1.run.app
```

**Docker/Cloud Run:**
The production host is automatically set in the Dockerfile for Cloud Run deployment.

## Configuration Files

### .env file
Add the following to your `.env` file:
```env
# Swagger Configuration
SWAGGER_HOST=localhost:8080
```

### .env.example file
The example file includes a template for the SWAGGER_HOST variable.

## How it Works

1. The application reads the `SWAGGER_HOST` environment variable
2. If not set, it falls back to the server configuration (`SERVER_HOST:SERVER_PORT`)
3. In production, the Dockerfile sets the correct Cloud Run URL
4. This ensures Swagger UI makes requests to the correct host in all environments

## Testing

After setting up the environment variable:

1. **Local Development:** Access Swagger at `http://localhost:8080/swagger/index.html`
2. **Production:** Access Swagger at `https://real-estate-backend-840370620772.us-central1.run.app/swagger/index.html`

The "Try it out" functionality should work correctly in both environments.
