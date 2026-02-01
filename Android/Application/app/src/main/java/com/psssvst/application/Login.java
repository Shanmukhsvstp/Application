package com.psssvst.application;

import android.os.Bundle;
import android.widget.Toast;

import androidx.activity.EdgeToEdge;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.graphics.Insets;
import androidx.core.view.ViewCompat;
import androidx.core.view.WindowInsetsCompat;

import com.google.android.material.button.MaterialButton;
import com.google.android.material.textfield.TextInputEditText;
import com.psssvst.application.APIs.ApiClient;
import com.psssvst.application.APIs.ApiService;
import com.psssvst.application.APIs.models.login.LoginRequest;
import com.psssvst.application.APIs.models.login.LoginResponse;
import com.psssvst.application.Managers.SessionManager;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;


public class Login extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        EdgeToEdge.enable(this);
        setContentView(R.layout.activity_login);
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main), (v, insets) -> {
            Insets systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars());
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom);
            return insets;
        });

        TextInputEditText usernameInput = findViewById(R.id.usernameEditText);
        TextInputEditText emailInput = findViewById(R.id.emailEditText);
        TextInputEditText passwordInput = findViewById(R.id.passwordEditText);
        MaterialButton loginBtn = findViewById(R.id.loginBtn);

        ApiService api = ApiClient.getClient().create(ApiService.class);


        loginBtn.setOnClickListener(v -> {
            String username = usernameInput.getText().toString();
            String email = emailInput.getText().toString();
            String password = passwordInput.getText().toString();

            LoginRequest req = new LoginRequest(username, email, password);

            api.login(req).enqueue(new Callback<LoginResponse>() {
                @Override
                public void onResponse(Call<LoginResponse> call, Response<LoginResponse> response) {
                    if (response.isSuccessful()){
                        LoginResponse loginResponse = response.body();
                        Toast.makeText(Login.this, "Login Successful", Toast.LENGTH_SHORT).show();
                        if (loginResponse != null && loginResponse.getToken() != null) {
                            SessionManager sessionManager = new SessionManager(Login.this);
                            sessionManager.saveToken(loginResponse.getToken());
                        } else {
                            Toast.makeText(Login.this, "Something's wrong on our end...", Toast.LENGTH_SHORT).show();
                        }
                    } else {
                        try {
                            Toast.makeText(Login.this, response.errorBody().string(), Toast.LENGTH_SHORT).show();
                        } catch (Exception e) {
                            e.printStackTrace();
                        }
                    }
                }

                @Override
                public void onFailure(Call<LoginResponse> call, Throwable t) {
                    Toast.makeText(Login.this, t.getLocalizedMessage(), Toast.LENGTH_SHORT).show();
                }
            });
        });



    }
}