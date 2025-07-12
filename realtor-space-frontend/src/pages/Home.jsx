import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

const Home = () => {
  const [featuredProperties, setFeaturedProperties] = useState([])

  useEffect(() => {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
    fetch(`${apiBaseUrl}/properties`)
        .then(response => response.json())
        .then(data => {
          setFeaturedProperties(data.properties.slice(0, 3))
        })
        .catch(error => {
          console.error('Error fetching properties:', error)
        })
  }, [])

  return (
      <>
        {/* Hero Section */}
        <section
            className="h-[60vh] bg-gradient-to-r from-black/50 to-black/50 bg-cover bg-center flex flex-col justify-center relative"
            style={{ backgroundImage: "url('https://images.unsplash.com/photo-1560448204-e02f11c3d0e2?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80')" }}
        >
          <div className="max-w-3xl mx-auto text-center mt-20 animate-slide-up">
            <h1 className="text-4xl md:text-5xl font-bold mb-5">Discover Your Dream Property</h1>
            <p className="text-xl mb-8">Find it. Love it. Live it.</p>
            <Link
                to="/listings"
                className="bg-blue-600 text-white px-8 py-3 rounded-full font-semibold uppercase tracking-wider hover:bg-blue-700 transition-colors shadow-lg animate-pulse"
            >
              View Properties
            </Link>
          </div>
        </section>

        {/* Featured Properties */}
        <section className="py-20 px-5 bg-blue-100">
          <div className="max-w-7xl mx-auto">
            <div className="text-center mb-12">
              <h2 className="text-3xl font-bold mb-3">Featured Properties</h2>
              <p className="text-gray-700">Curated selection of premium homes</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {featuredProperties.map(property => (
                  <PropertyCard key={property.id} property={property} />
              ))}
            </div>
          </div>
        </section>

        {/* About Section */}
        <section id="about" className="py-20 px-5 bg-white">
          <div className="max-w-7xl mx-auto flex flex-col md:flex-row gap-12">
            <div className="md:w-1/2">
              <h2 className="text-3xl font-bold mb-5">About Realtor's Space</h2>
              <p className="text-lg mb-8 leading-relaxed">
                Based in Nyeri, Kenya, Realtor's Space brings you the finest property listings with a premium user
                experience. Our platform combines luxury real estate with cutting-edge technology to help you find your
                perfect home.
              </p>

              <div className="flex flex-col md:flex-row gap-8 mt-10">
                <div className="text-center">
                  <h3 className="text-2xl text-blue-600 mb-1">40+</h3>
                  <p className="text-gray-600">Properties Listed</p>
                </div>
                <div className="text-center">
                  <h3 className="text-2xl text-blue-600 mb-1">98%</h3>
                  <p className="text-gray-600">Client Satisfaction</p>
                </div>
                <div className="text-center">
                  <h3 className="text-2xl text-blue-600 mb-1">2</h3>
                  <p className="text-gray-600">Years Experience</p>
                </div>
              </div>
            </div>

            <div className="md:w-1/2">
              <img
                  src="https://images.unsplash.com/photo-1560518883-ce09059eeffa?ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60"
                  alt="Luxury building"
                  className="w-full rounded-lg shadow-xl"
              />
            </div>
          </div>
        </section>

        {/* Services Section */}
        <section id="services" className="py-20 px-5 bg-blue-100">
          <div className="max-w-7xl mx-auto">
            <div className="text-center mb-12">
              <h2 className="text-3xl font-bold mb-3">Our Premium Services</h2>
              <p className="text-gray-700">Beyond listings - we provide complete property solutions</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              <ServiceCard
                  icon="fas fa-search-dollar"
                  title="Custom Property Search"
                  description="Can't find what you're looking for? Our expert agents will conduct a personalized property search based on your exact requirements. We'll handpick options that match your budget, location, and amenities needs."
                  price="From KES 1,500"
              />
              <ServiceCard
                  icon="fas fa-truck-moving"
                  title="Professional Moving Services"
                  description="Moving to your new property? Our trusted moving partners offer stress-free relocation services. From packing to transportation and setup, we handle everything with care and professionalism."
              />
              <ServiceCard
                  icon="fas fa-house-user"
                  title="Guided Property Viewing"
                  description="Want to see a property in person? Our agents will arrange and accompany you on professional viewings. We provide insights about the neighborhood, answer questions, and help evaluate the property's potential."
                  price="From KES 500 per viewing"
              />
            </div>
          </div>
        </section>

        {/* Testimonials */}
        <section className="py-20 px-5 bg-green-100 text-center">
          <div className="max-w-7xl mx-auto">
            <h2 className="text-3xl font-bold mb-12">What Our Clients Say</h2>

            <div className="flex flex-col md:flex-row justify-center gap-8">
              <TestimonialCard
                  quote="Realtor's Space made finding my dream home effortless. The virtual tours were incredible!"
                  name="Sarah W."
                  role="Homeowner"
              />
              <TestimonialCard
                  quote="As a landlord, I've never had such quality tenants. The platform attracts serious buyers."
                  name="James K."
                  role="Property Owner"
              />
            </div>
          </div>
        </section>

        {/* WhatsApp Float Button */}
        <a
            href="https://wa.me/254757577018"
            className="floating-whatsapp"
            target="_blank"
            rel="noopener noreferrer"
        >
          <i className="fab fa-whatsapp"></i>
        </a>
      </>
  )
}

const PropertyCard = ({ property }) => {
  return (
      <div className="bg-white rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow hover-grow">
        <div className="h-48 overflow-hidden">
          <img
              src={property.images?.[0]?.secure_url || property.images?.[0]?.image_url || '/placeholder-property.jpg'}
              alt={property.title}
              className="w-full h-full object-cover transition-transform duration-500 hover:scale-110"
          />
        </div>
        <div className="p-6">
          <h3 className="text-xl font-semibold mb-2">{property.title}</h3>
          <div className="flex items-center gap-2 text-gray-600 mb-3">
            <i className="fas fa-map-marker-alt"></i>
            <span>{property.county?.name || property.location_details}</span>
          </div>
          <div className="text-blue-600 text-lg font-bold mb-4">KES {property.rent_amount?.toLocaleString()}</div>
          <div className="flex justify-between mb-4">
          <span className="flex items-center gap-1">
            <i className="fas fa-bed"></i> {property.bedrooms} Bed
          </span>
            <span className="flex items-center gap-1">
            <i className="fas fa-bath"></i> {property.bathrooms} Bath
          </span>
            <span className="flex items-center gap-1">
            <i className="fas fa-ruler-combined"></i> {property.square_meters} m²
          </span>
          </div>
          <Link
              to={`/property/${property.id}`}
              className="block w-full bg-blue-600 text-white text-center py-2 rounded-full font-semibold hover:bg-blue-700 transition-colors"
          >
            View Details
          </Link>
        </div>
      </div>
  )
}

const ServiceCard = ({ icon, title, description, price }) => {
  return (
      <div className="bg-white rounded-xl p-8 shadow-lg hover:shadow-xl transition-all hover:translate-y-[-10px] text-center">
        <div className="text-blue-600 text-4xl mb-5">
          <i className={icon}></i>
        </div>
        <h3 className="text-xl font-bold mb-4">{title}</h3>
        <p className="text-gray-600 mb-6">{description}</p>
        {price && <p className="font-bold text-blue-600 my-5">{price}</p>}
        <a
            href="https://wa.me/254757577018?text=Hi%20RealtorSpace,%20I'm%20interested%20in%20your%20service"
            className="bg-green-500 text-white px-6 py-2 rounded-md flex items-center justify-center gap-2 hover:bg-green-600 transition-colors"
        >
          <i className="fab fa-whatsapp"></i> WhatsApp Us
        </a>
      </div>
  )
}

const TestimonialCard = ({ quote, name, role }) => {
  return (
      <div className="bg-white rounded-xl p-8 max-w-md shadow-md text-left relative">
        <div className="italic mb-6 relative">
          <span className="absolute -top-6 -left-3 text-blue-600 opacity-20 text-5xl">"</span>
          {quote}
        </div>
        <div className="flex items-center">
          <div>
            <h4 className="font-bold">{name}</h4>
            <p className="text-gray-600 text-sm">{role}</p>
          </div>
        </div>
      </div>
  )
}

export default Home