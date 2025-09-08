package com.test.realestate.repository

import com.test.realestate.data.Property
import com.test.realestate.network.NetworkModule

class PropertyRepository {
    private val apiService = NetworkModule.apiService

    suspend fun getProperties(page: Int = 1, search: String? = null) = 
        apiService.getProperties(page = page, search = search)

    suspend fun getProperty(id: String) = 
        apiService.getProperty(id)
}
