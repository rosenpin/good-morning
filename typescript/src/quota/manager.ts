export class QuotaManager {
  private state: DurableObjectState;

  constructor(state: DurableObjectState) {
    this.state = state;
  }

  async fetch(request: Request): Promise<Response> {
    const url = new URL(request.url);
    const action = url.pathname.substring(1);

    switch (action) {
      case 'increment':
        return this.handleIncrement();
      case 'reset':
        return this.handleReset();
      case 'get':
        return this.handleGet();
      default:
        return new Response('Not found', { status: 404 });
    }
  }

  private async handleIncrement(): Promise<Response> {
    const currentCount = await this.getCurrentCount();
    const newCount = currentCount + 1;
    
    await this.state.storage.put('dailyCounter', newCount);
    await this.state.storage.put('lastUpdated', Date.now());
    
    return new Response(JSON.stringify({ count: newCount }));
  }

  private async handleReset(): Promise<Response> {
    await this.state.storage.put('dailyCounter', 0);
    await this.state.storage.put('lastUpdated', Date.now());
    
    return new Response(JSON.stringify({ count: 0, reset: true }));
  }

  private async handleGet(): Promise<Response> {
    const count = await this.getCurrentCount();
    const lastUpdated = await this.state.storage.get('lastUpdated') || 0;
    
    return new Response(JSON.stringify({ 
      count, 
      lastUpdated,
      shouldReset: this.shouldReset(lastUpdated as number)
    }));
  }

  private async getCurrentCount(): Promise<number> {
    const count = await this.state.storage.get('dailyCounter');
    return (count as number) || 0;
  }

  private shouldReset(lastUpdated: number): boolean {
    const now = Date.now();
    const oneDayMs = 24 * 60 * 60 * 1000;
    return (now - lastUpdated) > oneDayMs;
  }
}

export class QuotaService {
  constructor(private quotaManager: DurableObjectNamespace) {}

  async incrementCounter(): Promise<number> {
    const id = this.quotaManager.idFromName('daily-quota');
    const stub = this.quotaManager.get(id);
    
    const response = await stub.fetch('http://quota/increment');
    const result = await response.json() as { count: number };
    
    return result.count;
  }

  async resetCounter(): Promise<void> {
    const id = this.quotaManager.idFromName('daily-quota');
    const stub = this.quotaManager.get(id);
    
    await stub.fetch('http://quota/reset');
  }

  async getCurrentCount(): Promise<{ count: number; shouldReset: boolean }> {
    const id = this.quotaManager.idFromName('daily-quota');
    const stub = this.quotaManager.get(id);
    
    const response = await stub.fetch('http://quota/get');
    const result = await response.json() as { count: number; shouldReset: boolean };
    
    if (result.shouldReset) {
      await this.resetCounter();
      return { count: 0, shouldReset: true };
    }
    
    return result;
  }
}