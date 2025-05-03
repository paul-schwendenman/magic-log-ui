import { persistentStore } from './persistentStores';

export const localStorageStore = <T>(key: string, initial: T) =>
	persistentStore(key, initial, localStorage);
