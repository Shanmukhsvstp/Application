package com.psssvst.application.APIs;

import com.psssvst.application.APIs.models.login.LoginRequest;
import com.psssvst.application.APIs.models.login.LoginResponse;
import com.psssvst.application.APIs.models.signup.SignupRequest;
import com.psssvst.application.APIs.models.signup.SignupResponse;

import retrofit2.Call;
import retrofit2.http.Body;
import retrofit2.http.GET;
import retrofit2.http.POST;

public interface ApiService {
    @POST("/api/auth/login")
    Call<LoginResponse> login(
            @Body LoginRequest request
    );

    @POST("/api/auth/signup")
    Call<SignupResponse> signup(
            @Body SignupRequest request
    );

    @GET("/api/auth/validate")
    Call<com.psssvst.application.APIs.models.AuthValidation.Response> validateAuthentication();
}
