UPDATE identity_verifiable_addresses SET code = LEFT(SHA2(RANDOM_BYTES(32), 256), 32) WHERE code IS NULL
