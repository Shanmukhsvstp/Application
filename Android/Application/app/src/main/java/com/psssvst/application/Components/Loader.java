package com.psssvst.application.Components;

import android.app.Dialog;
import android.content.Context;
import android.graphics.drawable.Animatable;
import android.graphics.drawable.Drawable;
import android.view.Window;
import android.view.WindowManager;
import android.widget.ImageView;
import android.widget.TextView;

import com.psssvst.application.R;

public class Loader {

    private static Loader instance;
    private Dialog dialog;

    private Loader() {}

    public static synchronized Loader getInstance() {
        if (instance == null) {
            instance = new Loader();
        }
        return instance;
    }

    public void show(Context context) {
        if (dialog != null && dialog.isShowing()) return;

        dialog = new Dialog(context);
        dialog.requestWindowFeature(Window.FEATURE_NO_TITLE);
        dialog.setContentView(R.layout.loader);
        dialog.setCancelable(false);

        if (dialog.getWindow() != null) {
            dialog.getWindow().setBackgroundDrawableResource(android.R.color.transparent);
            dialog.getWindow().clearFlags(WindowManager.LayoutParams.FLAG_DIM_BEHIND);
        }

        ImageView loader = dialog.findViewById(R.id.loader);
        Drawable drawable = loader.getDrawable();

        if (drawable instanceof Animatable) {
            ((Animatable) drawable).start();
        }

        dialog.show();
    }

    public void show(Context context, String message) {
        if (dialog != null && dialog.isShowing()) return;

        dialog = new Dialog(context);
        dialog.requestWindowFeature(Window.FEATURE_NO_TITLE);
        dialog.setContentView(R.layout.loader);
        dialog.setCancelable(false);

        if (dialog.getWindow() != null) {
            dialog.getWindow().setBackgroundDrawableResource(android.R.color.transparent);
            dialog.getWindow().clearFlags(WindowManager.LayoutParams.FLAG_DIM_BEHIND);
        }

        ImageView loader = dialog.findViewById(R.id.loader);
        TextView messageView = dialog.findViewById(R.id.loadingText);
        Drawable drawable = loader.getDrawable();

        messageView.setText(message);

        if (drawable instanceof Animatable) {
            ((Animatable) drawable).start();
        }

        dialog.show();
    }

    public void hide() {
        if (dialog == null) return;

        if (dialog.isShowing()) {
            dialog.dismiss();
        }
        dialog = null;
    }
}
