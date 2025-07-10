package utils

import "fmt"

// KenyanPropertyAmenities contains common amenities in Kenyan properties
var KenyanPropertyAmenities = map[string][]string{
	"security": {
		"24/7 Security",
		"CCTV Surveillance",
		"Electric Fence",
		"Security Guards",
		"Gated Community",
		"Access Control",
		"Perimeter Wall",
		"Security Lights",
	},
	"utilities": {
		"Borehole Water",
		"Mains Water",
		"Backup Generator",
		"Solar Water Heating",
		"Solar Power",
		"Prepaid Electricity",
		"Garbage Collection",
		"Internet Ready",
		"DSTV Ready",
		"Intercom System",
	},
	"kitchen": {
		"Modern Kitchen",
		"Kitchen Cabinets",
		"Granite Countertops",
		"Gas Cooker",
		"Electric Cooker",
		"Microwave",
		"Refrigerator",
		"Dishwasher",
		"Pantry",
		"Breakfast Bar",
	},
	"bathroom": {
		"En-suite Bathroom",
		"Guest Toilet",
		"Bathtub",
		"Shower Cubicle",
		"Hot Water",
		"Modern Fixtures",
		"Vanity Unit",
		"Bidet",
	},
	"outdoor": {
		"Garden",
		"Balcony",
		"Terrace",
		"Rooftop Access",
		"Compound Parking",
		"Carport",
		"Garage",
		"Servant Quarter",
		"Laundry Area",
		"Outdoor Kitchen",
		"Barbecue Area",
		"Children Play Area",
	},
	"flooring": {
		"Tiled Floors",
		"Wooden Floors",
		"Marble Floors",
		"Terrazzo Floors",
		"Carpeted Floors",
		"Ceramic Tiles",
		"Granite Floors",
	},
	"facilities": {
		"Swimming Pool",
		"Gym/Fitness Center",
		"Clubhouse",
		"Tennis Court",
		"Basketball Court",
		"Children Playground",
		"Jogging Track",
		"Spa",
		"Sauna",
		"Conference Room",
		"Business Center",
		"Lift/Elevator",
		"Backup Water Tank",
		"Waste Management",
	},
	"location": {
		"Near Shopping Mall",
		"Near School",
		"Near Hospital",
		"Near Public Transport",
		"Near Highway",
		"Near Airport",
		"Near CBD",
		"Quiet Neighborhood",
		"Residential Area",
		"Commercial Area",
		"Mixed Development",
	},
}

// KenyanPropertyTypes contains property types common in Kenya
var KenyanPropertyTypes = map[string]string{
	"bedsitter":   "A single room with a small kitchen area and private bathroom",
	"studio":      "Open plan living space with separate bathroom",
	"apartment":   "Multi-room unit in a building with shared facilities",
	"maisonette":  "Two-story apartment or house unit",
	"bungalow":    "Single-story detached house",
	"villa":       "Large, luxurious house often in gated community",
	"townhouse":   "Multi-story house sharing walls with neighbors",
	"penthouse":   "Luxury apartment on the top floor of a building",
	"duplex":      "Two-unit building with separate entrances",
	"commercial":  "Property for business use (shops, offices, warehouses)",
}

// KenyanUtilities contains utilities commonly included in Kenyan rentals
var KenyanUtilities = map[string]bool{
	"water":           true,
	"electricity":     false, // Usually separate
	"garbage":         true,
	"security":        true,
	"internet":        false,
	"cable_tv":        false,
	"gas":            false,
	"parking":        true,
	"garden_service": false,
	"cleaning":       false,
}

// KenyanRentalTerms contains common rental terms in Kenya
var KenyanRentalTerms = map[string]string{
	"deposit":        "Usually 1-2 months rent paid upfront as security deposit",
	"advance_rent":   "1-3 months rent paid in advance",
	"agent_fee":      "Usually 50% of one month's rent paid to agent",
	"lease_period":   "Typically 1-2 years with option to renew",
	"notice_period":  "Usually 1-3 months notice required to vacate",
	"maintenance":    "Tenant responsible for minor repairs, landlord for major",
	"utilities":      "Tenant usually pays electricity, water may be included",
	"pets":          "Usually not allowed or require additional deposit",
}

// GetAmenitiesByCategory returns amenities for a specific category
func GetAmenitiesByCategory(category string) []string {
	if amenities, exists := KenyanPropertyAmenities[category]; exists {
		return amenities
	}
	return []string{}
}

// GetAllAmenities returns all available amenities
func GetAllAmenities() map[string][]string {
	return KenyanPropertyAmenities
}

// GetPropertyTypeDescription returns description for a property type
func GetPropertyTypeDescription(propertyType string) string {
	if description, exists := KenyanPropertyTypes[propertyType]; exists {
		return description
	}
	return ""
}

// GetDefaultUtilities returns default utilities typically included
func GetDefaultUtilities() map[string]bool {
	return KenyanUtilities
}

// GetRentalTerms returns common rental terms in Kenya
func GetRentalTerms() map[string]string {
	return KenyanRentalTerms
}

// ValidateKenyanPhoneNumber validates a Kenyan phone number
func ValidateKenyanPhoneNumber(phone string) bool {
	// Remove any non-digit characters
	cleaned := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}

	// Check various Kenyan phone number formats
	switch {
	case len(cleaned) == 10 && cleaned[:1] == "0":
		// 0712345678, 0722345678, etc.
		return cleaned[1:3] == "71" || cleaned[1:3] == "72" || cleaned[1:3] == "73" ||
			   cleaned[1:3] == "74" || cleaned[1:3] == "75" || cleaned[1:3] == "76" ||
			   cleaned[1:3] == "77" || cleaned[1:3] == "78" || cleaned[1:3] == "79" ||
			   cleaned[1:3] == "70"
	case len(cleaned) == 9:
		// 712345678, 722345678, etc.
		return cleaned[:2] == "71" || cleaned[:2] == "72" || cleaned[:2] == "73" ||
			   cleaned[:2] == "74" || cleaned[:2] == "75" || cleaned[:2] == "76" ||
			   cleaned[:2] == "77" || cleaned[:2] == "78" || cleaned[:2] == "79" ||
			   cleaned[:2] == "70"
	case len(cleaned) == 12 && cleaned[:3] == "254":
		// 254712345678, 254722345678, etc.
		return cleaned[3:5] == "71" || cleaned[3:5] == "72" || cleaned[3:5] == "73" ||
			   cleaned[3:5] == "74" || cleaned[3:5] == "75" || cleaned[3:5] == "76" ||
			   cleaned[3:5] == "77" || cleaned[3:5] == "78" || cleaned[3:5] == "79" ||
			   cleaned[3:5] == "70"
	}
	return false
}

// FormatKenyanCurrency formats amount in Kenyan Shillings
func FormatKenyanCurrency(amount float64) string {
	if amount >= 1000000 {
		return fmt.Sprintf("KES %.1fM", amount/1000000)
	} else if amount >= 1000 {
		return fmt.Sprintf("KES %.0fK", amount/1000)
	}
	return fmt.Sprintf("KES %.0f", amount)
}

// GetPopularKenyanAreas returns popular residential areas by county
func GetPopularKenyanAreas() map[string][]string {
	return map[string][]string{
		"Nairobi": {
			"Westlands", "Karen", "Kilimani", "Lavington", "Kileleshwa",
			"Runda", "Muthaiga", "Spring Valley", "Loresho", "Gigiri",
			"Parklands", "Eastleigh", "South B", "South C", "Langata",
			"Kasarani", "Roysambu", "Thika Road", "Ngong Road", "Waiyaki Way",
		},
		"Mombasa": {
			"Nyali", "Bamburi", "Shanzu", "Diani", "Likoni",
			"Tudor", "Kizingo", "Ganjoni", "Mtwapa", "Kilifi",
		},
		"Kiambu": {
			"Thika", "Ruiru", "Juja", "Kikuyu", "Limuru",
			"Kiambu Town", "Githunguri", "Gatundu", "Lari",
		},
		"Nakuru": {
			"Nakuru Town", "Naivasha", "Gilgil", "Molo", "Njoro",
			"Bahati", "Rongai", "Subukia",
		},
		"Uasin Gishu": {
			"Eldoret", "Moiben", "Soy", "Turbo", "Kapseret",
		},
	}
}

