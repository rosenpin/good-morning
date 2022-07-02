# good-morning

A server side that provides a good morning image every day

Design can be found [here](https://www.lucidchart.com/documents/view/7f36558f-a031-4421-891c-6121e79f5d73/0_0)

Example can be found [here](https://image.rosenpin.io)

Example config
``` yaml
api:
    apikey: xxxxx
    cx: xxxxx
image:
    size: all # huge, icon, large, medium, small, xlarge, xxlarge or, use `all` for all sizes
    rights: cc_sharealike # cc_publicdomain, cc_attribute, cc_sharealike, cc_noncommercial, cc_nonderived
    lifeSpan: 16 # in hours
search:
    randomness: 100
    features: friend,flower,today,green,wish,card,smile,dogs,cute,funny,religious
    safe: active # or off
    baseQuery: good morning image
maxDailyReload: 3
```
