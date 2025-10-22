import { PrismaClient } from '@prisma/client';
import { config as loadEnvFile } from 'dotenv';
import { existsSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const ensureDatabaseUrl = () => {
  if (process.env.DATABASE_URL) return;

  const currentDir = dirname(fileURLToPath(import.meta.url));
  const packageRoot = resolve(currentDir, '..');
  const workspaceRoot = resolve(packageRoot, '..', '..');
  const candidateEnvFiles = [
    resolve(packageRoot, '.env.local'),
    resolve(packageRoot, '.env'),
    resolve(workspaceRoot, '.env.local'),
    resolve(workspaceRoot, '.env'),
  ];

  for (const envPath of candidateEnvFiles) {
    if (!existsSync(envPath)) continue;
    loadEnvFile({ path: envPath });
    if (process.env.DATABASE_URL) break;
  }
};

ensureDatabaseUrl();

const globalForPrisma = globalThis as unknown as { prisma?: PrismaClient };

export const prisma =
  globalForPrisma.prisma ??
  new PrismaClient({
    log: process.env.NODE_ENV === 'development' ? ['error', 'warn'] : ['error'],
  });

if (process.env.NODE_ENV !== 'production') globalForPrisma.prisma = prisma;

export * from '@prisma/client';
