package com.test.realestate.viewmodel

import androidx.compose.runtime.mutableStateOf
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.test.realestate.data.Property
import com.test.realestate.repository.PropertyRepository
import kotlinx.coroutines.launch

class PropertyDetailsViewModel : ViewModel() {
    private val repository = PropertyRepository()
    
    var property = mutableStateOf<Property?>(null)
        private set
    
    var isLoading = mutableStateOf(false)
        private set
    
    var errorMessage = mutableStateOf<String?>(null)
        private set

    fun loadProperty(id: String) {
        viewModelScope.launch {
            isLoading.value = true
            errorMessage.value = null
            
            try {
                val response = repository.getProperty(id)
                if (response.isSuccessful) {
                    property.value = response.body()
                } else {
                    errorMessage.value = "Failed to load property: ${response.code()}"
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
