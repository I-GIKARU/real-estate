-- Migration: 003_add_cloudinary_fields.sql
-- Add Cloudinary-specific fields to property_images table

-- Add new columns for Cloudinary integration
ALTER TABLE property_images 
ADD COLUMN secure_url TEXT,
ADD COLUMN public_id VARCHAR(255),
ADD COLUMN width INTEGER,
ADD COLUMN height INTEGER,
ADD COLUMN format VARCHAR(10),
ADD COLUMN bytes INTEGER;

-- Create index on public_id for efficient lookups
CREATE INDEX idx_property_images_public_id ON property_images(public_id);

-- Update existing records to have secure_url same as image_url (for backward compatibility)
UPDATE property_images SET secure_url = image_url WHERE secure_url IS NULL;

-- Add constraint to ensure public_id is unique when not null
CREATE UNIQUE INDEX idx_property_images_public_id_unique 
ON property_images(public_id) 
WHERE public_id IS NOT NULL;
