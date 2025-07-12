import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { propertyApi } from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';

const AgentProperties = () => {
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user, isAgent } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    // Redirect if not an agent
    if (!isAgent()) {
      navigate('/');
      return;
    }
    
    fetchProperties();
  }, [isAgent, navigate]);

  const fetchProperties = async () => {
    try {
      setLoading(true);
      const response = await propertyApi.getMyProperties();
      // Backend returns {properties: [...]} structure
      setProperties(response.data.properties || []);
    } catch (error) {
      console.error('Error fetching properties:', error);
      setError('Failed to load properties');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (propertyId) => {
    if (!window.confirm('Are you sure you want to delete this property?')) {
      return;
    }

    try {
      await propertyApi.delete(propertyId);
      setProperties(properties.filter(p => p.id !== propertyId));
    } catch (error) {
      console.error('Error deleting property:', error);
      setError('Failed to delete property');
    }
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(price);
  };

  if (loading) {
    return (
      <div className="min-h-screen pt-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen pt-20 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">My Properties</h1>
              <p className="text-gray-600 mt-2">Manage your property listings</p>
            </div>
            <Link
              to="/agent/properties/new"
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors duration-200 flex items-center gap-2"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v16m8-8H4" />
              </svg>
              Add New Property
            </Link>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {error}
          </div>
        )}

        {/* Properties Grid */}
        {properties.length === 0 ? (
          <div className="text-center py-12">
            <div className="max-w-md mx-auto">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
              <h3 className="mt-4 text-lg font-medium text-gray-900">No properties yet</h3>
              <p className="mt-2 text-gray-500">Get started by adding your first property listing.</p>
              <Link
                to="/agent/properties/new"
                className="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
              >
                Add Property
              </Link>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {properties.map((property) => (
              <div key={property.id} className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-200">
                {/* Property Image */}
                <div className="h-48 bg-gray-200 relative">
                  {property.images && property.images.length > 0 ? (
                    <img
                      src={property.images[0].secure_url || property.images[0].image_url}
                      alt={property.title}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400">
                      <svg className="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                      </svg>
                    </div>
                  )}
                  <div className="absolute top-2 right-2">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                      property.is_available ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                    }`}>
                      {property.is_available ? 'Available' : 'Unavailable'}
                    </span>
                  </div>
                </div>

                {/* Property Info */}
                <div className="p-4">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2 line-clamp-1">
                    {property.title}
                  </h3>
                  <p className="text-sm text-gray-600 mb-2 line-clamp-2">
                    {property.location_details || property.county?.name || 'Location not specified'}
                  </p>
                  <div className="flex justify-between items-center mb-4">
                    <span className="text-xl font-bold text-blue-600">
                      {formatPrice(property.rent_amount)}
                    </span>
                    <span className="text-sm text-gray-500 capitalize">
                      {property.property_type}
                    </span>
                  </div>

                  {/* Property Details */}
                  <div className="flex justify-between text-sm text-gray-600 mb-4">
                    <span>{property.bedrooms || 0} beds</span>
                    <span>{property.bathrooms || 0} baths</span>
                    <span>{property.square_meters || 0} sqm</span>
                  </div>

                  {/* Action Buttons */}
                  <div className="flex gap-2">
                    <Link
                      to={`/property/${property.id}`}
                      className="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 py-2 px-3 rounded text-center text-sm font-medium transition-colors duration-200"
                    >
                      View
                    </Link>
                    <Link
                      to={`/agent/properties/${property.id}/edit`}
                      className="flex-1 bg-blue-100 hover:bg-blue-200 text-blue-700 py-2 px-3 rounded text-center text-sm font-medium transition-colors duration-200"
                    >
                      Edit
                    </Link>
                    <button
                      onClick={() => handleDelete(property.id)}
                      className="flex-1 bg-red-100 hover:bg-red-200 text-red-700 py-2 px-3 rounded text-sm font-medium transition-colors duration-200"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default AgentProperties;
