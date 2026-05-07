import { FlatCompat } from '@eslint/eslintrc';
import js from '@eslint/js';
import globals from 'globals';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
});

const nextCompat = new FlatCompat({
  baseDirectory: path.join(__dirname, 'apps/web'),
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
});

const nextFiles = [
  'apps/web/**/*.{js,jsx,ts,tsx}',
  'apps/web/next.config.{js,ts,mjs,cjs}',
  'apps/web/postcss.config.{js,ts,mjs,cjs}',
];

const nodeFiles = ['apps/api/**/*.{js,ts}', 'packages/db/**/*.{js,ts}'];

export default [
  {
    ignores: [
      '**/node_modules/**',
      '**/.turbo/**',
      '**/.next/**',
      '**/dist/**',
      '**/build/**',
      '**/out/**',
      'pnpm-lock.yaml',
    ],
  },
  {
    files: ['**/*.{js,jsx,ts,tsx}'],
    languageOptions: {
      globals: {
        ...globals.es2024,
      },
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
    },
  },
  ...compat.extends('eslint:recommended', 'plugin:@typescript-eslint/recommended'),
  {
    files: nodeFiles,
    languageOptions: {
      globals: {
        ...globals.es2024,
        ...globals.node,
      },
    },
  },
  {
    files: nextFiles,
    languageOptions: {
      globals: {
        ...globals.es2024,
        ...globals.browser,
      },
    },
  },
  ...nextCompat.extends('next/core-web-vitals', 'next/typescript').map((config) => ({
    ...config,
    files: nextFiles,
    settings: {
      ...config.settings,
      next: {
        ...(config.settings?.next ?? {}),
        rootDir: ['apps/web'],
      },
    },
    rules: {
      ...config.rules,
      '@next/next/no-html-link-for-pages': 'off',
    },
  })),
  {
    files: ['apps/web/next-env.d.ts'],
    rules: {
      '@typescript-eslint/triple-slash-reference': 'off',
    },
  },
];
