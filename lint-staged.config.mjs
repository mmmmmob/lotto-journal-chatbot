/** @type {import('lint-staged').Config} */
const config = {
  '**/*.{js,jsx,ts,tsx}': [
    'pnpm exec prettier --write',
    'pnpm exec eslint --fix --max-warnings=0',
  ],
  '**/*.{json,md,mdx,css,scss,html,yml,yaml}': ['pnpm exec prettier --write'],
};

export default config;
