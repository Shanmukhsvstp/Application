package com.psssvst.application.APIs;

import com.psssvst.application.APIs.models.login.LoginRequest;
import com.psssvst.application.APIs.models.login.LoginResponse;
import com.psssvst.application.APIs.models.signup.SignupRequest;
import com.psssvst.application.APIs.models.signup.SignupResponse;

import retrofit2.Call;
import retrofit2.http.Body;
import retrofit2.http.POST;

public interface ApiService {
    Call<LoginResponse> login(
            @Body LoginRequest request
    );
    Call<SignupResponse> signup(
            @Body SignupRequest request
    );
}
