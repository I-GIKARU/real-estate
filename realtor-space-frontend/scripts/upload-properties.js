const axios = require('axios');
const fs = require('fs');
const path = require('path');

// Configuration
const API_BASE_URL = 'https://realtor-space-backend.onrender.com/api';
const AGENT_EMAIL = 'ict622j1a@gmail.com';
const AGENT_PASSWORD = 'SecurePassword123!';

// Read properties data
const propertiesData = JSON.parse(fs.readFileSync(path.join(__dirname, '../public/properties.json'), 'utf8'));

// Helper function to convert local image paths to accessible URLs
function convertImagePaths(images) {
    return images.map(imagePath => {
        // Remove relative path prefixes and convert to absolute URLs
        const cleanPath = imagePath.replace(/^\.\.?\/?/, '').replace(/^\//, '');
        // Convert to accessible URL (assuming images are served from public folder)
        return `${API_BASE_URL.replace('/api', '')}/assets/images/${cleanPath.replace('assets/images/', '')}`;
    });
}

// Transform property data to match API format
function transformProperty(property) {
    return {
        title: property.name,
        description: property.description,
        property_type: property.category.toLowerCase().replace(/\s+/g, '_'),
        listing_type: "rent", // Assuming all are for rent
        rent_amount: property.price,
        deposit_amount: property.price * 0.5, // Assuming deposit is 50% of rent
        location: {
            address: property.location,
            city: property.location.split(',').pop().trim(),
            state: "Kenya",
            country: "Kenya"
        },
        bedrooms: property.bedrooms,
        bathrooms: property.bathrooms,
        area: property.area,
        amenities: property.amenities,
        images: property.images.some(img => img.startsWith('http')) 
            ? property.images 
            : convertImagePaths(property.images),
        virtual_tour_url: property.virtualTour || null,
        contact_info: {
            name: property.management.name,
            phone: property.management.contact,
            email: `${property.management.name.toLowerCase().replace(/\s+/g, '.')}@example.com`
        },
        availability_status: "available",
        furnished: property.amenities.some(amenity => 
            amenity.toLowerCase().includes('furnished')
        ),
        pets_allowed: false,
        utilities_included: property.amenities.some(amenity => 
            amenity.toLowerCase().includes('water') || 
            amenity.toLowerCase().includes('electricity')
        )
    };
}

// Register or login agent
async function authenticateAgent() {
    try {
        console.log('🔐 Attempting to authenticate agent...');
        
        // Try to login first
        try {
            const loginResponse = await axios.post(`${API_BASE_URL}/auth/login`, {
                email: AGENT_EMAIL,
                password: AGENT_PASSWORD
            });
            
            console.log('✅ Agent login successful!');
            return loginResponse.data.token;
        } catch (loginError) {
            console.log('📝 Login failed, attempting registration...');
            
            // If login fails, try to register
            const registerResponse = await axios.post(`${API_BASE_URL}/auth/register`, {
                first_name: 'Agent',
                last_name: 'User',
                email: AGENT_EMAIL,
                password: AGENT_PASSWORD,
                phone: '+254700000000',
                user_type: 'agent'
            });
            
            console.log('✅ Agent registration successful!');
            return registerResponse.data.token;
        }
    } catch (error) {
        console.error('❌ Authentication failed:', error.response?.data || error.message);
        throw error;
    }
}

// Upload a single property
async function uploadProperty(property, token) {
    try {
        const transformedProperty = transformProperty(property);
        
        console.log(`📤 Uploading property: ${property.name}`);
        
        const response = await axios.post(`${API_BASE_URL}/properties`, transformedProperty, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });
        
        console.log(`✅ Successfully uploaded: ${property.name}`);
        return response.data;
    } catch (error) {
        console.error(`❌ Failed to upload ${property.name}:`, error.response?.data || error.message);
        throw error;
    }
}

// Main execution function
async function main() {
    try {
        console.log('🚀 Starting property upload process...');
        console.log(`📊 Found ${propertiesData.properties.length} properties to upload`);
        
        // Authenticate agent
        const token = await authenticateAgent();
        
        // Upload each property
        const uploadResults = [];
        for (const property of propertiesData.properties) {
            try {
                const result = await uploadProperty(property, token);
                uploadResults.push({ success: true, property: property.name, data: result });
                
                // Add delay between uploads to avoid rate limiting
                await new Promise(resolve => setTimeout(resolve, 1000));
            } catch (error) {
                uploadResults.push({ success: false, property: property.name, error: error.message });
            }
        }
        
        // Summary
        const successful = uploadResults.filter(r => r.success).length;
        const failed = uploadResults.filter(r => !r.success).length;
        
        console.log('\n📋 Upload Summary:');
        console.log(`✅ Successful: ${successful}`);
        console.log(`❌ Failed: ${failed}`);
        
        if (failed > 0) {
            console.log('\n❌ Failed uploads:');
            uploadResults.filter(r => !r.success).forEach(result => {
                console.log(`- ${result.property}: ${result.error}`);
            });
        }
        
        console.log('\n🎉 Property upload process completed!');
        
    } catch (error) {
        console.error('💥 Fatal error:', error.message);
        process.exit(1);
    }
}

// Run the script
if (require.main === module) {
    main();
}

module.exports = { main, transformProperty, authenticateAgent };
