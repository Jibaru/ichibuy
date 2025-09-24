import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';
import * as crypto from 'crypto';

interface JWK {
  kid: string;
  n: string;
  e: string;
  kty: string;
  use: string;
}

interface JWKS {
  keys: JWK[];
}

interface AuthenticatedRequest extends Request {
  userId?: string;
}

class JWTAuthMiddleware {
  private authBaseUrl: string;
  private cache: Map<string, string>;

  constructor(authBaseUrl: string) {
    this.authBaseUrl = authBaseUrl;
    this.cache = new Map<string, string>();
  }

  public validateToken() {
    return async (req: AuthenticatedRequest, res: Response, next: NextFunction) => {
      try {
        const authHeader = req.headers.authorization;
        if (!authHeader) {
          return res.status(401).json({ error: 'Authorization header is required' });
        }

        const parts = authHeader.split(' ');
        if (parts.length !== 2 || parts[0] !== 'Bearer') {
          return res.status(401).json({ error: 'Invalid authorization header format' });
        }

        const tokenString = parts[1];

        const payload = await new Promise<any>((resolve, reject) => {
          jwt.verify(tokenString, async (header, callback) => {
            try {
              const kid = header.kid;
              if (!kid) {
                return callback(new Error('kid not found in token header'));
              }

              const publicKey = await this.getPublicKey(kid);
              if (!publicKey) {
                return callback(new Error(`Public key not found for kid: ${kid}`));
              }

              callback(null, publicKey);
            } catch (error) {
              callback(error as Error);
            }
          }, { algorithms: ['RS256'] }, (error, decoded) => {
            if (error) {
              reject(error);
            } else {
              resolve(decoded);
            }
          });
        });

        if (!payload || !payload.user_id) {
          return res.status(401).json({ error: 'user_id not found in token' });
        }

        console.log('Token validated successfully', { user_id: payload.user_id });
        req.userId = payload.user_id as string;
        next();
      } catch (error) {
        console.error('JWT validation error:', error);
        return res.status(401).json({ error: `Invalid token: ${error instanceof Error ? error.message : 'Unknown error'}` });
      }
    };
  }

  private async getPublicKey(kid: string): Promise<string | null> {
    if (this.cache.has(kid)) {
      return this.cache.get(kid)!;
    }

    try {
      const jwksUrl = `${this.authBaseUrl}/api/v1/auth/.well-known/jwks.json`;
      const response = await fetch(jwksUrl);

      if (!response.ok) {
        throw new Error(`Failed to fetch JWKS: ${response.status}`);
      }

      const jwks: JWKS = await response.json();

      for (const key of jwks.keys) {
        if (key.kid === kid) {
          const publicKey = this.jwkToPem(key);
          this.cache.set(kid, publicKey);
          return publicKey;
        }
      }

      throw new Error(`Public key not found for kid: ${kid}`);
    } catch (error) {
      console.error('Error fetching public key:', error);
      return null;
    }
  }

  private jwkToPem(jwk: JWK): string {
    try {
      const keyObject = crypto.createPublicKey({
        key: {
          kty: 'RSA',
          n: jwk.n,
          e: jwk.e
        },
        format: 'jwk'
      });

      return keyObject.export({ type: 'spki', format: 'pem' }) as string;
    } catch (error) {
      console.error('Error converting JWK to PEM:', error);
      throw error;
    }
  }
}

export { JWTAuthMiddleware, AuthenticatedRequest };