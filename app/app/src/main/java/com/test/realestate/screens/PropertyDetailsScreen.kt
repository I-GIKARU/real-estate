package com.test.realestate.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.Info
import androidx.compose.material.icons.filled.Place
import androidx.compose.material.icons.filled.Phone
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import coil.compose.AsyncImage
import com.test.realestate.data.Property
import com.test.realestate.viewmodel.PropertyDetailsViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PropertyDetailsScreen(
    propertyId: String,
    onBackClick: () -> Unit,
    viewModel: PropertyDetailsViewModel = viewModel()
) {
    val property by viewModel.property
    val isLoading by viewModel.isLoading
    val errorMessage by viewModel.errorMessage

    LaunchedEffect(propertyId) {
        viewModel.loadProperty(propertyId)
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Property Details") },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Back")
                    }
                }
            )
        }
    ) { paddingValues ->
        when {
            isLoading -> {
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(paddingValues),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator()
                }
            }
            errorMessage != null -> {
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(paddingValues),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = errorMessage!!,
                        color = MaterialTheme.colorScheme.error
                    )
                }
            }
            property != null -> {
                PropertyDetailsContent(
                    property = property!!,
                    modifier = Modifier.padding(paddingValues)
                )
            }
        }
    }
}

@Composable
fun PropertyDetailsContent(
    property: Property,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
    ) {
        // Image Gallery
        if (property.images?.isNotEmpty() == true) {
            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                contentPadding = PaddingValues(16.dp)
            ) {
                items(property.images) { image ->
                    AsyncImage(
                        model = image.secureUrl.ifEmpty { image.imageUrl },
                        contentDescription = property.title,
                        modifier = Modifier
                            .width(300.dp)
                            .height(200.dp),
                        contentScale = ContentScale.Crop
                    )
                }
            }
        }

        // Property Info
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            Text(
                text = property.title,
                style = MaterialTheme.typography.headlineMedium,
                fontWeight = FontWeight.Bold
            )
            
            Spacer(modifier = Modifier.height(8.dp))
            
            Text(
                text = "KSh ${String.format("%,.0f", property.rentAmount)}/month",
                style = MaterialTheme.typography.headlineSmall,
                color = MaterialTheme.colorScheme.primary,
                fontWeight = FontWeight.Bold
            )
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // Property Features
            Row(
                horizontalArrangement = Arrangement.spacedBy(24.dp)
            ) {
                PropertyFeature(
                    icon = Icons.Default.Home,
                    label = "${property.bedrooms} Beds"
                )
                PropertyFeature(
                    icon = Icons.Default.Info,
                    label = "${property.bathrooms} Baths"
                )
                PropertyFeature(
                    icon = Icons.Default.Place,
                    label = "${property.squareMeters}m²"
                )
            }
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // Location
            Card(
                modifier = Modifier.fillMaxWidth()
            ) {
                Column(
                    modifier = Modifier.padding(16.dp)
                ) {
                    Text(
                        text = "Location",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = property.locationDetails)
                    Text(text = "${property.subCounty.name}, ${property.county.name}")
                }
            }
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // Description
            Card(
                modifier = Modifier.fillMaxWidth()
            ) {
                Column(
                    modifier = Modifier.padding(16.dp)
                ) {
                    Text(
                        text = "Description",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = property.description)
                }
            }
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // Agent Info
            Card(
                modifier = Modifier.fillMaxWidth()
            ) {
                Column(
                    modifier = Modifier.padding(16.dp)
                ) {
                    Text(
                        text = "Agent Information",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = "${property.agent.firstName} ${property.agent.lastName}")
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            Icons.Default.Phone,
                            contentDescription = "Phone",
                            modifier = Modifier.size(16.dp)
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(text = property.agent.phoneNumber)
                    }
                }
            }
        }
    }
}

@Composable
fun PropertyFeature(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    label: String
) {
    Row(
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            icon,
            contentDescription = null,
            modifier = Modifier.size(20.dp),
            tint = MaterialTheme.colorScheme.primary
        )
        Spacer(modifier = Modifier.width(8.dp))
        Text(
            text = label,
            style = MaterialTheme.typography.bodyMedium
        )
    }
}
