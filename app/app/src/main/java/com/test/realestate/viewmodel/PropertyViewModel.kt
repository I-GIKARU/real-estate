package com.test.realestate.viewmodel

import androidx.compose.runtime.mutableStateOf
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.test.realestate.data.Property
import com.test.realestate.repository.PropertyRepository
import kotlinx.coroutines.launch

class PropertyViewModel : ViewModel() {
    private val repository = PropertyRepository()
    
    var properties = mutableStateOf<List<Property>>(emptyList())
        private set
    
    var isLoading = mutableStateOf(false)
        private set
    
    var errorMessage = mutableStateOf<String?>(null)
        private set

    init {
        loadProperties()
    }

    fun loadProperties(search: String? = null) {
        viewModelScope.launch {
            isLoading.value = true
            errorMessage.value = null
            
            try {
                val response = repository.getProperties(search = search)
                if (response.isSuccessful) {
                    val body = response.body()
                    properties.value = body?.properties ?: emptyList()
                } else {
                    errorMessage.value = "Failed to load properties: ${response.code()}"
                }
            } catch (e: Exception) {
                errorMessage.value = "Network error: ${e.message}"
                e.printStackTrace()
            } finally {
                isLoading.value = false
            }
        }
    }
}
