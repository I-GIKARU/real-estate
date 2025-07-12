import { Link } from 'react-router-dom';

const PropertyCard = ({ property, className = "" }) => {
  const formatPrice = (price) => {
    return new Intl.NumberFormat('en-KE', {
      style: 'currency',
      currency: 'KES',
      minimumFractionDigits: 0,
    }).format(price);
  };

  return (
    <div className={`card-modern hover-lift group ${className}`}>
      {/* Property Image */}
      <div className="relative h-64 overflow-hidden">
        <img
          src={property.images?.[0]?.secure_url || property.images?.[0]?.image_url || '/placeholder-property.jpg'}
          alt={property.title}
          className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/50 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
        
        {/* Category Badge */}
        <div className="absolute top-4 right-4">
          <span className="glass text-white px-3 py-1 rounded-full text-sm font-semibold backdrop-blur-md">
            {property.property_type}
          </span>
        </div>
        
        {/* Heart Icon */}
        <button className="absolute top-4 left-4 w-10 h-10 rounded-full glass backdrop-blur-md flex items-center justify-center text-white hover:text-red-400 transition-colors duration-200">
          <i className="far fa-heart"></i>
        </button>

        {/* Floating Info */}
        <div className="absolute bottom-4 left-4 right-4 transform translate-y-8 group-hover:translate-y-0 opacity-0 group-hover:opacity-100 transition-all duration-300">
          <div className="glass backdrop-blur-md rounded-lg p-3">
            <div className="flex items-center text-white text-sm">
              <i className="fas fa-map-marker-alt mr-2"></i>
              <span className="truncate">{property.county?.name || property.location_details}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Property Info */}
      <div className="p-6">
        <h3 className="text-xl font-bold text-gray-900 mb-3 line-clamp-2 group-hover:text-purple-600 transition-colors duration-200">
          {property.title}
        </h3>

        <div className="mb-4">
          <div className="text-3xl font-bold gradient-text mb-1">
            {formatPrice(property.rent_amount)}
          </div>
          <span className="text-sm text-gray-500">/month</span>
        </div>

        {/* Property Features */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center space-x-4">
            <div className="flex items-center text-gray-600">
              <div className="w-8 h-8 rounded-full bg-purple-100 flex items-center justify-center mr-2">
                <i className="fas fa-bed text-purple-600 text-sm"></i>
              </div>
              <span className="text-sm font-medium">{property.bedrooms}</span>
            </div>
            <div className="flex items-center text-gray-600">
              <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center mr-2">
                <i className="fas fa-bath text-blue-600 text-sm"></i>
              </div>
              <span className="text-sm font-medium">{property.bathrooms}</span>
            </div>
            <div className="flex items-center text-gray-600">
              <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center mr-2">
                <i className="fas fa-ruler-combined text-green-600 text-sm"></i>
              </div>
              <span className="text-sm font-medium">{property.square_meters}m²</span>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex space-x-3">
          <Link
            to={`/property/${property.id}`}
            className="flex-1 btn-modern text-center hover-glow"
          >
            View Details
          </Link>
          <button className="px-4 py-3 rounded-xl border-2 border-purple-200 text-purple-600 hover:bg-purple-50 transition-colors duration-200 font-medium">
            <i className="fas fa-share-alt"></i>
          </button>
        </div>
      </div>
    </div>
  );
};

export default PropertyCard;
