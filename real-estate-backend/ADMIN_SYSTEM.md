# Admin System Documentation

## Overview

The admin system provides administrative control over the real estate platform, primarily focused on agent approval workflows. Only one admin user exists in the system to maintain security and control.

## Admin Features

### 1. Agent Approval
- View all agents waiting for approval
- View all registered agents (approved and pending)
- Approve agents to allow them to create and manage properties
- Real-time statistics dashboard

### 2. Dashboard Features
- **Pending Approvals Count**: Shows agents waiting for approval
- **Approved Agents Count**: Shows total approved agents
- **Total Agents Count**: Shows all registered agents
- **Detailed Agent Information**: Name, email, phone, verification status, registration date
- **Action Buttons**: Quick approve functionality

## Access Information

### Default Admin Credentials
- **Email**: `admin@realestate.com`
- **Password**: `Admin123!`

### Admin Routes
- **Login**: `/admin/login`
- **Dashboard**: `/admin/dashboard`

### API Endpoints

#### Authentication
```bash
POST /api/v1/login
Content-Type: application/json

{
  "email": "admin@realestate.com",
  "password": "Admin123!"
}
```

#### Admin Endpoints (Require Admin Token)
```bash
# Get pending agents
GET /api/v1/admin/pending-agents
Authorization: Bearer <admin_token>

# Get all agents
GET /api/v1/admin/agents
Authorization: Bearer <admin_token>

# Approve an agent
POST /api/v1/admin/approve-agent/{agentId}
Authorization: Bearer <admin_token>
```

## Database Schema

The admin user is stored in the `users` table with:
- `user_type = 'admin'`
- `is_verified = true`
- `is_approved = true`
- `is_active = true`

## Security Features

1. **Single Admin Policy**: Only one admin user exists in the system
2. **JWT Authentication**: Secure token-based authentication
3. **Role-based Access**: Admin-only endpoints protected by middleware
4. **Secure Password**: Bcrypt hashed passwords
5. **Session Management**: Token-based sessions with expiry

## Agent Approval Workflow

1. **Agent Registration**: Agents register through the normal registration flow
2. **Email Verification**: Agents must verify their email addresses
3. **Pending State**: Verified agents appear in the admin pending list
4. **Admin Review**: Admin reviews agent information in the dashboard
5. **Approval**: Admin clicks "Approve" to activate the agent
6. **Active State**: Approved agents can create and manage properties

## Frontend Integration

### Admin Login Page
- Clean, professional login interface
- Secure credential validation
- Error handling and feedback
- Automatic redirect to dashboard

### Admin Dashboard
- Modern, responsive design
- Real-time data loading
- Interactive tables with agent details
- Quick approval actions
- Statistics cards
- Logout functionality

## Testing

Use the provided test script to verify admin functionality:

```bash
./test_admin_login.sh
```

This script tests:
- Admin login authentication
- Pending agents endpoint
- All agents endpoint
- Token validation

## Database Migration

The admin user is created via migration `006_create_admin_user.sql`:

```sql
INSERT INTO users (
    email, password_hash, first_name, last_name,
    phone_number, user_type, is_verified, 
    is_approved, is_active
) VALUES (
    'admin@realestate.com',
    '<bcrypt_hash>',
    'System', 'Administrator',
    '+254700000000', 'admin',
    true, true, true
);
```

## Production Considerations

1. **Change Default Password**: Update the admin password in production
2. **Secure Access**: Use HTTPS for all admin operations
3. **Monitoring**: Log all admin actions for audit trails
4. **Backup**: Ensure admin credentials are securely backed up
5. **Access Control**: Limit admin access to authorized personnel only

## Troubleshooting

### Common Issues

1. **Login Failed**: Verify credentials and check password hash
2. **Token Expired**: Re-login to get a new token
3. **Unauthorized Access**: Ensure user has admin role
4. **Network Errors**: Check API endpoint connectivity

### Logs

Check application logs for:
- Authentication failures
- Authorization errors
- Database connection issues
- API request/response details

---

For technical support or questions about the admin system, refer to the main project documentation or contact the development team.
