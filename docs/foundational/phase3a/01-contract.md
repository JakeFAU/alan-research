# Frozen first-slice contract

## State and identity

Let d identify the original delegation by caller namespace, caller generation,
operation sequence, callee identity/generation and immutable request. Physical
recovery changes A's activation, not d. A fresh invocation is not recovery.
The finite model fixes d and the grant; it abstracts payloads to two distinct
opaque values r1 and r2. The actual world result is an oracle fact unavailable
to either agent except through B's retained receipt and authentic messages.

A retains its request, optional receipt D_A, and optional consumption marker C_A.
It may also hold a volatile received result V_A. B retains admission, a once-only
application tombstone, and optional result D_B. Request, result and ACK channels
are separate one-slot asynchronous channels. Crash erases A's volatile state,
leaves B and messages in flight unchanged, and makes A unavailable until recover.

## Transitions

Commit/send the original request; B admits it; B atomically applies it once and
records either r1 or r2; B transmits its retained result; A receives it in V_A;
A may durably record it in D_A; A may consume D_A once. Repeated original requests
do not apply again. B can resend the same retained result. Delivery, recording
and consumption are different events. Channel duplication/retransmission loops
are explored as finite-state cycles, not discarded as completed executions.

In the **delivery-ACK** profile, A may send ACK after volatile receipt. In the
**durable-ACK** profile, A may send ACK only after D_A records that receipt.
B's receipt-release transition consumes ACK and removes D_B while preserving
the application tombstone. No transition can recreate a released receipt by
reading the oracle's world fact. Consumption and ACK do not authorize mutation.

The model must permit B and the channels to progress while A is down. It must
not give B an enabling guard that reads A's private durable or volatile state.
Transport receipt authenticates the fixed d; no value or correlation is invented.

## Properties and evidence categories

Local safety: committed A facts survive crash; B's tombstone is monotone;
application count is at most one; every A receipt/consumption equals an authentic
result for d; consumption occurs at most once; authority is unchanged.

The candidate edge invariant after application is

    D_B = r  OR  D_A = r

where a consumed result remains retained in D_A. Volatile or in-flight bytes do
not discharge durable custody. Delivery-ACK is expected to violate this invariant
and to reach a recovered state with no receipt available to either agent or the
channels. Durable-ACK is expected to preserve it on every reachable state.

Check completion reachability separately: from each reachable state of the
durable-ACK model, some service schedule should reach A's authentic consumption.
This is an existential finite-graph check, not inevitable progress. Any liveness
argument additionally needs recovery and eventual service of enabled original
request, application, retransmission, delivery, recording and consumption work;
indefinite delay remains possible without those assumptions.

## Conditional argument to attack

Assume B initially establishes D_B with application; only authenticated copies
of r can populate D_A; A never deletes D_A; B releases only on an ACK causally
after D_A's durable commit; crashes preserve durable stores. Induction over local
steps preserves custody. This proof uses the edge contract plus local guarantees,
not a scan of a complete global execution history. It is a narrow safety lemma,
not a general composition or recovery theorem.

For impossibility, compare r1/r2 executions after B releases its receipt and A
loses its volatile copy, with empty result channel and no durable A receipt.
Erase the inaccessible oracle result when comparing recoverable states. If these
states coincide, any evidence-respecting recovery using only that state cannot
guarantee returning the actual result in both worlds. Guessing is not evidence;
withholding is safe but fails useful recovery. Allowing an authenticated external
lookup or retained B receipt changes the premises and removes this particular
argument. Do not present it as a new general distributed impossibility theorem.

## Deferred obligations

Executable correspondence to public ALAN programs and durable Store traces;
B-only and simultaneous crashes; non-atomic B effect/result recording; grant
attenuation and conservation; cancellation racing admission and results; two
outstanding opposite-direction obligations; resource bounds and fair service;
three-agent closure, hiding and associativity. None is discharged by this slice.
