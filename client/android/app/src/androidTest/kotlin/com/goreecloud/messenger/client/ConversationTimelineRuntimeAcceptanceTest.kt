package com.goreecloud.messenger.client

import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.test.core.app.ActivityScenario
import androidx.test.ext.junit.runners.AndroidJUnit4
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class ConversationTimelineRuntimeAcceptanceTest {
    @Test
    fun conversationTimelineIsVisibleButComposerRemainsUnavailable() {
        ActivityScenario.launch(MessengerClientActivity::class.java).use { scenario ->
            scenario.onActivity { activity ->
                val visibleText = collectText(activity.window.decorView)
                assertTrue(visibleText.any { it.trim() == "Conversation preview" })
                assertTrue(visibleText.any { it.contains("Are we still on for six?") })
                assertTrue(visibleText.any { it.contains("Replying to Alex") })
                assertTrue(visibleText.any { it.contains("Delivered") })
                assertTrue(visibleText.any { it.contains("Message composer unavailable") })
                assertTrue(visibleText.any { it.contains("Protection not verified") })
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
