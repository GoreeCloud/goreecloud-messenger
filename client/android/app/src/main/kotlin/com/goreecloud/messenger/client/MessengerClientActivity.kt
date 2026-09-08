package com.goreecloud.messenger.client

import android.app.Activity
import android.content.res.Configuration
import android.graphics.Typeface
import android.graphics.drawable.GradientDrawable
import android.os.Bundle
import android.view.Gravity
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
        content.addView(surface("Development boundary", getString(R.string.development_detail), colors))

        content.addView(spacer(22))
        content.addView(text("Conversations", 20f, colors.text, Typeface.BOLD))
        content.addView(spacer(5))
        content.addView(text("Native conversation-list structure using Development-only presentation records. No live account or message history is loaded.", 14f, colors.muted, Typeface.NORMAL))
        content.addView(spacer(12))
        val previews = DevelopmentConversationCatalog.previews()
        previews.forEachIndexed { index, preview ->
            content.addView(conversationRow(preview, colors))
            if (index != previews.lastIndex) content.addView(spacer(10))
        }

        content.addView(spacer(22))
        content.addView(text("Conversation preview", 20f, colors.text, Typeface.BOLD))
        content.addView(spacer(5))
        content.addView(text("Read-only native message timeline. Text below is static Development presentation data, not retrieved conversation content.", 14f, colors.muted, Typeface.NORMAL))
        content.addView(spacer(12))
        content.addView(surface("Alex · @alex", "Data · Protection not verified · Development preview", colors))
        content.addView(spacer(10))
        DevelopmentConversationMessageCatalog.messages().forEachIndexed { index, message ->
            content.addView(messageBubble(message, colors))
            if (index != DevelopmentConversationMessageCatalog.messages().lastIndex) content.addView(spacer(8))
        }
        content.addView(spacer(10))
        content.addView(disabledComposer(colors))

        content.addView(spacer(22))
        content.addView(text(getString(R.string.readiness_heading), 20f, colors.text, Typeface.BOLD))
        content.addView(spacer(5))
        content.addView(text(getString(R.string.readiness_summary), 14f, colors.muted, Typeface.NORMAL))
        content.addView(spacer(12))
        content.addView(surface("Data send unavailable", readinessExplanation(disconnectedReadiness()), colors))

        content.addView(spacer(22))
        content.addView(text(getString(R.string.provenance_heading), 20f, colors.text, Typeface.BOLD))
        content.addView(spacer(5))
        content.addView(text(getString(R.string.provenance_summary), 14f, colors.muted, Typeface.NORMAL))
        content.addView(spacer(12))

        val examples = listOf(
            CommunicationProvenance(CommunicationTransport.DATA, CommunicationProtection.E2EE_ACTIVE) to "Allowed only after the client can verify the GoreeCloud E2EE state.",
            CommunicationProvenance(CommunicationTransport.DATA, CommunicationProtection.UNKNOWN) to "Used when Data transport is known but protection has not been verified.",
            CommunicationProvenance(CommunicationTransport.SMS, CommunicationProtection.UNKNOWN) to "Carrier transport remains visibly SMS; this client cannot label it GoreeCloud E2EE.",
            CommunicationProvenance(CommunicationTransport.RCS, CommunicationProtection.UNKNOWN) to "RCS is shown only as a transport example and is not claimed available in this build.",
        )
        examples.forEachIndexed { index, (provenance, explanation) ->
            content.addView(surface(provenance.displayLabel(), explanation, colors))
            if (index != examples.lastIndex) content.addView(spacer(10))
        }

        content.addView(spacer(22))
        content.addView(text(getString(R.string.platform_heading), 20f, colors.text, Typeface.BOLD))
        content.addView(spacer(10))
        content.addView(surface("Not Release Candidate", getString(R.string.platform_summary), colors))

        root.addView(content)
        return root
    }

    private fun conversationRow(preview: ConversationPreview, colors: Palette): View = LinearLayout(this).apply {
        orientation = LinearLayout.VERTICAL
        val padding = dp(16)
        setPadding(padding, padding, padding, padding)
        minimumHeight = dp(GlazeClientTokens.InteractionFloorDp)
        contentDescription = preview.accessibilitySummary()
        isClickable = false
        isFocusable = true
        background = cardBackground(colors)

        val heading = LinearLayout(this@MessengerClientActivity).apply {
            orientation = LinearLayout.HORIZONTAL
            gravity = Gravity.CENTER_VERTICAL
        }
        heading.addView(text(preview.displayName, 16f, colors.text, Typeface.BOLD), LinearLayout.LayoutParams(0, ViewGroup.LayoutParams.WRAP_CONTENT, 1f))
        heading.addView(text(preview.timestampLabel, 12f, colors.muted, Typeface.NORMAL))
        addView(heading)
        preview.secondaryIdentity?.let { addView(spacer(3)); addView(text(it, 13f, colors.muted, Typeface.NORMAL)) }
        addView(spacer(7))
        addView(text(preview.previewText, 14f, colors.text, Typeface.NORMAL))
        addView(spacer(7))
        val stateParts = buildList {
            add(preview.provenance.displayLabel())
            if (preview.unreadCount > 0) add("${preview.unreadCount} unread")
            if (preview.isPinned) add("Pinned")
            if (preview.isMuted) add("Muted")
            add("Development preview")
        }
        addView(text(stateParts.joinToString(" · "), 12f, colors.muted, Typeface.BOLD))
    }

    private fun messageBubble(message: ConversationMessagePreview, colors: Palette): View = LinearLayout(this).apply {
        orientation = LinearLayout.VERTICAL
        val padding = dp(14)
        setPadding(padding, padding, padding, padding)
        val params = LinearLayout.LayoutParams(dp(300), ViewGroup.LayoutParams.WRAP_CONTENT)
        params.gravity = if (message.direction == MessageDirection.OUTGOING) Gravity.END else Gravity.START
        layoutParams = params
        background = cardBackground(colors)
        addView(text(message.senderLabel, 12f, colors.muted, Typeface.BOLD))
        message.replyToLabel?.let {
            addView(spacer(4))
            addView(text(it, 12f, colors.muted, Typeface.BOLD))
        }
        addView(spacer(5))
        addView(text(message.body, 15f, colors.text, Typeface.NORMAL))
        addView(spacer(6))
        addView(text("${message.timestampLabel} · ${message.statusLabel()}", 11f, colors.muted, Typeface.NORMAL))
    }

    private fun disabledComposer(colors: Palette): View = LinearLayout(this).apply {
        orientation = LinearLayout.VERTICAL
        val padding = dp(16)
        setPadding(padding, padding, padding, padding)
        background = cardBackground(colors)
        isEnabled = false
        contentDescription = "Message composer unavailable. Identity, conversation authorization, Data transport, and verified active E2EE are not available."
        addView(text("Message composer unavailable", 15f, colors.text, Typeface.BOLD))
        addView(spacer(5))
        addView(text("No message input or Send control is enabled until live Identity, exact conversation authorization, GoreeCloud Data transport, and verified active E2EE are independently available.", 13f, colors.muted, Typeface.NORMAL))
    }

    private fun disconnectedReadiness(): DataMessagingReadiness.Result = DataMessagingReadiness.evaluate(
        DataMessagingReadiness.Evidence(
            identity = DataMessagingReadiness.IdentityState.UNKNOWN,
            conversationAccess = DataMessagingReadiness.ConversationAccessState.UNKNOWN,
            transport = DataMessagingReadiness.DataTransportState.UNKNOWN,
            cryptography = DataMessagingReadiness.CryptographicState.UNKNOWN,
        ),
    )

    private fun readinessExplanation(result: DataMessagingReadiness.Result): String = when (result) {
        is DataMessagingReadiness.Result.Ready -> "Data messaging prerequisites are independently verified (${result.provenance.displayLabel()})."
        is DataMessagingReadiness.Result.Blocked -> {
            val labels = listOfNotNull(
                "Identity authentication".takeIf { DataMessagingReadiness.BlockReason.IDENTITY_NOT_AUTHENTICATED in result.reasons },
                "conversation authorization".takeIf { DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED in result.reasons },
                "GoreeCloud Data transport".takeIf { DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE in result.reasons },
                "verified active E2EE".takeIf { DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE in result.reasons },
            )
            "Missing verified prerequisites: ${labels.joinToString(" · ")}. No Send action is exposed."
        }
    }

    private fun surface(title: String, body: String, colors: Palette): View = LinearLayout(this).apply {
        orientation = LinearLayout.VERTICAL
        val padding = dp(18)
        setPadding(padding, padding, padding, padding)
        minimumHeight = dp(GlazeClientTokens.InteractionFloorDp)
        background = cardBackground(colors)
        addView(text(title, 16f, colors.text, Typeface.BOLD))
        addView(spacer(6))
        addView(text(body, 14f, colors.muted, Typeface.NORMAL))
    }

    private fun cardBackground(colors: Palette) = GradientDrawable().apply {
        shape = GradientDrawable.RECTANGLE
        cornerRadius = dp(GlazeClientTokens.SurfaceRadiusDp).toFloat()
        setColor(colors.surface)
        setStroke(dp(1), colors.border)
    }

    private fun text(value: String, sizeSp: Float, color: Int, style: Int): TextView = TextView(this).apply {
        text = value
        textSize = sizeSp
        setTextColor(color)
        typeface = Typeface.create(Typeface.DEFAULT, style)
        setLineSpacing(0f, 1.08f)
    }

    private fun spacer(heightDp: Int): View = View(this).apply {
        layoutParams = LinearLayout.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, dp(heightDp))
    }

    private fun dp(value: Int): Int = (value * resources.displayMetrics.density).toInt()
    private fun dp(value: Float): Int = (value * resources.displayMetrics.density).toInt()

    private fun palette(): Palette = if (isDark) {
        Palette(GlazeClientTokens.DarkCanvas.toInt(), GlazeClientTokens.DarkSurface.toInt(), GlazeClientTokens.DarkText.toInt(), GlazeClientTokens.DarkMutedText.toInt(), GlazeClientTokens.DarkBorder.toInt())
    } else {
        Palette(GlazeClientTokens.LightCanvas.toInt(), GlazeClientTokens.LightSurface.toInt(), GlazeClientTokens.LightText.toInt(), GlazeClientTokens.LightMutedText.toInt(), GlazeClientTokens.LightBorder.toInt())
    }

    private data class Palette(val canvas: Int, val surface: Int, val text: Int, val muted: Int, val border: Int)
}
