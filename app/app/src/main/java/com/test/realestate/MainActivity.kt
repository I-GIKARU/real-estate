package com.test.realestate

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.test.realestate.screens.PropertyDetailsScreen
import com.test.realestate.screens.PropertyListScreen
import com.test.realestate.ui.theme.RealestateTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            RealestateTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    RealEstateApp(
                        modifier = Modifier.padding(innerPadding)
                    )
                }
            }
        }
    }
}

@Composable
fun RealEstateApp(
    modifier: Modifier = Modifier,
    navController: NavHostController = rememberNavController()
) {
    NavHost(
        navController = navController,
        startDestination = "property_list",
        modifier = modifier
    ) {
        composable("property_list") {
            PropertyListScreen(
                onPropertyClick = { propertyId ->
                    navController.navigate("property_details/$propertyId")
                }
            )
        }
        composable("property_details/{propertyId}") { backStackEntry ->
            val propertyId = backStackEntry.arguments?.getString("propertyId") ?: ""
            PropertyDetailsScreen(
                propertyId = propertyId,
                onBackClick = {
                    navController.popBackStack()
                }
            )
        }
    }
}