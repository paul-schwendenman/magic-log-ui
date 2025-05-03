import { persistentStore } from './persistentStores';

export const sessionStorageStore = <T>(key: string, initial: T) =>
	persistentStore(key, initial, () => sessionStorage);
