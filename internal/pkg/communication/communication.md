# Communication

Packet will be organized like so :

┌─────────┬────────┬─────────┬─────────┬─────────┐
│ Version │  Type  │ LenAuth │ LenData │  Auth   │
│         │        │   (n)   │   (N)   │         │
│ 1 byte  │ 1 byte │ 4 bytes │ 4 bytes │ n bytes │
├─────────┴────────┴─────────┴─────────┴─────────┤
│                                                │
│                                                │
│                                                │
│                  Data / Payload                │
│                                                │
│                     N bytes                    │
│                                                │
│                                                │
│                                                │
└────────────────────────────────────────────────┘

When packet will be sent, they will be put in big endian.
Inversely, when received, packet will need to be put in little endian or 
converted directly to `Packet struct`
