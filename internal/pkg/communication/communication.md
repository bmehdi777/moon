# Communication

Packet will be organized like so :

┌─────────┬────────┬─────────┬─────────┬─────────┐
│ Version │  Type  │ LenAuth │ LenData │  Auth   │
│         │        │   (n)   │   (N)   │         │
│ 1 byte  │ 1 byte │ 4 bytes │ 8 bytes │ n bytes │
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

## Authentification / start

┌─────────┐         ┌─────────────┐       ┌──────────┐
│ Client  │         │ Auth server │       │  Server  │
└────┬────┘         └──────┬──────┘       └─────┬────┘
     │                     │                    │     
     │    Login:           │                    │     
     ├────────────────────►│                    │     
     │    OIDC authent     │                    │     
     │◄────────────────────┤                    │     
     │                     │                    │     
     │    Start:           │                    │     
     ├─────────────────────┼───────────────────►│     
     │    Access token     │                    │     
     │                     │ Access token valid?│     
     │                     │◄───────────────────┤     
     │                     │                    │     
     │◄────────────────────┼────────────────────┤     
     │                     │                    │     
     │                     │                    │     
     │                     │                    │

The authentification process starts with the *user*. He starts an oidc process
with the keycloak server. Keycloak is configured to deliver an *offline_token* 
(i.e. refresh_token that doesn't expire).

When the user will start his tunnel, the client will use the *refresh_token* to
acquire an *access_token*. This *access_token* will be in the header of each packet.
When the server will see that someone want to start a tunnel, it will check if the 
*access_token* is valid with the Keycloak public certificate.

If the *access_token* is valid, the tunnel will starts.
Otherwise, the server will refuse the communication.

## Logout

┌─────────┐         ┌─────────────┐
│ Client  │         │ Auth server │
└────┬────┘         └──────┬──────┘
     │                     │       
     │                     │       
     │    Revoke token     │       
     ├────────────────────►│       
     │                     │       
     │                     │       
     │                     │       
     ├─────┐               │   
     │     │               │  
     │     │Remove         │  
     │     │refresh_token  │  
     │     │               │  
     │◄────┘               │   
     │                     │ 

The user will ask the Keycloak server to revoke his *refresh_token*.
When this is done, the *refresh_token* will be removed from the local cache.
