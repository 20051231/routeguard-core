# RouteGuard core derivative

This public derivative preserves the Hiddify and sing-box copyright notices and
GPL licenses. It pins its sing-box submodule and the RouteGuard sing-quic fork.

The owned `routeguard-normal` selector can retain all configured terminal
Hysteria2 clients as idle references. Only in-process Mobile methods may change
this selector, and they require its exact positive physical-network epoch.
`RouteGuardSelectKeep` deliberately preserves existing normal application flows;
new flows use the newly selected terminal. Immediate migration is implemented by
the Android application's full socket fence and core restart, never by this API.
AI and DNS rules must reference their independently fixed terminal, not this
normal selector. Android must check its generation before and after Mobile calls.

Session snapshots report published authenticated local QUIC context, an opaque
client-scoped session ID, authentication time and local stream/session counts.
They do not promise fresh remote reachability. Compare IDs only together with
core epoch and terminal tag. `Warm` joins normal coalesced authentication without
opening an application stream; idle references are established first. Physical
network transitions still revoke all old transports before new authentication.

Connection audit uses the router's captured outbound chain, so later KEEP
selection changes cannot rewrite an existing flow's reported terminal.

The dedicated Actions build includes gVisor, QUIC, uTLS, authenticated Clash,
gRPC and connection tracking, targeting Android API23 ARM64 and x86_64. Optional
Chromium, OpenVPN, OpenConnect, WireGuard and embedded Tor helpers are outside
this RouteGuard build. The strict Android compiler accepts only its constrained
Hysteria2 terminals and owned direct rules; arbitrary imported helper/detour
configuration is never admitted to the standby pool.
