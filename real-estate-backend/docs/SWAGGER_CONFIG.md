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

**Google Cloud Run:**
Set the environment variable in your Google Cloud Run service configuration.

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
3. For Google Cloud Run, set the environment variable in the service configuration
4. This ensures Swagger UI makes requests to the correct host in all environments

## Google Cloud Run Configuration

To set the environment variable in Google Cloud Run:

1. Go to Google Cloud Console → Cloud Run
2. Select your service
3. Click "Edit & Deploy New Revision"
4. Go to "Variables & Secrets" tab
5. Add environment variable:
   - Name: `SWAGGER_HOST`
   - Value: `real-estate-backend-840370620772.us-central1.run.app`
6. Deploy the revision

## Testing

After setting up the environment variable:

1. **Local Development:** Access Swagger at `http://localhost:8080/swagger/index.html`


The "Try it out" functionality should work correctly in both environments.
