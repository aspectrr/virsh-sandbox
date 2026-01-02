┌────────────────────────────────────────────────────────────────────────────┐
│                         Agent Connection Flow                               │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                            │
│  1. Agent calls tmux-client API                                            │
│     POST /v1/sandbox/sessions/create                                       │
│     { "sandbox_id": "SBX-abc123" }                                         │
│              │                                                             │
│              ▼                                                             │
│  2. tmux-client generates ephemeral key pair                               │
│     ┌─────────────────────────────────────────────┐                        │
│     │ /tmp/sandbox-keys/sandbox_SBX-abc123_...    │                        │
│     │   - Private key (ed25519)                   │                        │
│     │   - Public key                              │                        │
│     └─────────────────────────────────────────────┘                        │
│              │                                                             │
│              ▼                                                             │
│  3. tmux-client calls virsh-sandbox API                                    │
│     POST /v1/access/request                                                │
│     { "sandbox_id": "...", "public_key": "ssh-ed25519 ..." }              │
│              │                                                             │
│              ▼                                                             │
│  4. virsh-sandbox issues certificate (5 min TTL)                           │
│     Returns: { "certificate": "ssh-ed25519-cert-v01...",                   │
│                "vm_ip_address": "192.168.122.10" }                         │
│              │                                                             │
│              ▼                                                             │
│  5. tmux-client saves certificate                                          │
│     /tmp/sandbox-keys/sandbox_SBX-abc123_...-cert.pub                      │
│              │                                                             │
│              ▼                                                             │
│  6. tmux-client creates tmux session with SSH command                      │
│     tmux new-session -d -s sandbox_SBX-abc123 \                            │
│       "ssh -i /tmp/.../key -o CertificateFile=/tmp/.../key-cert.pub \      │
│        sandbox@192.168.122.10"                                             │
│              │                                                             │
│              ▼                                                             │
│  7. Agent is now in a tmux session connected to the sandbox VM             │
│     ┌────────────────────────────────────────────┐                         │
│     │ sandbox@vm:~$                              │                         │
│     │ (tmux session - no shell escape)           │                         │
│     └────────────────────────────────────────────┘                         │
│                                                                            │
│  8. When done: DELETE /v1/sandbox/sessions/sandbox_SBX-abc123              │
│     - Kills tmux session                                                   │
│     - Deletes ephemeral keys                                               │
│                                                                            │
└────────────────────────────────────────────────────────────────────────────┘
