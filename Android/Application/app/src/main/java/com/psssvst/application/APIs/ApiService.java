package com.psssvst.application.APIs;

import com.psssvst.application.APIs.models.login.LoginRequest;
import com.psssvst.application.APIs.models.login.LoginResponse;

import retrofit2.Call;
import retrofit2.http.Body;
import retrofit2.http.POST;

public interface ApiService {
    @POST("api/auth/login")
    Call<LoginResponse> login(
      @Body LoginRequest request
    );
}
