import { loadConfig } from './config/loader';
import { newImageCache } from './caching/cache';
import { JSONQuerier } from './querier/querier';
import { GoogleImagesCreator } from './url/creator';
import { GoogleImagesResultParser } from './result/parser';
import { newImageProvider } from './provider/image';
import { QuotaManager, QuotaService } from './quota/manager';

export { QuotaManager };

interface Env {
  CACHE_KV: KVNamespace;
  QUOTA_MANAGER: DurableObjectNamespace;
  API_KEY: string;
  API_CX: string;
  IMAGE_SIZE: string;
  IMAGE_RIGHTS: string;
  IMAGE_LIFESPAN: string;
  SEARCH_RANDOMNESS: string;
  SEARCH_FEATURES: string;
  SEARCH_SAFE: string;
  SEARCH_BASE_QUERY: string;
  MAX_DAILY_RELOAD: string;
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    const pathname = url.pathname;

    try {
      const config = loadConfig(env);
      const cache = newImageCache(env.CACHE_KV);
      const querier = new JSONQuerier();
      const urlCreator = new GoogleImagesCreator();
      const parser = new GoogleImagesResultParser();
      const provider = newImageProvider(querier, urlCreator, parser, config, cache);
      const quotaService = new QuotaService(env.QUOTA_MANAGER);

      if (pathname === '/change') {
        return await handleChange(provider, quotaService, config.maxDailyReload);
      } else if (pathname === '/') {
        return await handleRoot(provider);
      } else if (pathname === '/favicon.ico') {
        return new Response('', { status: 204 });
      } else {
        return new Response('Not Found', { status: 404 });
      }
    } catch (error) {
      console.error('Error:', error);
      return new Response(`Internal Server Error: ${error}`, { status: 500 });
    }
  },

  async scheduled(controller: ScheduledController, env: Env): Promise<void> {
    const quotaService = new QuotaService(env.QUOTA_MANAGER);
    await quotaService.resetCounter();
    console.log('Daily quota counter reset');
  },
};

async function handleChange(
  provider: any,
  quotaService: QuotaService,
  maxDailyReload: number
): Promise<Response> {
  const { count } = await quotaService.getCurrentCount();
  
  if (count >= maxDailyReload) {
    return new Response(
      'You can no longer request reloads today. Daily quota exceeded',
      { status: 403 }
    );
  }

  await quotaService.incrementCounter();

  try {
    await provider.forceReload();
    return new Response('Image reloaded successfully', { status: 200 });
  } catch (error) {
    return new Response(
      `Error reload image: ${error}`,
      { status: 501 }
    );
  }
}

async function handleRoot(provider: any): Promise<Response> {
  try {
    const imageStream = await provider.provide();
    
    return new Response(imageStream, {
      status: 200,
      headers: {
        'Content-Type': 'image/jpeg',
        'Cache-Control': 'public, max-age=3600',
      },
    });
  } catch (error) {
    return new Response(`${error}`, { status: 501 });
  }
}