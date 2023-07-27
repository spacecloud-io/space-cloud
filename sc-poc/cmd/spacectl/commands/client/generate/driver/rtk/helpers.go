package rtk

var HelperTS string = `
import { FetchArgs } from "@reduxjs/toolkit/dist/query";

interface cacheTag {
  type: string;
  key: string;
}

export const prepareQuery = (path: string, method: string) => (arg: any): FetchArgs => {
  const fetchArgs: FetchArgs = {url: path, method: method};
  if (!arg) {
    return fetchArgs;
  }

  if (["GET", "DELETE"].includes(method)) {
    let params: Record<string, any> = {};
    Object.keys(arg).forEach(key => {
      let v = arg[key];
      if (Array.isArray(v) || (typeof v === "object")) {
        v = JSON.stringify(v);
      }
      params[key] = v;
    })
    fetchArgs.params = params;
  } else {
    fetchArgs.body = arg;
  }

  return fetchArgs;
}

export const getTags = (tags: cacheTag[]): any => {
  return (result: any, error: any, arg: any) => tags.reduce((prev: any, cur) => {
    if (cur.key.startsWith('$')) {
      return [...prev, { type: cur.type, id: loadValue(cur.key.slice("$".length), arg) }];
    }

    if (cur.key.startsWith("/") && !error) {
      const keys: string[] = cur.key.slice("/".length).split("/");
      return [...prev, ...prepareTags(keys, cur.type, result)];
    }

    return [...prev, { type: cur.type, id: cur.key }];
  }, []);
};

const prepareTags = (keys: string[], type: string, obj: any): any => {
  if (keys.length == 0) {
    return [{ type, id: obj }];
  }

  if (Array.isArray(obj)) {
    let tags: any = [];
    obj.forEach(item => {
      tags = [...tags, ...prepareTags(keys.slice(1), type, item[keys[0]])];
    });
    return tags;
  }

  return prepareTags(keys.slice(1), type, obj[keys[0]]);
};

const loadValue = (key: string, obj: any) => (
  key
    .replace(/\[([^\[\]]*)\]/g, '.$1.')
    .split('.')
    .filter(t => t !== '')
    .reduce((prev, cur) => prev && prev[cur], obj)
);

`
