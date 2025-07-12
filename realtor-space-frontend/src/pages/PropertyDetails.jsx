import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'

const PropertyDetails = () => {
  const { id } = useParams()
  const [property, setProperty] = useState(null)
  const [currentImageIndex, setCurrentImageIndex] = useState(0)

  useEffect(() => {
    fetch('/data/properties.json')
        .then(response => response.json())
        .then(data => {
          const foundProperty = data.properties.find(p => p.id == id)
          if (foundProperty) {
            setProperty(foundProperty)
            document.title = `${foundProperty.name} | RealtorSpace`
          }
        })
  }, [id])

  if (!property) {
    return <div className="py-20 text-center">Loading property details...</div>
  }

  const nextImage = () => {
    setCurrentImageIndex((prev) => (prev + 1) % property.images.length)
  }

  const prevImage = () => {
    setCurrentImageIndex((prev) => (prev - 1 + property.images.length) % property.images.length)
  }

  const selectImage = (index) => {
    setCurrentImageIndex(index)
  }

  return (
      <div className="pt-20">
        {/* Property Gallery */}
        <section className="max-w-6xl mx-auto px-5">
          <div className="relative h-[60vh] rounded-xl overflow-hidden mb-5">
            <img
                src={property.images[currentImageIndex]}
                alt={property.name}
                className="w-full h-full object-cover"
            />
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
          </div>

          <div className="flex gap-3 overflow-x-auto pb-3">
            {property.images.map((image, index) => (
                <div
                    key={index}
                    onClick={() => selectImage(index)}
                    className={`w-24 h-16 rounded-md overflow-hidden cursor-pointer flex-shrink-0 border-2 ${currentImageIndex === index ? 'border-blue-500' : 'border-transparent'}`}
                >
                  <img
                      src={image}
                      alt={`Property view ${index + 1}`}
                      className="w-full h-full object-cover"
                  />
                </div>
            ))}
          </div>
        </section>

        {/* Property Details */}
        <section className="max-w-6xl mx-auto px-5 py-10">
          <div className="mb-10">
            <a
                href={`https://wa.me/254757577018?text=Hi%20RealtorSpace,%20I'm%20interested%20in%20${encodeURIComponent(property.name)}`}
                className="bg-green-500 text-white px-6 py-3 rounded-md flex items-center justify-center gap-2 hover:bg-green-600 transition-colors mb-5"
            >
              <i className="fab fa-whatsapp"></i> WhatsApp Us
            </a>

            <h1 className="text-3xl font-bold mb-3">{property.name}</h1>

            <div className="flex items-center gap-5 mb-4">
              <span className="text-blue-600 text-xl font-bold">KES {property.price.toLocaleString()}</span>
              <span className="bg-green-100 px-4 py-1 rounded-full text-sm">{property.category}</span>
            </div>

            <div className="flex items-center gap-2 text-gray-600">
              <i className="fas fa-map-marker-alt"></i>
              <span>{property.location}</span>
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
                {property.amenities.map((amenity, index) => (
                    <li key={index} className="flex items-center gap-3">
                      <i className="fas fa-check text-blue-500"></i>
                      <span>{amenity}</span>
                    </li>
                ))}
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
            </div>
          </div>
        </section>
      </div>
  )
}

export default PropertyDetails