import { Config } from './types';

export function loadConfig(env: any): Config {
  return {
    api: {
      apikey: env.API_KEY || '',
      cx: env.API_CX || '',
    },
    image: {
      size: env.IMAGE_SIZE || 'all',
      rights: env.IMAGE_RIGHTS || 'cc_sharealike',
      lifeSpan: parseInt(env.IMAGE_LIFESPAN || '16', 10),
    },
    search: {
      randomness: parseInt(env.SEARCH_RANDOMNESS || '100', 10),
      safe: env.SEARCH_SAFE || 'active',
      baseQuery: env.SEARCH_BASE_QUERY || 'good morning image',
      features: (env.SEARCH_FEATURES || '').split(',').filter((f: string) => f.trim()),
    },
    maxDailyReload: parseInt(env.MAX_DAILY_RELOAD || '3', 10),
  };
}