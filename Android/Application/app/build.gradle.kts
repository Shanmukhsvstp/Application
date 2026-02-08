plugins {
    alias(libs.plugins.android.application)
}

android {
    namespace = "com.psssvst.application"
    compileSdk = 36

    defaultConfig {
        applicationId = "com.psssvst.application"
        minSdk = 24
        targetSdk = 36
        versionCode = 1
        versionName = "1.0"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    buildTypes {
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    buildFeatures {
        viewBinding = true
    }
}

dependencies {

    implementation(libs.appcompat)
    implementation(libs.material)
    implementation(libs.activity)
    implementation(libs.constraintlayout)
    implementation(libs.lifecycle.livedata.ktx)
    implementation(libs.lifecycle.viewmodel.ktx)
    implementation(libs.navigation.fragment)
    implementation(libs.navigation.ui)
    testImplementation(libs.junit)
    androidTestImplementation(libs.ext.junit)
    androidTestImplementation(libs.espresso.core)

    // Retrofit (API layer)
    implementation("com.squareup.retrofit2:retrofit:3.0.0")

    // Gson converter (JSON ↔ objects)
    implementation("com.squareup.retrofit2:converter-gson:3.0.0")

    // OkHttp (HTTP engine)
    implementation("com.squareup.okhttp3:okhttp:5.3.2")

    // OkHttp logging interceptor
    implementation("com.squareup.okhttp3:logging-interceptor:5.3.2")
}