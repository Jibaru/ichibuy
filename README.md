# Ichibuy

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub last commit](https://img.shields.io/github/last-commit/Jibaru/ichibuy.svg)](https://github.com/Jibaru/ichibuy/commits/main)

This repository contains the source code for the Ichibuy microservices-based application.

## C4 Diagrams

### System Context Diagram

```mermaid
C4Context
  title System Context diagram for Ichibuy

  Person(user, "User", "A user of the Ichibuy platform.")
  System(ichibuy, "Ichibuy", "The main Ichibuy system.")

  System_Ext(google, "Google", "Provides OAuth 2.0 for authentication.")
  System_Ext(uploadthing, "UploadThing", "Provides file storage.")

  Rel(user, ichibuy, "Uses")
  Rel(ichibuy, google, "Uses for authentication")
  Rel(ichibuy, uploadthing, "Uses for file storage")
```

### Container Diagram

```mermaid
C4Container
  title Container diagram for Ichibuy

  Person(user, "User", "A user of the Ichibuy platform.")

  System_Boundary(ichibuy, "Ichibuy") {
    Container(auth, "Auth Service", "Go", "Handles user authentication and authorization.")
    Container(fstorage, "File Storage Service", "TypeScript", "Provides an API for uploading and managing files.")
    Container(store, "Store Service", "Go", "Manages stores, customers, and products.")
  }

  System_Ext(google, "Google", "Provides OAuth 2.0 for authentication.")
  System_Ext(uploadthing, "UploadThing", "Provides file storage.")

  Rel(user, auth, "Authenticates using")
  Rel(user, fstorage, "Uploads files to")
  Rel(user, store, "Manages stores, customers, and products via")

  Rel(auth, google, "Uses", "OAuth 2.0")
  Rel(fstorage, uploadthing, "Stores files in")
  Rel(store, auth, "Validates JWT tokens from")
```

## Microservices

| Service | Language | Description | Documentation |
| --- | --- | --- | --- |
| [Authentication Service](/auth) | ![Language](https://img.shields.io/badge/language-Go-blue.svg) | Handles user authentication and authorization using Google OAuth2 and JWT. | [API Docs](https://ichibuy-auth.vercel.app/api/swagger/index.html) / [README](/auth/README.md) |
| [File Storage Service](/fstorage) | ![Language](https://img.shields.io/badge/language-TypeScript-blue.svg) | Provides an API for uploading and managing files. | [API Docs](https://ichibuy-fstorage.vercel.app/api/swagger) / [README](/fstorage/README.md) |
| [Store Service](/store) | ![Language](https://img.shields.io/badge/language-Go-blue.svg) | Manages stores, customers, and products. | [API Docs](https://ichibuy-store.vercel.app/api/swagger/index.html) / [README](/store/README.md) |