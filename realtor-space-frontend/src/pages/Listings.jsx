import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import PropertyCard from '../components/PropertyCard'

const Listings = () => {
  const [properties, setProperties] = useState([])
  const [filteredProperties, setFilteredProperties] = useState([])
  const [filters, setFilters] = useState({
    propertyType: 'all',
    location: 'all',
    priceRange: 'all'
  })

  useEffect(() => {
    fetch('/data/properties.json')
        .then(response => response.json())
        .then(data => {
          setProperties(data.properties)
          setFilteredProperties(data.properties)
        })
  }, [])

  const handleFilterChange = (e) => {
    const { name, value } = e.target
    setFilters(prev => ({
      ...prev,
      [name]: value
    }))
  }

  const applyFilters = () => {
    let filtered = [...properties]

    if (filters.propertyType !== 'all') {
      filtered = filtered.filter(property => {
        if (filters.propertyType === 'single') return property.category === 'Single Room'
        if (filters.propertyType === 'one-bedroom') return property.category === 'One Bedroom'
        if (filters.propertyType === 'two-bedroom') return property.category === 'Two Bedroom'
        if (filters.propertyType === 'three-bedroom') return property.category.includes('Three Bedroom') || property.category.includes('+')
        if (filters.propertyType === 'apartment') return property.category === 'Apartment'
        if (filters.propertyType === 'house') return property.category === 'House'
        return true
      })
    }

    if (filters.location !== 'all') {
      filtered = filtered.filter(property =>
          property.location.toLowerCase().includes(filters.location.toLowerCase())
      )
    }

    if (filters.priceRange !== 'all') {
      const [min, max] = filters.priceRange.split('-').map(val =>
          val.endsWith('+') ? parseInt(val) : parseInt(val)
      )

      filtered = filtered.filter(property => {
        if (filters.priceRange.endsWith('+')) {
          return property.price >= min
        } else {
          return property.price >= min && property.price <= max
        }
      })
    }

    setFilteredProperties(filtered)
  }

  return (
      <>
        {/* Hero Section */}
        <section
            className="h-[60vh] bg-gradient-to-r from-black/50 to-black/50 bg-cover bg-center flex flex-col justify-center relative"
            style={{ backgroundImage: "url('https://images.unsplash.com/photo-1512917774080-9991f1c4c750?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80')" }}
        >
          <div className="max-w-3xl mx-auto text-center mt-20">
            <h1 className="text-4xl md:text-5xl font-bold mb-5">Premium Property Listings</h1>
            <p className="text-xl mb-8">Find your perfect home with our curated selection</p>
          </div>
        </section>

        {/* Filters */}
        <section className="py-10 px-5 bg-blue-100">
          <div className="max-w-7xl mx-auto">
            <div className="flex flex-wrap justify-center gap-5">
              <div className="min-w-[200px]">
                <label className="block mb-2 font-medium">Property Type</label>
                <select
                    name="propertyType"
                    value={filters.propertyType}
                    onChange={handleFilterChange}
                    className="w-full p-3 rounded border border-gray-300 bg-white"
                >
                  <option value="all">All Types</option>
                  <option value="single">Single Room</option>
                  <option value="one-bedroom">One Bedroom</option>
                  <option value="two-bedroom">Two Bedroom</option>
                  <option value="three-bedroom">Three Bedroom+</option>
                  <option value="apartment">Apartment</option>
                  <option value="house">House</option>
                </select>
              </div>

              <div className="min-w-[200px]">
                <label className="block mb-2 font-medium">Location</label>
                <select
                    name="location"
                    value={filters.location}
                    onChange={handleFilterChange}
                    className="w-full p-3 rounded border border-gray-300 bg-white"
                >
                  <option value="all">All Locations</option>
                  <option value="nyeri">Nyeri</option>
                  <option value="nanyuki">Nairobi</option>
                </select>
              </div>

              <div className="min-w-[200px]">
                <label className="block mb-2 font-medium">Price Range</label>
                <select
                    name="priceRange"
                    value={filters.priceRange}
                    onChange={handleFilterChange}
                    className="w-full p-3 rounded border border-gray-300 bg-white"
                >
                  <option value="all">All Prices</option>
                  <option value="0-10000">Under KES 10,000</option>
                  <option value="10000-30000">KES 10,000 - 30,000</option>
                  <option value="30000-50000">KES 30,000 - 50,000</option>
                  <option value="50000-100000">KES 50,000 - 100,000</option>
                  <option value="100000+">KES 100,000+</option>
                </select>
              </div>

              <button
                  onClick={applyFilters}
                  className="self-end bg-blue-600 text-white px-6 py-3 rounded-full font-semibold hover:bg-blue-700 transition-colors"
              >
                Apply Filters
              </button>
            </div>
          </div>
        </section>

        {/* Listings */}
        <section className="py-16 px-5">
          <div className="max-w-7xl mx-auto">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {filteredProperties.length > 0 ? (
                  filteredProperties.map(property => (
                      <PropertyCard key={property.id} property={property} />
                  ))
              ) : (
                  <div className="col-span-full text-center py-10">
                    <p className="text-gray-600">No properties match your filters. Try adjusting your search criteria.</p>
                  </div>
              )}
            </div>
          </div>
        </section>
      </>
  )
}

export default Listings