export interface Creator {
  create(params: any): Promise<string>;
}

export interface GoogleImagesParams {
  apiKey: string;
  apiCX: string;
  baseQuery: string;
  imgSize: string;
  randomness: number;
  resultNum: number;
  features: string[];
}

export class GoogleImagesCreator implements Creator {
  private static readonly REQUEST_URL = 'https://www.googleapis.com/customsearch/v1';
  private static readonly IMAGE_SEARCH_TYPE = 'image';

  async create(rawParams: any): Promise<string> {
    const params = rawParams as GoogleImagesParams;
    
    if (!params.apiKey || !params.apiCX) {
      throw new Error('GoogleImagesCreator requires apiKey and apiCX');
    }

    const baseURL = new URL(GoogleImagesCreator.REQUEST_URL);
    const searchParams = this.addParams(params);
    
    searchParams.forEach((value, key) => {
      baseURL.searchParams.append(key, value);
    });

    return baseURL.toString();
  }

  private addParams(params: GoogleImagesParams): URLSearchParams {
    const searchParams = new URLSearchParams();
    
    searchParams.set('key', params.apiKey);
    searchParams.set('cx', params.apiCX);
    
    const randomFeature = params.features[Math.floor(Math.random() * params.features.length)];
    searchParams.set('q', `${params.baseQuery} ${randomFeature}`);
    
    if (params.imgSize !== 'all') {
      searchParams.set('imgSize', params.imgSize);
    }
    
    searchParams.set('start', Math.floor(Math.random() * params.randomness).toString());
    searchParams.set('searchType', GoogleImagesCreator.IMAGE_SEARCH_TYPE);
    searchParams.set('num', params.resultNum.toString());

    return searchParams;
  }
}