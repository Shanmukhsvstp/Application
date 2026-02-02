package com.psssvst.application;

import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.text.SpannableString;
import android.text.TextPaint;
import android.text.style.ClickableSpan;
import android.text.style.ForegroundColorSpan;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;

import androidx.activity.EdgeToEdge;
import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.graphics.Insets;
import androidx.core.view.ViewCompat;
import androidx.core.view.WindowInsetsCompat;

import com.google.android.material.textfield.TextInputEditText;
import com.psssvst.application.APIs.ApiClient;
import com.psssvst.application.APIs.ApiService;
import com.psssvst.application.APIs.models.login.LoginRequest;
import com.psssvst.application.APIs.models.login.LoginResponse;
import com.psssvst.application.APIs.models.signup.SignupRequest;
import com.psssvst.application.APIs.models.signup.SignupResponse;
import com.psssvst.application.Managers.SessionManager;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class Signup extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        EdgeToEdge.enable(this);
        setContentView(R.layout.activity_signup);
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main), (v, insets) -> {
            Insets systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars());
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom);
            return insets;
        });

        TextInputEditText usernameInput = findViewById(R.id.usernameEditText);
        TextInputEditText emailInput = findViewById(R.id.emailEditText);
        TextInputEditText passwordInput = findViewById(R.id.passwordEditText);

        TextView goToLogin = findViewById(R.id.loginRedirectText);
        Button signupBtn = findViewById(R.id.signupBtn);

        ApiService api = ApiClient.getClient().create(ApiService.class);

        initTextView(goToLogin);
        signupBtn.setOnClickListener(v->{
            String username = usernameInput.getText().toString();
            String email = emailInput.getText().toString();
            String password = passwordInput.getText().toString();
            SignupUser(username, email, password, api);
        });
    }
    private void initTextView(TextView textView) {
        String text = "Already have an account? Login";
        int startIndex = text.indexOf("Login");
        int endIndex = startIndex + "Login".length();
        ClickableSpan clickableSpan = new ClickableSpan() {
            @Override
            public void onClick(@NonNull View view) {
                startActivity(new Intent(Signup.this, Login.class));
                finish();
            }

            @Override
            public void updateDrawState(@NonNull TextPaint ds) {
                super.updateDrawState(ds);
                ds.setColor(Color.BLUE);
                ds.setUnderlineText(false);
            }
        };
        SpannableString spannableString = new SpannableString(text);
        spannableString.setSpan(
                clickableSpan,
                startIndex,
                endIndex,
                SpannableString.SPAN_EXCLUSIVE_EXCLUSIVE
        );

        textView.setText(text);
    }

    private void SignupUser(String username, String email, String password, ApiService api) {

        SignupRequest req = new SignupRequest(username, email, password);

        api.signup(req).enqueue(new Callback<SignupResponse>() {
            @Override
            public void onResponse(Call<SignupResponse> call, Response<SignupResponse> response) {
                if (response.isSuccessful()){
                    SignupResponse signupResponse = response.body();
                    Toast.makeText(Signup.this, "Signup Successful", Toast.LENGTH_SHORT).show();
                    if (signupResponse != null && signupResponse.getToken() != null) {
                        SessionManager sessionManager = new SessionManager(Signup.this);
                        sessionManager.saveToken(signupResponse.getToken());
                    } else {
                        Toast.makeText(Signup.this, "Something's wrong on our end...", Toast.LENGTH_SHORT).show();
                    }
                } else {
                    try {
                        Toast.makeText(Signup.this, response.errorBody().string(), Toast.LENGTH_SHORT).show();
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                }
            }

            @Override
            public void onFailure(Call<SignupResponse> call, Throwable t) {
                Toast.makeText(Signup.this, t.getLocalizedMessage(), Toast.LENGTH_SHORT).show();
            }
        });

    }
}