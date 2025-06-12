export interface Cache {
  save(key: string, value: any): Promise<void>;
  load(key: string): Promise<any>;
  age(key: string): Promise<number>;
}

export const IMAGE_KEY = 'ImageKey';

interface CacheData {
  url: string;
  imageData: ArrayBuffer;
  creation: number;
}

export class ImageCache implements Cache {
  constructor(private kv: KVNamespace) {}

  async age(key: string): Promise<number> {
    const data = await this.getCacheData(key);
    if (!data) {
      throw new Error(`invalid key provided ${key}`);
    }

    return Date.now() - data.creation;
  }

  async save(key: string, value: any): Promise<void> {
    const url = value as string;
    if (typeof url !== 'string') {
      throw new Error('invalid value provided to cache manager');
    }

    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`failed to fetch image: ${response.status}`);
    }

    const imageData = await response.arrayBuffer();
    
    const cacheData: CacheData = {
      url,
      imageData,
      creation: Date.now(),
    };

    await this.kv.put(key, JSON.stringify({
      url,
      creation: cacheData.creation,
    }));
    
    await this.kv.put(`${key}_data`, imageData);
  }

  async load(key: string): Promise<ReadableStream> {
    const data = await this.getCacheData(key);
    if (!data) {
      throw new Error(`invalid key provided ${key}`);
    }

    const imageData = await this.kv.get(`${key}_data`, 'arrayBuffer');
    if (!imageData) {
      throw new Error(`image data not found for key ${key}`);
    }

    return new ReadableStream({
      start(controller) {
        controller.enqueue(new Uint8Array(imageData));
        controller.close();
      }
    });
  }

  private async getCacheData(key: string): Promise<{ url: string; creation: number } | null> {
    const data = await this.kv.get(key);
    if (!data) {
      return null;
    }

    try {
      return JSON.parse(data);
    } catch {
      return null;
    }
  }
}

export function newImageCache(kv: KVNamespace): Cache {
  return new ImageCache(kv);
}