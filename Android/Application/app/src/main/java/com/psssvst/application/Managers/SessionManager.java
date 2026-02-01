package com.psssvst.application.Managers;

import android.content.Context;
import android.content.SharedPreferences;

public class SessionManager {
    String prefName = "prefs";

    String tokenName = "auth_token";

    private final SharedPreferences prefs;
    private final SharedPreferences.Editor editor;

    public SessionManager(Context context) {
        prefs = context.getSharedPreferences(prefName, Context.MODE_PRIVATE);
        editor = prefs.edit();
    }
    public void saveToken(String token) {
        editor.putString(tokenName, token);
        editor.apply();
    }

    public String getToken() {
        return prefs.getString(tokenName, null);
    }
}
