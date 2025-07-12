import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { propertyApi, locationApi } from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';
import ImageUpload from '../../components/ImageUpload';

const PropertyForm = () => {
  const { id } = useParams();
  const isEditing = Boolean(id);
  const navigate = useNavigate();
  const { user, isAgent } = useAuth();
  
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    rent_amount: '',
    property_type: 'house',
    bedrooms: '',
    bathrooms: '',
    square_meters: '',
    location_details: '',
    county_id: '',
    sub_county_id: '',
    amenities: {},
    is_furnished: false,
    is_available: true,
    parking_spaces: 0,
    deposit_amount: ''
  });
  
  const [counties, setCounties] = useState([]);
  const [subCounties, setSubCounties] = useState([]);
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({});
  const [featuresInput, setFeaturesInput] = useState('');
  const [images, setImages] = useState([]);

  const propertyTypes = [
    { value: 'house', label: 'House' },
    { value: 'apartment', label: 'Apartment' },
    { value: 'condo', label: 'Condo' },
    { value: 'townhouse', label: 'Townhouse' },
    { value: 'villa', label: 'Villa' },
    { value: 'land', label: 'Land' },
    { value: 'commercial', label: 'Commercial' }
  ];

  const statusOptions = [
    { value: 'active', label: 'Active' },
    { value: 'pending', label: 'Pending' },
    { value: 'sold', label: 'Sold' },
    { value: 'draft', label: 'Draft' }
  ];

  useEffect(() => {
    // Redirect if not an agent
    if (!isAgent()) {
      navigate('/');
      return;
    }

    fetchCounties();
    
    if (isEditing) {
      fetchProperty();
    }
  }, [id, isEditing, isAgent, navigate]);

  useEffect(() => {
    if (formData.county_id) {
      fetchSubCounties(formData.county_id);
    }
  }, [formData.county_id]);

  const fetchCounties = async () => {
    try {
      const response = await locationApi.getCounties();
      setCounties(response.data.counties || []);
    } catch (error) {
      console.error('Error fetching counties:', error);
    }
  };

  const fetchSubCounties = async (countyId) => {
    try {
      const response = await locationApi.getSubCounties(countyId);
      setSubCounties(response.data.sub_counties || []);
    } catch (error) {
      console.error('Error fetching sub-counties:', error);
      setSubCounties([]); // Reset sub-counties on error
    }
  };

  const fetchProperty = async () => {
    try {
      setLoading(true);
      const response = await propertyApi.getById(id);
      const property = response.data.property || response.data;
      
      setFormData({
        title: property.title || '',
        description: property.description || '',
        rent_amount: property.rent_amount || '',
        property_type: property.property_type || 'house',
        bedrooms: property.bedrooms || '',
        bathrooms: property.bathrooms || '',
        square_meters: property.square_meters || '',
        location_details: property.location_details || '',
        county_id: property.county_id || '',
        sub_county_id: property.sub_county_id || '',
        amenities: property.amenities || {},
        is_furnished: property.is_furnished || false,
        is_available: property.is_available !== undefined ? property.is_available : true,
        parking_spaces: property.parking_spaces || 0,
        deposit_amount: property.deposit_amount || ''
      });
      
      if (property.features && Array.isArray(property.features)) {
        setFeaturesInput(property.features.join(', '));
      }
      
      // Load existing images
      if (property.images && Array.isArray(property.images)) {
        const existingImages = property.images.map((img, index) => ({
          id: img.id || Date.now() + index,
          url: img.url || img,
          preview: img.url || img,
          isNew: false
        }));
        setImages(existingImages);
      }
    } catch (error) {
      console.error('Error fetching property:', error);
      setErrors({ general: 'Failed to load property data' });
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    
    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  const handleFeaturesChange = (e) => {
    const value = e.target.value;
    setFeaturesInput(value);
    
    // Convert comma-separated string to array
    const featuresArray = value.split(',').map(feature => feature.trim()).filter(feature => feature.length > 0);
    setFormData(prev => ({
      ...prev,
      features: featuresArray
    }));
  };

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.title.trim()) {
      newErrors.title = 'Title is required';
    }
    
    if (!formData.description.trim()) {
      newErrors.description = 'Description is required';
    }
    
    if (!formData.rent_amount || formData.rent_amount <= 0) {
      newErrors.rent_amount = 'Valid rent amount is required';
    }
    
    if (!formData.location_details.trim()) {
      newErrors.location_details = 'Location details are required';
    }
    
    if (!formData.bedrooms || formData.bedrooms < 0) {
      newErrors.bedrooms = 'Valid number of bedrooms is required';
    }
    
    if (!formData.bathrooms || formData.bathrooms < 0) {
      newErrors.bathrooms = 'Valid number of bathrooms is required';
    }
    
    if (!formData.square_meters || formData.square_meters <= 0) {
      newErrors.square_meters = 'Valid area is required';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const uploadImages = async (propertyId) => {
    const newImages = images.filter(img => img.isNew && img.file);
    
    if (newImages.length === 0) {
      return;
    }

    try {
      const formData = new FormData();
      newImages.forEach((image, index) => {
        formData.append('images', image.file);
      });

      await propertyApi.uploadImages(propertyId, formData);
    } catch (error) {
      console.error('Error uploading images:', error);
      throw new Error('Failed to upload images');
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }
    
    setLoading(true);
    
    try {
      const propertyData = {
        title: formData.title,
        description: formData.description,
        property_type: formData.property_type,
        rent_amount: parseFloat(formData.rent_amount),
        deposit_amount: formData.deposit_amount ? parseFloat(formData.deposit_amount) : null,
        bedrooms: parseInt(formData.bedrooms),
        bathrooms: parseInt(formData.bathrooms),
        square_meters: formData.square_meters ? parseFloat(formData.square_meters) : null,
        location_details: formData.location_details,
        county_id: formData.county_id ? parseInt(formData.county_id) : null,
        sub_county_id: formData.sub_county_id ? parseInt(formData.sub_county_id) : null,
        amenities: formData.amenities || {},
        is_furnished: formData.is_furnished,
        is_available: formData.is_available,
        parking_spaces: parseInt(formData.parking_spaces) || 0
      };
      
      let propertyId;
      
      if (isEditing) {
        const response = await propertyApi.update(id, propertyData);
        propertyId = id;
      } else {
        const response = await propertyApi.create(propertyData);
        propertyId = response.data.id || response.data.property?.id;
      }
      
      // Upload images if there are any new ones
      if (propertyId && images.some(img => img.isNew)) {
        await uploadImages(propertyId);
      }
      
      navigate('/agent/properties');
    } catch (error) {
      console.error('Error saving property:', error);
      setErrors({
        general: error.response?.data?.message || error.message || `Failed to ${isEditing ? 'update' : 'create'} property`
      });
    } finally {
      setLoading(false);
    }
  };

  if (loading && isEditing) {
    return (
      <div className="min-h-screen pt-20 bg-gray-50">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen pt-20 bg-gray-50">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">
            {isEditing ? 'Edit Property' : 'Add New Property'}
          </h1>
          <p className="text-gray-600 mt-2">
            {isEditing ? 'Update your property listing' : 'Create a new property listing'}
          </p>
        </div>

        {/* Form */}
        <div className="bg-white rounded-lg shadow-md p-6">
          {errors.general && (
            <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
              {errors.general}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Basic Information */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="md:col-span-2">
                <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
                  Property Title *
                </label>
                <input
                  type="text"
                  id="title"
                  name="title"
                  value={formData.title}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Beautiful 3BR house in Karen"
                />
                {errors.title && <p className="mt-1 text-sm text-red-600">{errors.title}</p>}
              </div>

              <div>
                <label htmlFor="property_type" className="block text-sm font-medium text-gray-700 mb-2">
                  Property Type *
                </label>
                <select
                  id="property_type"
                  name="property_type"
                  value={formData.property_type}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  {propertyTypes.map(type => (
                    <option key={type.value} value={type.value}>{type.label}</option>
                  ))}
                </select>
              </div>

              <div>
                <label htmlFor="is_available" className="block text-sm font-medium text-gray-700 mb-2">
                  Availability
                </label>
                <select
                  id="is_available"
                  name="is_available"
                  value={formData.is_available}
                  onChange={(e) => setFormData(prev => ({ ...prev, is_available: e.target.value === 'true' }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value={true}>Available</option>
                  <option value={false}>Not Available</option>
                </select>
              </div>

              <div>
                <label htmlFor="rent_amount" className="block text-sm font-medium text-gray-700 mb-2">
                  Rent Amount (USD) *
                </label>
                <input
                  type="number"
                  id="rent_amount"
                  name="rent_amount"
                  value={formData.rent_amount}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="50000"
                  min="0"
                />
                {errors.rent_amount && <p className="mt-1 text-sm text-red-600">{errors.rent_amount}</p>}
              </div>

              <div>
                <label htmlFor="square_meters" className="block text-sm font-medium text-gray-700 mb-2">
                  Area (sqm) *
                </label>
                <input
                  type="number"
                  id="square_meters"
                  name="square_meters"
                  value={formData.square_meters}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="250"
                  min="0"
                />
                {errors.square_meters && <p className="mt-1 text-sm text-red-600">{errors.square_meters}</p>}
              </div>

              <div>
                <label htmlFor="bedrooms" className="block text-sm font-medium text-gray-700 mb-2">
                  Bedrooms *
                </label>
                <input
                  type="number"
                  id="bedrooms"
                  name="bedrooms"
                  value={formData.bedrooms}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="3"
                  min="0"
                />
                {errors.bedrooms && <p className="mt-1 text-sm text-red-600">{errors.bedrooms}</p>}
              </div>

              <div>
                <label htmlFor="bathrooms" className="block text-sm font-medium text-gray-700 mb-2">
                  Bathrooms *
                </label>
                <input
                  type="number"
                  id="bathrooms"
                  name="bathrooms"
                  value={formData.bathrooms}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="2"
                  min="0"
                  step="0.5"
                />
                {errors.bathrooms && <p className="mt-1 text-sm text-red-600">{errors.bathrooms}</p>}
              </div>
            </div>

            {/* Location Information */}
            <div className="space-y-4">
              <h3 className="text-lg font-medium text-gray-900">Location Information</h3>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="md:col-span-2">
                  <label htmlFor="location_details" className="block text-sm font-medium text-gray-700 mb-2">
                    Location Details *
                  </label>
                  <input
                    type="text"
                    id="location_details"
                    name="location_details"
                    value={formData.location_details}
                    onChange={handleChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Karen, Nairobi - Near Karen Shopping Centre"
                  />
                  {errors.location_details && <p className="mt-1 text-sm text-red-600">{errors.location_details}</p>}
                </div>

                <div>
                  <label htmlFor="county_id" className="block text-sm font-medium text-gray-700 mb-2">
                    County
                  </label>
                  <select
                    id="county_id"
                    name="county_id"
                    value={formData.county_id}
                    onChange={handleChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="">Select County</option>
                    {counties.map(county => (
                      <option key={county.id} value={county.id}>{county.name}</option>
                    ))}
                  </select>
                </div>

                <div>
                  <label htmlFor="sub_county_id" className="block text-sm font-medium text-gray-700 mb-2">
                    Sub County
                  </label>
                  <select
                    id="sub_county_id"
                    name="sub_county_id"
                    value={formData.sub_county_id}
                    onChange={handleChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    disabled={!formData.county_id}
                  >
                    <option value="">Select Sub County</option>
                    {subCounties.map(subCounty => (
                      <option key={subCounty.id} value={subCounty.id}>{subCounty.name}</option>
                    ))}
                  </select>
                </div>
              </div>
            </div>

            {/* Description */}
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
                Description *
              </label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleChange}
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Describe the property features, amenities, and highlights..."
              />
              {errors.description && <p className="mt-1 text-sm text-red-600">{errors.description}</p>}
            </div>

            {/* Features */}
            <div>
              <label htmlFor="features" className="block text-sm font-medium text-gray-700 mb-2">
                Features (comma-separated)
              </label>
              <input
                type="text"
                id="features"
                value={featuresInput}
                onChange={handleFeaturesChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Swimming pool, Garden, Parking, Security, Modern kitchen"
              />
              <p className="mt-1 text-sm text-gray-500">
                Enter features separated by commas (e.g., Swimming pool, Garden, Parking)
              </p>
            </div>

            {/* Property Images */}
            <div className="space-y-4">
              <h3 className="text-lg font-medium text-gray-900">Property Images</h3>
              <ImageUpload 
                images={images} 
                onImagesChange={setImages} 
                maxImages={10}
              />
            </div>

            {/* Submit Buttons */}
            <div className="flex justify-end space-x-4 pt-6">
              <button
                type="button"
                onClick={() => navigate('/agent/properties')}
                className="px-6 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 font-medium transition-colors duration-200"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={loading}
                className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? 'Saving...' : isEditing ? 'Update Property' : 'Create Property'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default PropertyForm;
