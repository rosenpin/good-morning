export interface API {
  apikey: string;
  cx: string;
}

export interface Image {
  size: string;
  rights: string;
  lifeSpan: number;
}

export interface Search {
  randomness: number;
  safe: string;
  baseQuery: string;
  features: string[];
}

export interface Config {
  api: API;
  image: Image;
  search: Search;
  maxDailyReload: number;
}