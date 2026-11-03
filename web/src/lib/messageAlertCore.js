// Pure helpers for new-message alerts; kept free of Svelte-rune imports so
// node:test can load them (glue lives in messageAlert.js).

/** Title/body/tag for a message. */
export function describeAlert(msg) {
  const tactical = msg.thread_kind === 'tactical';
  return {
    title: tactical ? msg.thread_key : msg.from_call,
    body: tactical ? `${msg.from_call}: ${msg.text}` : msg.text,
    tag: `${msg.thread_kind}:${msg.thread_key}`,
  };
}

/**
 * Should this transport change raise an alert? An unseen id after the feed has
 * caught up is what marks a message as new; the change kind is unreliable (a
 * message acked before the next poll arrives as "updated").
 * @param {{message?: object}} change
 * @param {{baselineDone: boolean, known: (id:number)=>boolean, muted: boolean}} ctx
 */
export function isAlertable(change, ctx) {
  const m = change?.message;
  if (!ctx.baselineDone || !m) return false;
  if (m.direction !== 'in' || !m.unread || m.is_ack || m.is_bulletin) return false;
  if (typeof m.id === 'number' && ctx.known(m.id)) return false;
  return !ctx.muted;
}
