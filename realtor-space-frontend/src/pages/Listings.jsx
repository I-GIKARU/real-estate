import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import PropertyCard from '../components/PropertyCard'

const Listings = () => {
  const [properties, setProperties] = useState([])
  const [filteredProperties, setFilteredProperties] = useState([])
  const [counties, setCounties] = useState([])
  const [subCounties, setSubCounties] = useState([])
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState({
    propertyType: 'all',
    county: 'all',
    subCounty: 'all',
    priceRange: 'all'
  })

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        
        // Fetch properties and counties in parallel
        const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
        const [propertiesResponse, countiesResponse] = await Promise.all([
          fetch(`${apiBaseUrl}/properties`),
          fetch(`${apiBaseUrl}/counties`)
        ])
        
        const propertiesData = await propertiesResponse.json()
        const countiesData = await countiesResponse.json()
        
        setProperties(propertiesData.properties)
        setFilteredProperties(propertiesData.properties)
        setCounties(countiesData.counties)
      } catch (error) {
        console.error('Error fetching data:', error)
      } finally {
        setLoading(false)
      }
    }
    
    fetchData()
  }, [])

  const handleFilterChange = async (e) => {
    const { name, value } = e.target
    setFilters(prev => ({
      ...prev,
      [name]: value,
      // Reset sub-county when county changes
      ...(name === 'county' && { subCounty: 'all' })
    }))
    
    // Load sub-counties when county is selected
    if (name === 'county' && value !== 'all') {
      try {
        const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
        const response = await fetch(`${apiBaseUrl}/counties/${value}/sub-counties`)
        const data = await response.json()
        setSubCounties(data.sub_counties)
      } catch (error) {
        console.error('Error fetching sub-counties:', error)
        setSubCounties([])
      }
    } else if (name === 'county' && value === 'all') {
      setSubCounties([])
    }
  }

  const applyFilters = () => {
    let filtered = [...properties]

    // Filter by property type
    if (filters.propertyType !== 'all') {
      filtered = filtered.filter(property => {
        return property.property_type === filters.propertyType
      })
    }

    // Filter by county
    if (filters.county !== 'all') {
      filtered = filtered.filter(property => {
        return property.county_id === parseInt(filters.county)
      })
    }

    // Filter by sub-county
    if (filters.subCounty !== 'all') {
      filtered = filtered.filter(property => {
        return property.sub_county_id === parseInt(filters.subCounty)
      })
    }

    // Filter by price range
    if (filters.priceRange !== 'all') {
      const [min, max] = filters.priceRange.split('-').map(val =>
          val.endsWith('+') ? parseInt(val) : parseInt(val)
      )

      filtered = filtered.filter(property => {
        const price = property.rent_amount || 0
        if (filters.priceRange.endsWith('+')) {
          return price >= min
        } else {
          return price >= min && price <= max
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
                  <option value="apartment">Apartment</option>
                  <option value="house">House</option>
                  <option value="villa">Villa</option>
                  <option value="land">Land</option>
                </select>
              </div>

              <div className="min-w-[200px]">
                <label className="block mb-2 font-medium">County</label>
                <select
                    name="county"
                    value={filters.county}
                    onChange={handleFilterChange}
                    className="w-full p-3 rounded border border-gray-300 bg-white"
                >
                  <option value="all">All Counties</option>
                  {counties.map(county => (
                    <option key={county.id} value={county.id}>{county.name}</option>
                  ))}
                </select>
              </div>

              {subCounties.length > 0 && (
                <div className="min-w-[200px]">
                  <label className="block mb-2 font-medium">Sub-County</label>
                  <select
                      name="subCounty"
                      value={filters.subCounty}
                      onChange={handleFilterChange}
                      className="w-full p-3 rounded border border-gray-300 bg-white"
                  >
                    <option value="all">All Sub-Counties</option>
                    {subCounties.map(subCounty => (
                      <option key={subCounty.id} value={subCounty.id}>{subCounty.name}</option>
                    ))}
                  </select>
                </div>
              )}

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