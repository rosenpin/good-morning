export interface Parser {
  parse(data: any): Promise<any>;
}

const RESPONSE_ERR = 'unexpected google images response format';

export class GoogleImagesResultParser implements Parser {
  async parse(raw: any): Promise<string> {
    if (!raw || typeof raw !== 'object') {
      throw new Error(`${RESPONSE_ERR}: ${JSON.stringify(raw)}`);
    }

    const items = raw.items;
    if (!Array.isArray(items)) {
      throw new Error(`${RESPONSE_ERR}: items not found or not array`);
    }

    if (items.length < 1) {
      throw new Error(`${RESPONSE_ERR}: no items returned`);
    }

    const item = items[0];
    if (!item || typeof item !== 'object') {
      throw new Error(`${RESPONSE_ERR}: invalid item format`);
    }

    const link = item.link;
    if (typeof link !== 'string') {
      throw new Error(`${RESPONSE_ERR}: link not found or not string`);
    }

    return link;
  }
}