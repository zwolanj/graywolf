import { test } from 'node:test';
import assert from 'node:assert/strict';
import { describeAlert, isAlertable } from './messageAlertCore.js';

const msg = (over = {}) => ({
  id: 5, direction: 'in', unread: true, from_call: 'W1ABC-9', text: 'hello',
  thread_kind: 'dm', thread_key: 'W1ABC-9', ...over,
});
const ctx = (over = {}) => ({ baselineDone: true, known: () => false, muted: false, ...over });

test('describeAlert: DM uses sender as title', () => {
  assert.deepEqual(describeAlert(msg()), { title: 'W1ABC-9', body: 'hello', tag: 'dm:W1ABC-9' });
});

test('describeAlert: tactical uses label as title and prefixes sender', () => {
  const m = msg({ thread_kind: 'tactical', thread_key: 'NW5WOPS' });
  assert.deepEqual(describeAlert(m), { title: 'NW5WOPS', body: 'W1ABC-9: hello', tag: 'tactical:NW5WOPS' });
});

test('isAlertable: unseen unread inbound alerts regardless of change kind', () => {
  assert.equal(isAlertable({ kind: 'created', message: msg() }, ctx()), true);
  assert.equal(isAlertable({ kind: 'updated', message: msg() }, ctx()), true);
});

test('isAlertable: baseline sync, outbound, read, ack and bulletin are silent', () => {
  assert.equal(isAlertable({ kind: 'created', message: msg() }, ctx({ baselineDone: false })), false);
  assert.equal(isAlertable({ kind: 'created', message: msg({ direction: 'out' }) }, ctx()), false);
  assert.equal(isAlertable({ kind: 'created', message: msg({ unread: false }) }, ctx()), false);
  assert.equal(isAlertable({ kind: 'created', message: msg({ is_ack: true }) }, ctx()), false);
  assert.equal(isAlertable({ kind: 'created', message: msg({ is_bulletin: true }) }, ctx()), false);
});

test('isAlertable: already-seen id and muted thread are silent', () => {
  assert.equal(isAlertable({ kind: 'created', message: msg() }, ctx({ known: () => true })), false);
  assert.equal(isAlertable({ kind: 'created', message: msg() }, ctx({ muted: true })), false);
});
