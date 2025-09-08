package com.test.realestate.data

import com.google.gson.annotations.SerializedName

data class Property(
    val id: String,
    @SerializedName("agent_id") val agentId: String,
    val title: String,
    val description: String,
    @SerializedName("property_type") val propertyType: String,
    val bedrooms: Int,
    val bathrooms: Int,
    @SerializedName("square_meters") val squareMeters: Int,
    @SerializedName("rent_amount") val rentAmount: Double,
    @SerializedName("deposit_amount") val depositAmount: Double? = null,
    @SerializedName("county_id") val countyId: Int,
    @SerializedName("sub_county_id") val subCountyId: Int,
    @SerializedName("location_details") val locationDetails: String,
    val amenities: Map<String, Any>? = null,
    @SerializedName("utilities_included") val utilitiesIncluded: Map<String, Any>? = null,
    @SerializedName("parking_spaces") val parkingSpaces: Int,
    @SerializedName("is_furnished") val isFurnished: Boolean,
    @SerializedName("is_available") val isAvailable: Boolean,
    @SerializedName("created_at") val createdAt: String,
    @SerializedName("updated_at") val updatedAt: String,
    val county: County,
    @SerializedName("sub_county") val subCounty: SubCounty,
    val agent: Agent,
    val images: List<PropertyImage> = emptyList()
)

data class PropertyImage(
    val id: String,
    @SerializedName("property_id") val propertyId: String,
    @SerializedName("image_url") val imageUrl: String,
    @SerializedName("secure_url") val secureUrl: String,
    @SerializedName("public_id") val publicId: String,
    @SerializedName("is_primary") val isPrimary: Boolean,
    @SerializedName("display_order") val displayOrder: Int,
    val width: Int,
    val height: Int,
    val format: String,
    val bytes: Int,
    @SerializedName("created_at") val createdAt: String
)

data class County(
    val id: Int,
    val name: String,
    val code: String,
    @SerializedName("created_at") val createdAt: String
)

data class SubCounty(
    val id: Int,
    @SerializedName("county_id") val countyId: Int,
    val name: String,
    @SerializedName("created_at") val createdAt: String
)

data class Agent(
    val id: String,
    val email: String,
    @SerializedName("first_name") val firstName: String,
    @SerializedName("last_name") val lastName: String,
    @SerializedName("phone_number") val phoneNumber: String,
    @SerializedName("user_type") val userType: String,
    @SerializedName("is_verified") val isVerified: Boolean,
    @SerializedName("is_approved") val isApproved: Boolean,
    @SerializedName("approved_at") val approvedAt: String?,
    @SerializedName("approved_by") val approvedBy: String?,
    @SerializedName("is_active") val isActive: Boolean,
    @SerializedName("created_at") val createdAt: String,
    @SerializedName("updated_at") val updatedAt: String
)

data class PropertiesResponse(
    val properties: List<Property>,
    val filters: Map<String, Any>? = null
)
