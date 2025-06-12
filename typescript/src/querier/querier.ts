export interface Querier {
  query(url: string): Promise<any>;
}

export class JSONQuerier implements Querier {
  async query(url: string): Promise<any> {
    const response = await fetch(url);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const result = await response.json();
    return result;
  }
}