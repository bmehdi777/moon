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

## Authentification

┌─────────┐         ┌─────────────┐
│ Client  │         │ Auth server │
└────┬────┘         └──────┬──────┘
     │                     │       
     │                     │       
     ├────────────────────►│       
     │     OIDC Authent    │       
     │◄────────────────────┤       
     │                     │       
     │                     │       
     │                     │       
     │                     │       
     │                     │

At the end of this process, client will have an `access token`. This token will
be sent in every request in a header field.
The server will use the public key of the auth server to verify the `access token`.

## 
