package com.nw5w.graywolf.platformsvc

import com.nw5w.graywolf.platformproto.MessageNotification
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class MessageNotificationFormatTest {
    private fun note(kind: String, key: String, from: String, text: String) =
        MessageNotification.newBuilder()
            .setThreadKind(kind).setThreadKey(key).setFromCall(from).setText(text).build()

    @Test
    fun dmUsesPeerAsTitleAndPlainBody() {
        val n = note("dm", "W1ABC-9", "W1ABC-9", "hello")
        assertEquals("dm:W1ABC-9", MessageNotificationFormat.threadId(n))
        assertEquals("W1ABC-9", MessageNotificationFormat.title(n))
        assertEquals("hello", MessageNotificationFormat.body(n))
    }

    @Test
    fun tacticalPrefixesSenderInBody() {
        val n = note("tactical", "NW5WOPS", "W1ABC-9", "net open")
        assertEquals("NW5WOPS", MessageNotificationFormat.title(n))
        assertEquals("W1ABC-9: net open", MessageNotificationFormat.body(n))
    }

    @Test
    fun threadIdValidationRejectsAnythingNotSpaShaped() {
        assertTrue(MessageNotificationFormat.isValidThreadId("dm:W1ABC-9"))
        assertTrue(MessageNotificationFormat.isValidThreadId("tactical:NW5WOPS"))
        assertFalse(MessageNotificationFormat.isValidThreadId(null))
        assertFalse(MessageNotificationFormat.isValidThreadId("bogus:W1ABC"))
        assertFalse(MessageNotificationFormat.isValidThreadId("dm:w1abc"))
        assertFalse(MessageNotificationFormat.isValidThreadId("dm:W1ABC';alert(1);'"))
        assertFalse(MessageNotificationFormat.isValidThreadId("dm:TOOLONGCALLSIGN"))
    }
}
