package com.goreecloud.messenger.client

import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.test.core.app.ActivityScenario
import androidx.test.ext.junit.runners.AndroidJUnit4
import org.junit.Assert.assertTrue
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class MessengerClientActivityRuntimeTest {
    @Test
    fun recreationPreservesDisconnectedReadOnlyBoundary() {
        ActivityScenario.launch(MessengerClientActivity::class.java).use { scenario ->
            assertDisconnectedBoundary(scenario)
            scenario.recreate()
            assertDisconnectedBoundary(scenario)
        }
    }

    private fun assertDisconnectedBoundary(scenario: ActivityScenario<MessengerClientActivity>) {
        scenario.onActivity { activity ->
            val root = activity.findViewById<ViewGroup>(android.R.id.content)
            val renderedText = collectText(root)
            assertTrue(renderedText.contains("Development"))
            assertTrue(renderedText.contains("Data send unavailable"))
            assertTrue(renderedText.contains("Not Release Candidate"))
        }
    }

    private fun collectText(view: View): String = buildString {
        when (view) {
            is TextView -> append(view.text).append('\n')
            is ViewGroup -> {
                for (index in 0 until view.childCount) {
                    append(collectText(view.getChildAt(index)))
                }
            }
        }
    }
}
