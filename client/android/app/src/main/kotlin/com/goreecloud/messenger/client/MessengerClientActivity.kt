package com.goreecloud.messenger.client

import android.app.Activity
import android.content.res.Configuration
import android.graphics.Typeface
import android.graphics.drawable.GradientDrawable
import android.os.Bundle
import android.view.View
import android.view.ViewGroup
import android.widget.LinearLayout
import android.widget.ScrollView
import android.widget.TextView

class MessengerClientActivity : Activity() {
    private val isDark: Boolean
        get() = resources.configuration.uiMode and Configuration.UI_MODE_NIGHT_MASK ==
            Configuration.UI_MODE_NIGHT_YES

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        window.navigationBarColor = palette().canvas
        window.statusBarColor = palette().canvas
        setContentView(buildContent())
    }

    private fun buildContent(): View {
        val colors = palette()
        val root = ScrollView(this).apply {
            setBackgroundColor(colors.canvas)
            isFillViewport = true
        }
        val content = LinearLayout(this).apply {
            orientation = LinearLayout.VERTICAL
            val gutter = dp(GlazeClientTokens.ScreenGutterDp)
            setPadding(gutter, dp(28), gutter, dp(36))
            layoutParams = ViewGroup.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.WRAP_CONTENT,
            )
        }

        content.addView(text("GoreeCloud Messenger", 30f, colors.text, Typeface.BOLD))
        content.addView(spacer(8))
        content.addView(text(getString(R.string.development_title), 17f, colors.text, Typeface.BOLD))
        content.addView(spacer(4))
        content.addView(text(getString(R.string.development_summary), 15f, colors.muted, Typeface.NORMAL))
        content.addView(spacer(14))
        content.addView(
            surface(
                title = "Development boundary",
                body = getString(R.string.development_detail),
                colors = colors,
            ),
        )

        content.addView(spacer(22))
        content.addView(text(getString(R.string.readiness_heading), 20f, colors.text, Typeface.BOLD))
        content.addView(spacer(5))
        content.addView(text(getString(R.string.readiness_summary), 14f, colors.muted, Typeface.NORMAL))
        content.addView(spacer(12))
        content.addView(
            surface(
                title = "Data send unavailable",
                body = readinessExplanation(disconnectedReadiness()),
                colors = colors,
            ),
        )

        content.addView(spacer(22))
        content.addView(text(getString(R.string.provenance_heading), 20f, colors.text, Typeface.BOLD))
        content.addView(spacer(5))
        content.addView(text(getString(R.string.provenance_summary), 14f, colors.muted, Typeface.NORMAL))
        content.addView(spacer(12))

        val examples = listOf(
            CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.E2EE_ACTIVE,
            ) to "Allowed only after the client can verify the GoreeCloud E2EE state.",
            CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.UNKNOWN,
            ) to "Used when Data transport is known but protection has not been verified.",
            CommunicationProvenance(
                CommunicationTransport.SMS,
                CommunicationProtection.UNKNOWN,
            ) to "Carrier transport remains visibly SMS; this client cannot label it GoreeCloud E2EE.",
            CommunicationProvenance(
                CommunicationTransport.RCS,
                CommunicationProtection.UNKNOWN,
            ) to "RCS is shown only as a transport example and is not claimed available in this build.",
        )
        examples.forEachIndexed { index, (provenance, explanation) ->
            content.addView(
                surface(
                    title = provenance.displayLabel(),
                    body = explanation,
                    colors = colors,
                ),
            )
            if (index != examples.lastIndex) content.addView(spacer(10))
        }

        content.addView(spacer(22))
        content.addView(text(getString(R.string.platform_heading), 20f, colors.text, Typeface.BOLD))
        content.addView(spacer(10))
        content.addView(
            surface(
                title = "Not Release Candidate",
                body = getString(R.string.platform_summary),
                colors = colors,
            ),
        )

        root.addView(content)
        return root
    }

    /**
     * This Development shell intentionally has no live authorities to supply these prerequisites.
     * Keep every value unknown rather than synthesizing readiness from the mere existence of client
     * code or from a future transport connection.
     */
    private fun disconnectedReadiness(): DataMessagingReadiness.Result =
        DataMessagingReadiness.evaluate(
            DataMessagingReadiness.Evidence(
                identity = DataMessagingReadiness.IdentityState.UNKNOWN,
                conversationAccess = DataMessagingReadiness.ConversationAccessState.UNKNOWN,
                transport = DataMessagingReadiness.DataTransportState.UNKNOWN,
                cryptography = DataMessagingReadiness.CryptographicState.UNKNOWN,
            ),
        )

    private fun readinessExplanation(result: DataMessagingReadiness.Result): String = when (result) {
        is DataMessagingReadiness.Result.Ready ->
            "Data messaging prerequisites are independently verified (${result.provenance.displayLabel()})."

        is DataMessagingReadiness.Result.Blocked -> {
            val labels = listOfNotNull(
                "Identity authentication".takeIf {
                    DataMessagingReadiness.BlockReason.IDENTITY_NOT_AUTHENTICATED in result.reasons
                },
                "conversation authorization".takeIf {
                    DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED in result.reasons
                },
                "GoreeCloud Data transport".takeIf {
                    DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE in result.reasons
                },
                "verified active E2EE".takeIf {
                    DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE in result.reasons
                },
            )
            "Missing verified prerequisites: ${labels.joinToString(" · ")}. No Send action is exposed."
        }
    }

    private fun surface(title: String, body: String, colors: Palette): View =
        LinearLayout(this).apply {
            orientation = LinearLayout.VERTICAL
            val padding = dp(18)
            setPadding(padding, padding, padding, padding)
            minimumHeight = dp(GlazeClientTokens.InteractionFloorDp)
            background = GradientDrawable().apply {
                shape = GradientDrawable.RECTANGLE
                cornerRadius = dp(GlazeClientTokens.SurfaceRadiusDp).toFloat()
                setColor(colors.surface)
                setStroke(dp(1), colors.border)
            }
            addView(text(title, 16f, colors.text, Typeface.BOLD))
            addView(spacer(6))
            addView(text(body, 14f, colors.muted, Typeface.NORMAL))
        }

    private fun text(value: String, sizeSp: Float, color: Int, style: Int): TextView =
        TextView(this).apply {
            text = value
            textSize = sizeSp
            setTextColor(color)
            typeface = Typeface.create(Typeface.DEFAULT, style)
            setLineSpacing(0f, 1.08f)
        }

    private fun spacer(heightDp: Int): View = View(this).apply {
        layoutParams = LinearLayout.LayoutParams(
            ViewGroup.LayoutParams.MATCH_PARENT,
            dp(heightDp),
        )
    }

    private fun dp(value: Int): Int = (value * resources.displayMetrics.density).toInt()

    private fun dp(value: Float): Int = (value * resources.displayMetrics.density).toInt()

    private fun palette(): Palette = if (isDark) {
        Palette(
            canvas = GlazeClientTokens.DarkCanvas.toInt(),
            surface = GlazeClientTokens.DarkSurface.toInt(),
            text = GlazeClientTokens.DarkText.toInt(),
            muted = GlazeClientTokens.DarkMutedText.toInt(),
            border = GlazeClientTokens.DarkBorder.toInt(),
        )
    } else {
        Palette(
            canvas = GlazeClientTokens.LightCanvas.toInt(),
            surface = GlazeClientTokens.LightSurface.toInt(),
            text = GlazeClientTokens.LightText.toInt(),
            muted = GlazeClientTokens.LightMutedText.toInt(),
            border = GlazeClientTokens.LightBorder.toInt(),
        )
    }

    private data class Palette(
        val canvas: Int,
        val surface: Int,
        val text: Int,
        val muted: Int,
        val border: Int,
    )
}
