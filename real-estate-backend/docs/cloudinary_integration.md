# Cloudinary Integration for Property Image Management

This document describes the Cloudinary integration for handling property image uploads in the Kenyan Real Estate Backend.

## Overview

The system now supports direct file uploads for property images using Cloudinary as the storage and CDN provider. This replaces the previous URL-only approach with a more robust solution that includes:

- Direct file upload handling
- Automatic image optimization
- CDN delivery
- Image transformations
- Secure storage with metadata

## Configuration

### Environment Variables

Add the following to your `.env` file:

```env
# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
CLOUDINARY_FOLDER=real-estate-properties
```

### Obtaining Cloudinary Credentials

1. Sign up at [Cloudinary](https://cloudinary.com/)
2. Navigate to your Dashboard
3. Copy the Cloud Name, API Key, and API Secret
4. Set the folder name for organizing your property images

## Database Changes

Run the migration to add Cloudinary-specific fields:

```bash
psql -U postgres -d kenyan_real_estate -f migrations/003_add_cloudinary_fields.sql
```

This adds the following columns to `property_images`:
- `secure_url` - HTTPS URL from Cloudinary
- `public_id` - Cloudinary's unique identifier
- `width` - Image width in pixels
- `height` - Image height in pixels
- `format` - Image format (jpg, png, webp, etc.)
- `bytes` - File size in bytes

## API Usage

### Upload Property Image

**Endpoint**: `POST /api/v1/properties/{property_id}/images`

**Authentication**: Required (Landlord only)

**Content-Type**: `multipart/form-data`

**Form Fields**:
- `image` (file, required) - The image file to upload
- `caption` (string, optional) - Image caption
- `is_primary` (boolean, optional) - Whether this is the primary image
- `display_order` (integer, optional) - Display order for sorting

**Example using curl**:
```bash
curl -X POST \
  http://localhost:8080/api/v1/properties/{property_id}/images \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "image=@/path/to/your/image.jpg" \
  -F "caption=Beautiful living room" \
  -F "is_primary=true" \
  -F "display_order=1"
```

**Example using JavaScript**:
```javascript
const formData = new FormData();
formData.append('image', imageFile);
formData.append('caption', 'Beautiful living room');
formData.append('is_primary', 'true');
formData.append('display_order', '1');

fetch('/api/v1/properties/{property_id}/images', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + token
  },
  body: formData
})
.then(response => response.json())
.then(data => console.log(data));
```

**Response**:
```json
{
  "message": "Image added successfully",
  "image": {
    "id": "uuid-here",
    "property_id": "property-uuid",
    "image_url": "https://res.cloudinary.com/your-cloud/image/upload/v1234567890/real-estate-properties/property_uuid_filename_timestamp.jpg",
    "secure_url": "https://res.cloudinary.com/your-cloud/image/upload/v1234567890/real-estate-properties/property_uuid_filename_timestamp.jpg",
    "public_id": "real-estate-properties/property_uuid_filename_timestamp",
    "width": 1920,
    "height": 1080,
    "format": "jpg",
    "bytes": 245760,
    "caption": "Beautiful living room",
    "is_primary": true,
    "display_order": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Delete Property Image

**Endpoint**: `DELETE /api/v1/properties/{property_id}/images/{image_id}`

**Authentication**: Required (Landlord only)

**Example**:
```bash
curl -X DELETE \
  http://localhost:8080/api/v1/properties/{property_id}/images/{image_id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "message": "Image deleted successfully"
}
```

## File Validation

### Supported Formats
- JPEG (.jpg, .jpeg)
- PNG (.png)
- WebP (.webp)

### File Size Limits
- Maximum file size: 10MB
- Configurable via `MAX_FILE_SIZE` environment variable

### Automatic Optimizations
Cloudinary automatically applies:
- Quality optimization (`q_auto`)
- Format optimization (`f_auto`)
- Compression based on content

## Image Transformations

### Generating Thumbnails

The Cloudinary service provides methods for generating transformed images:

```go
// Generate a 300x200 thumbnail
thumbnailURL := cloudinaryService.GetImageThumbnail(publicID, 300, 200)
```

### Custom Transformations

You can generate custom transformation URLs:

```go
// Apply custom transformations
transformedURL := cloudinaryService.GenerateTransformationURL(publicID, "w_500,h_300,c_fill,q_80")
```

Common transformations:
- `w_500,h_300` - Resize to 500x300
- `c_fill` - Crop and fill the dimensions
- `c_fit` - Fit within dimensions
- `q_80` - Set quality to 80%
- `f_webp` - Convert to WebP format

## Error Handling

### Validation Errors
- File too large: `"file size exceeds maximum limit of 10MB"`
- Invalid format: `"file type image/gif not allowed"`
- Missing file: `"Image file is required"`

### Upload Errors
- Cloudinary service issues are logged and return generic error messages
- Database failures are handled gracefully
- Partial failures (e.g., Cloudinary upload succeeds but database fails) are handled

### Delete Errors
- If Cloudinary deletion fails, the database record is still removed
- This prevents orphaned database records
- Cloudinary cleanup can be handled separately if needed

## Security Features

### Access Control
- Only property owners can upload/delete images
- JWT authentication required
- Property ownership verification

### File Security
- Cloudinary provides secure URLs
- Public IDs are generated with UUIDs to prevent guessing
- File validation prevents malicious uploads

### Rate Limiting
Standard API rate limits apply to upload endpoints.

## Monitoring and Logging

### Cloudinary Dashboard
Monitor usage, transformations, and storage in your Cloudinary dashboard.

### Application Logs
Upload and deletion operations are logged for debugging and monitoring.

## Best Practices

### Frontend Implementation
1. Show upload progress to users
2. Validate files client-side before upload
3. Display thumbnail versions for better performance
4. Handle upload errors gracefully

### Performance
1. Use Cloudinary's automatic optimization features
2. Implement lazy loading for image galleries
3. Use appropriate image sizes for different contexts

### Error Recovery
1. Implement retry logic for failed uploads
2. Provide clear error messages to users
3. Log detailed error information for debugging

## Migration from URL-based System

If you have existing properties with image URLs:

1. Run the migration to add new columns
2. Existing `image_url` values remain functional
3. New uploads use the enhanced Cloudinary fields
4. Gradually migrate existing images if needed

The system is backward compatible with existing URL-based images while providing enhanced functionality for new uploads.
