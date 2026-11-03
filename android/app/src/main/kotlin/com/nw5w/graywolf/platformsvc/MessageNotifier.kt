package com.nw5w.graywolf.platformsvc

import android.Manifest
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.util.Log
import androidx.core.content.ContextCompat
import com.nw5w.graywolf.MainActivity
import com.nw5w.graywolf.R
import com.nw5w.graywolf.platformproto.MessageNotification

/** Pure text/validation helpers, kept Android-free so JVM unit tests can cover them. */
object MessageNotificationFormat {
    private val THREAD_ID_RE = Regex("^(dm|tactical):[A-Z0-9-]{1,9}$")

    fun threadId(n: MessageNotification) = "${n.threadKind}:${n.threadKey}"

    fun title(n: MessageNotification) = n.threadKey

    fun body(n: MessageNotification) =
        if (n.threadKind == "tactical") "${n.fromCall}: ${n.text}" else n.text

    // The id ends up in a WebView URL, so only accept the exact shape the SPA emits.
    fun isValidThreadId(id: String?) = id != null && THREAD_ID_RE.matches(id)
}

/**
 * Posts a system notification for each inbound APRS message the Go side
 * reports. The WebView is not alive in the background, so this is the only
 * path that can alert the operator.
 */
class MessageNotifier(
    private val context: Context,
    private val isAppForeground: () -> Boolean,
) {
    init {
        val channel = NotificationChannel(
            CHANNEL_ID,
            context.getString(R.string.message_channel_name),
            NotificationManager.IMPORTANCE_HIGH,
        )
        context.getSystemService(NotificationManager::class.java)?.createNotificationChannel(channel)
    }

    fun handle(n: MessageNotification) {
        if (isAppForeground()) return
        if (ContextCompat.checkSelfPermission(context, Manifest.permission.POST_NOTIFICATIONS)
            != PackageManager.PERMISSION_GRANTED) {
            Log.d(TAG, "POST_NOTIFICATIONS not granted; dropping message notification")
            return
        }
        val threadId = MessageNotificationFormat.threadId(n)
        val open = PendingIntent.getActivity(
            context,
            threadId.hashCode(),
            Intent(context, MainActivity::class.java)
                .putExtra(MainActivity.EXTRA_THREAD_ID, threadId)
                .addFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_SINGLE_TOP),
            PendingIntent.FLAG_IMMUTABLE or PendingIntent.FLAG_UPDATE_CURRENT,
        )
        val notif = android.app.Notification.Builder(context, CHANNEL_ID)
            .setContentTitle(MessageNotificationFormat.title(n))
            .setContentText(MessageNotificationFormat.body(n))
            .setSmallIcon(R.drawable.ic_notification)
            .setCategory(android.app.Notification.CATEGORY_MESSAGE)
            .setAutoCancel(true)
            .setContentIntent(open)
            .build()
        // One notification per thread: a new message replaces and re-alerts.
        context.getSystemService(NotificationManager::class.java)?.notify(threadId, NOTIF_ID, notif)
    }

    companion object {
        private const val TAG = "MessageNotifier"
        const val CHANNEL_ID = "graywolf-messages"
        private const val NOTIF_ID = 0x6758

        fun cancelAll(context: Context) {
            val nm = context.getSystemService(NotificationManager::class.java) ?: return
            for (sbn in nm.activeNotifications) {
                if (sbn.notification.channelId == CHANNEL_ID) nm.cancel(sbn.tag, sbn.id)
            }
        }
    }
}
