package com.psssvst.application;

import android.content.Intent;
import android.os.Bundle;

import androidx.activity.EdgeToEdge;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.graphics.Insets;
import androidx.core.view.ViewCompat;
import androidx.core.view.WindowInsetsCompat;

import com.psssvst.application.APIs.ApiClient;
import com.psssvst.application.APIs.ApiService;
import com.psssvst.application.APIs.models.AuthValidation.Response;
import com.psssvst.application.Managers.SessionManager;

import retrofit2.Call;
import retrofit2.Callback;

public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        EdgeToEdge.enable(this);
        setContentView(R.layout.activity_main);
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main), (v, insets) -> {
            Insets systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars());
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom);
            return insets;
        });
        ApiClient.init(this);
        ApiService apiService = ApiClient.getClient().create(ApiService.class);
        apiService.validateAuthentication().enqueue(new Callback<Response>() {
            @Override
            public void onResponse(Call<Response> call, retrofit2.Response<Response> response) {
                if (response.isSuccessful()) {
                    com.psssvst.application.APIs.models.AuthValidation.Response response1 = response.body();
                    if (response1 != null && response1.getToken() != null) {
                        SessionManager sessionManager = new SessionManager(MainActivity.this);
                        sessionManager.saveToken(response1.getToken());
                    }

                    startActivity(new Intent(MainActivity.this, Home.class));
                    finish();
                }
                else {
                    startActivity(new Intent(MainActivity.this, Signup.class));
                    finish();
                }
            }


            @Override
            public void onFailure(Call<Response> call, Throwable t) {
                startActivity(new Intent(MainActivity.this, Signup.class));
                finish();
            }
        });
        startActivity(new Intent(this, Signup.class));
        finish();
    }
}