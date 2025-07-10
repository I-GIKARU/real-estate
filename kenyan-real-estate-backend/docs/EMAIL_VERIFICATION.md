# Email Verification System

This document explains how the email verification system works in the Kenyan Real Estate API.

## Overview

The email verification system ensures that users have access to the email address they registered with. This is important for:
- Account security
- Password recovery
- Important notifications
- Building trust between landlords and tenants

## How It Works

### 1. User Registration
When a user registers:
1. User account is created with `is_verified = false`
2. User receives JWT token (can use the API)
3. Response includes `email_verification_required: true`

### 2. Email Verification Process
1. User calls `/send-verification-email` (authenticated endpoint)
2. System generates secure token and sends email
3. User clicks link in email or uses token via `/verify-email`
4. User's `is_verified` status is updated to `true`
5. Welcome email is sent

## API Endpoints

### Send Verification Email
**POST** `/api/v1/send-verification-email`
- **Auth**: Required (Bearer token)
- **Description**: Sends verification email to authenticated user
- **Rate Limit**: 1 email per 5 minutes
- **Response**: Success message

```json
{
  "message": "Verification email sent successfully"
}
```

### Verify Email
**POST** `/api/v1/verify-email`
- **Auth**: Not required (public endpoint)
- **Body**: 
```json
{
  "token": "verification_token_from_email"
}
```
- **Response**: Success message and updated user info

```json
{
  "message": "Email verified successfully",
  "user": {
    "id": "user-id",
    "email": "user@example.com",
    "is_verified": true,
    ...
  }
}
```

### Check Verification Status
**GET** `/api/v1/verification-status`
- **Auth**: Required (Bearer token)
- **Response**: Current verification status

```json
{
  "is_verified": true,
  "pending_verification": false,
  "can_resend": false
}
```

## Email Configuration

Add these environment variables to your `.env` file:

```env
# Email Configuration
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your_email@gmail.com
EMAIL_PASSWORD=your_app_password
EMAIL_FROM=noreply@kenyanrealestate.com
EMAIL_SUPPORT=support@kenyanrealestate.com
BASE_URL=http://localhost:3000
```

### Using Gmail SMTP

1. Enable 2-factor authentication on your Gmail account
2. Generate an App Password:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate password for "Mail"
   - Use this password in `EMAIL_PASSWORD`

### Using Other Email Providers

Common SMTP settings:
- **Outlook/Hotmail**: `smtp-mail.outlook.com:587`
- **Yahoo**: `smtp.mail.yahoo.com:587`
- **SendGrid**: `smtp.sendgrid.net:587`
- **Mailgun**: `smtp.mailgun.org:587`

## Security Features

### Token Security
- Tokens are cryptographically secure (256-bit)
- Tokens expire after 24 hours
- Tokens are single-use only

### Rate Limiting
- Max 1 verification email per 5 minutes per user
- Prevents spam and abuse

### Email Templates
- Professional HTML emails
- Plain text fallback
- Clear instructions and branding

## Database Schema

The `email_verifications` table stores:
- `id`: UUID primary key
- `user_id`: Foreign key to users table
- `token`: Unique verification token
- `expires_at`: Token expiration time
- `is_used`: Whether token has been used
- `created_at`, `updated_at`: Timestamps

## Frontend Integration

### Registration Flow
1. User submits registration form
2. Show success message: "Account created! Please check your email."
3. Redirect to email verification page
4. Show verification status and resend option

### Verification Page
1. Check verification status on page load
2. Show appropriate UI based on status:
   - Already verified: Show success
   - Pending: Show "Check your email" with resend button
   - Can resend: Enable resend button
   - Rate limited: Show countdown timer

### Example Frontend Code
```javascript
// Check verification status
const checkVerificationStatus = async () => {
  const response = await fetch('/api/v1/verification-status', {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  const data = await response.json();
  
  if (data.is_verified) {
    showVerifiedUI();
  } else if (data.pending_verification && !data.can_resend) {
    showPendingUI(data.retry_after_seconds);
  } else {
    showResendUI();
  }
};

// Send verification email
const sendVerificationEmail = async () => {
  const response = await fetch('/api/v1/send-verification-email', {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` }
  });
  
  if (response.ok) {
    showEmailSentMessage();
  }
};
```

## Testing

### Development Testing
1. Use a service like [Mailtrap](https://mailtrap.io/) or [MailHog](https://github.com/mailhog/MailHog)
2. Set up local SMTP server for testing
3. Use real email for end-to-end testing

### Manual Testing
1. Register a new user
2. Check that verification email is sent
3. Verify email using token
4. Confirm user status is updated
5. Test rate limiting by requesting multiple emails

## Troubleshooting

### Common Issues

**Emails not sending:**
- Check SMTP credentials
- Verify firewall settings
- Check email provider limits

**Token validation errors:**
- Verify token hasn't expired
- Check token hasn't been used
- Ensure proper URL encoding

**Rate limiting issues:**
- Wait 5 minutes between requests
- Check server time settings
- Verify database timestamps

### Logs
Check application logs for:
- Email sending errors
- Database connection issues
- Token generation problems

## Production Considerations

1. **Email Provider**: Use dedicated email service (SendGrid, Mailgun, SES)
2. **Domain**: Set up proper SPF, DKIM, DMARC records
3. **Monitoring**: Track email delivery rates
4. **Cleanup**: Regularly clean expired verification records
5. **Scaling**: Consider queue-based email sending for high volume

## Future Enhancements

- Email template customization
- Multi-language support
- SMS verification option
- Admin verification dashboard
- Bulk verification management
