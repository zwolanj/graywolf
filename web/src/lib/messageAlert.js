// New-message alerts for the desktop SPA. Uses the browser Notification API
// when permission is granted, otherwise falls back to a short synthesized
// chime. Android never reaches this: its alerts are native (GraywolfService).
//
// Notification needs a secure context (HTTPS or localhost); graywolf serves
// plain HTTP, so on a LAN address the API is unavailable and only the chime plays.

import { messages } from './messagesStore.svelte.js';
import { messagesPreferencesState } from './settings/messages-preferences-store.svelte.js';
import { Platform } from './platform.js';
import { describeAlert, isAlertable } from './messageAlertCore.js';

/** @returns {'granted'|'denied'|'default'|'unavailable'} */
export function notificationState() {
  if (typeof window === 'undefined' || typeof Notification === 'undefined' || !window.isSecureContext) {
    return 'unavailable';
  }
  return Notification.permission;
}

/** Must be called from a user gesture (browsers ignore it otherwise). */
export async function requestNotificationPermission() {
  if (notificationState() === 'unavailable') return 'unavailable';
  try {
    return await Notification.requestPermission();
  } catch {
    return Notification.permission;
  }
}

let audioCtx = null;

// Browsers keep an AudioContext suspended until a user gesture.
function unlockAudio() {
  try {
    audioCtx ??= new (globalThis.AudioContext || globalThis.webkitAudioContext)();
    if (audioCtx.state === 'suspended') audioCtx.resume();
  } catch { /* no audio support */ }
}

export function armChime() {
  if (typeof window === 'undefined') return;
  for (const ev of ['pointerdown', 'keydown']) {
    window.addEventListener(ev, unlockAudio, { once: true, passive: true });
  }
}

export function playChime() {
  if (!audioCtx || audioCtx.state !== 'running') return;
  const t0 = audioCtx.currentTime;
  [[880, 0], [1175, 0.14]].forEach(([freq, offset]) => {
    const osc = audioCtx.createOscillator();
    const gain = audioCtx.createGain();
    osc.type = 'sine';
    osc.frequency.value = freq;
    gain.gain.setValueAtTime(0.0001, t0 + offset);
    gain.gain.exponentialRampToValueAtTime(0.2, t0 + offset + 0.02);
    gain.gain.exponentialRampToValueAtTime(0.0001, t0 + offset + 0.25);
    osc.connect(gain).connect(audioCtx.destination);
    osc.start(t0 + offset);
    osc.stop(t0 + offset + 0.3);
  });
}

function showNotification(msg) {
  const { title, body, tag } = describeAlert(msg);
  const n = new Notification(title, { body, tag });
  n.onclick = () => {
    window.focus();
    window.location.hash = `#/messages?thread=${encodeURIComponent(tag)}`;
    n.close();
  };
}

/** Called by the transport for every change before it is applied to the store. */
export function maybeAlert(change, baselineDone) {
  if (Platform.isAndroid || !messagesPreferencesState.notifyOnMessage) return;
  const m = change?.message;
  if (!m) return;
  const threadId = messages.threadIdFor(m.thread_kind, m.thread_key);
  const ok = isAlertable(change, {
    baselineDone,
    known: (id) => messages.messageById.has(id),
    muted: !!messages.conversations.get(threadId)?.muted,
  });
  if (!ok) return;
  if (notificationState() === 'granted') showNotification(m);
  else playChime();
}
