import { Config } from '../config/types';
import { Cache, IMAGE_KEY } from '../caching/cache';
import { Querier } from '../querier/querier';
import { Creator, GoogleImagesParams } from '../url/creator';
import { Parser } from '../result/parser';

export interface ImageProvider {
  provide(): Promise<ReadableStream>;
  forceReload(): Promise<ReadableStream>;
}

export class ImageProviderImpl implements ImageProvider {
  constructor(
    private querier: Querier,
    private urlCreator: Creator,
    private parser: Parser,
    private cache: Cache,
    private config: Config
  ) {}

  async provide(): Promise<ReadableStream> {
    if (await this.isCacheValid()) {
      return await this.loadFromCache();
    }

    console.log('cache invalid, reloading..');
    return await this.getNewImage();
  }

  async forceReload(): Promise<ReadableStream> {
    console.log('force reloading image');
    return await this.getNewImage();
  }

  private async getNewImage(): Promise<ReadableStream> {
    const link = await this.getImageURL();
    await this.cache.save(IMAGE_KEY, link);
    return await this.loadFromCache();
  }

  private async isCacheValid(): Promise<boolean> {
    try {
      const cacheAge = await this.cache.age(IMAGE_KEY);
      const maxAge = this.config.image.lifeSpan * 60 * 60 * 1000; // Convert hours to milliseconds
      return cacheAge < maxAge;
    } catch {
      return false;
    }
  }

  private async loadFromCache(): Promise<ReadableStream> {
    try {
      return await this.cache.load(IMAGE_KEY);
    } catch (error) {
      throw new Error(`Failed to load from cache: ${error}`);
    }
  }

  private async getImageURL(): Promise<string> {
    const query = await this.urlCreator.create(this.configToParams(this.config));
    
    console.log('sending request:', query);
    const result = await this.querier.query(query);
    const link = await this.parser.parse(result);
    
    if (typeof link !== 'string') {
      throw new Error(`unexpected result returned from parser, ${typeof link}:${link}`);
    }

    return link;
  }

  private configToParams(config: Config): GoogleImagesParams {
    return {
      apiKey: config.api.apikey,
      apiCX: config.api.cx,
      baseQuery: config.search.baseQuery,
      imgSize: config.image.size,
      randomness: config.search.randomness,
      features: config.search.features,
      resultNum: 1,
    };
  }
}

export function newImageProvider(
  querier: Querier,
  urlCreator: Creator,
  parser: Parser,
  config: Config,
  cache: Cache
): ImageProvider {
  return new ImageProviderImpl(querier, urlCreator, parser, cache, config);
}