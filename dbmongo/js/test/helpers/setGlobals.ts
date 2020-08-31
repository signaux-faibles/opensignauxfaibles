/*global globalThis*/

export const setGlobals = (globals: unknown): typeof globalThis =>
  Object.assign(globalThis, globals)
