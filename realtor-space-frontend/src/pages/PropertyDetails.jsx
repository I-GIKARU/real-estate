import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'

const PropertyDetails = () => {
  const { id } = useParams()
  const [property, setProperty] = useState(null)
  const [currentImageIndex, setCurrentImageIndex] = useState(0)

  useEffect(() => {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'https://real-estate-backend-840370620772.us-central1.run.app'
    fetch(`${apiBaseUrl}/api/v1/properties/${id}`)
        .then(response => response.json())
        .then(data => {
          if (data.property) {
            setProperty(data.property)
            document.title = `${data.property.title} | RealtorSpace`
            // Debug: Log the number of images received
            console.log('Property images count:', data.property.images?.length || 0)
            console.log('Property images:', data.property.images)
          }
        })
        .catch(error => {
          console.error('Error fetching property:', error)
        })
  }, [id])

  if (!property) {
    return <div className="py-20 text-center">Loading property details...</div>
  }

  const nextImage = () => {
    setCurrentImageIndex((prev) => (prev + 1) % (property.images?.length || 1))
  }

  const prevImage = () => {
    setCurrentImageIndex((prev) => (prev - 1 + (property.images?.length || 1)) % (property.images?.length || 1))
  }

  const selectImage = (index) => {
    setCurrentImageIndex(index)
  }

  return (
      <div className="pt-20">
        {/* Property Gallery */}
        <section className="max-w-6xl mx-auto px-5">
          {property.images && property.images.length > 0 ? (
            <>
              <div className="relative h-[60vh] rounded-xl overflow-hidden mb-5">
                <img
src={property.images[currentImageIndex]?.secure_url || property.images[currentImageIndex]?.image_url || property.image_url}
                    alt={property.title}
                    className="w-full h-full object-cover"
                />
                {property.images.length > 1 && (
                  <div className="absolute top-1/2 left-0 right-0 flex justify-between px-5 -translate-y-1/2">
                    <button
                        onClick={prevImage}
                        className="bg-white bg-opacity-70 w-12 h-12 rounded-full flex items-center justify-center hover:bg-white transition-colors"
                    >
                      <i className="fas fa-chevron-left"></i>
                    </button>
                    <button
                        onClick={nextImage}
                        className="bg-white bg-opacity-70 w-12 h-12 rounded-full flex items-center justify-center hover:bg-white transition-colors"
                    >
                      <i className="fas fa-chevron-right"></i>
                    </button>
                  </div>
                )}
              </div>

              {property.images.length > 1 && (
                <div className="flex gap-3 overflow-x-auto pb-3">
                  {property.images.map((image, index) => (
                      <div
                          key={index}
                          onClick={() => selectImage(index)}
                          className={`w-24 h-16 rounded-md overflow-hidden cursor-pointer flex-shrink-0 border-2 ${currentImageIndex === index ? 'border-blue-500' : 'border-transparent'}`}
                      >
                        <img
src={image.secure_url || image.image_url}
                            alt={`Property view ${index + 1}`}
                            className="w-full h-full object-cover"
                        />
                      </div>
                  ))}
                </div>
              )}
            </>
          ) : (
            <div className="h-[60vh] bg-gray-200 rounded-xl flex items-center justify-center mb-5">
              <div className="text-center text-gray-500">
                <i className="fas fa-image text-6xl mb-4"></i>
                <p>No images available for this property</p>
              </div>
            </div>
          )}
        </section>

        {/* Property Details */}
        <section className="max-w-6xl mx-auto px-5 py-10">
          <div className="mb-10">
            <a
                href={`https://wa.me/254757577018?text=Hi%20RealtorSpace,%20I'm%20interested%20in%20${encodeURIComponent(property.title)}`}
                className="bg-green-500 text-white px-6 py-3 rounded-md flex items-center justify-center gap-2 hover:bg-green-600 transition-colors mb-5"
            >
              <i className="fab fa-whatsapp"></i> WhatsApp Us
            </a>

            <h1 className="text-3xl font-bold mb-3">{property.title}</h1>

            <div className="flex items-center gap-5 mb-4">
              <span className="text-blue-600 text-xl font-bold">KES {property.rent_amount ? property.rent_amount.toLocaleString() : 'N/A'}</span>
              <span className="bg-green-100 px-4 py-1 rounded-full text-sm">{property.property_type}</span>
            </div>

            <div className="flex items-center gap-2 text-gray-600">
              <i className="fas fa-map-marker-alt"></i>
              <span>{property.county?.name || property.location}</span>
            </div>
          </div>

          <div className="space-y-10">
            {/* Description */}
            <div>
              <h2 className="text-2xl font-bold mb-4 pb-2 border-b border-gray-200">Description</h2>
              <p className="text-gray-700 leading-relaxed">{property.description}</p>
            </div>

            {/* Amenities */}
            <div>
              <h2 className="text-2xl font-bold mb-4 pb-2 border-b border-gray-200">Amenities</h2>
              <ul className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
                {property.amenities && typeof property.amenities === 'object' ? 
                  Object.entries(property.amenities)
                    .filter(([key, value]) => value === true)
                    .map(([amenity, _], index) => (
                      <li key={index} className="flex items-center gap-3">
                        <i className="fas fa-check text-blue-500"></i>
                        <span>{amenity.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}</span>
                      </li>
                    ))
                  : <p>No amenities available.</p>
                }
              </ul>
            </div>

            {/* Virtual Tour */}
            <div>
              <h2 className="text-2xl font-bold mb-4 pb-2 border-b border-gray-200">Virtual Tour</h2>
              {property.virtualTour ? (
                  <div className="relative pb-[56.25%] h-0 overflow-hidden rounded-xl">
                    <iframe
                        src={property.virtualTour}
                        className="absolute top-0 left-0 w-full h-full"
                        frameBorder="0"
                        allowFullScreen
                    ></iframe>
                  </div>
              ) : (
                  <p className="text-gray-600">Virtual tour not available for this property.</p>
              )}
            </div>

            {/* Property Management */}
            <div>
              <h2 className="text-2xl font-bold mb-4 pb-2 border-b border-gray-200">Property Management</h2>
              {property.management ? (
                <div className="flex gap-5 items-center bg-blue-100 p-5 rounded-xl">
                  <div>
                    <img
                        src={property.management.photo}
                        alt={property.management.name}
                        className="w-20 h-20 rounded-full object-cover"
                    />
                  </div>
                  <div>
                    <h3 className="font-bold">{property.management.name}</h3>
                    <p className="text-gray-600">{property.management.type}</p>
                    <p className="text-gray-600 mt-1">
                      <i className="fas fa-phone mr-2"></i>
                      {property.management.contact}
                    </p>
                  </div>
                </div>
              ) : (
                <p className="text-gray-600">Property management information not available.</p>
              )}
            </div>
          </div>
        </section>
      </div>
  )
}

export default PropertyDetails