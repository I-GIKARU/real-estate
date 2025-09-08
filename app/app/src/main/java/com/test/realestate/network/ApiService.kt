package com.test.realestate.network

import com.test.realestate.data.PropertiesResponse
import com.test.realestate.data.Property
import retrofit2.Response
import retrofit2.http.GET
import retrofit2.http.Path
import retrofit2.http.Query

interface ApiService {
    @GET("api/v1/properties")
    suspend fun getProperties(
        @Query("page") page: Int = 1,
        @Query("limit") limit: Int = 10,
        @Query("search") search: String? = null
    ): Response<PropertiesResponse>

    @GET("api/v1/properties/{id}")
    suspend fun getProperty(@Path("id") id: String): Response<Property>
}
