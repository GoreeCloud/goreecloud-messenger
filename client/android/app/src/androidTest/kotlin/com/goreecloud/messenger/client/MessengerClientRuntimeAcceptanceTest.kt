package com.goreecloud.messenger.client

import android.Manifest
import android.content.Context
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.test.core.app.ActivityScenario
import androidx.test.core.app.ApplicationProvider
import androidx.test.ext.junit.runners.AndroidJUnit4
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test
import org.junit.runner.RunWith

/**
 * Emulator acceptance for the disconnected Development shell only.
 *
 * This does not establish account/session binding, Data transport, message persistence, E2EE,
 * carrier messaging, calling, accessibility certification, physical-device acceptance, or RC.
 */
@RunWith(AndroidJUnit4::class)
class MessengerClientRuntimeAcceptanceTest {
    @Test
    fun launchPreservesVisibleDevelopmentBoundaryAndRestrictedAuthority() {
        val context = ApplicationProvider.getApplicationContext<Context>()
        val packageInfo = context.packageManager.getPackageInfo(context.packageName, 0x00001000)
        val requestedPermissions = packageInfo.requestedPermissions.orEmpty().toSet()
        val forbiddenPermissions = setOf(
            Manifest.permission.INTERNET,
            Manifest.permission.READ_CONTACTS,
            Manifest.permission.READ_SMS,
            Manifest.permission.SEND_SMS,
            Manifest.permission.RECEIVE_SMS,
            Manifest.permission.CALL_PHONE,
            Manifest.permission.RECORD_AUDIO,
            Manifest.permission.CAMERA,
        )
        assertTrue(
            "Disconnected Development shell must not request live communication authority",
            requestedPermissions.intersect(forbiddenPermissions).isEmpty(),
        )

        ActivityScenario.launch(MessengerClientActivity::class.java).use { scenario ->
            scenario.onActivity { activity ->
                val visibleText = collectText(activity.window.decorView)
                assertTrue(visibleText.any { it.contains("Native Android Development preview") })
                assertTrue(visibleText.any { it.contains("Disconnected shell") })
                assertTrue(visibleText.any { it.contains("Development boundary") })
                assertTrue(visibleText.any { it.contains("Not Release Candidate") })
                assertTrue(visibleText.any { it.contains("Provenance examples") })

                // A disconnected provenance preview must not grow a live message-send control.
                assertFalse(visibleText.any { it.trim().equals("Send", ignoreCase = true) })
            }
        }
    }

    private fun collectText(view: View): List<String> = when (view) {
        is TextView -> listOf(view.text?.toString().orEmpty())
        is ViewGroup -> buildList {
            repeat(view.childCount) { index -> addAll(collectText(view.getChildAt(index))) }
        }
        else -> emptyList()
    }
}
