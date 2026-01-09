import * as SecureStore from "expo-secure-store"

export const secureStorage = {
  async getItem(key: string) {
    const value = await SecureStore.getItemAsync(key);
    return value;
  },
  async setItem(key: string, value: string) {
    await SecureStore.setItemAsync(key, value);
  },
  async removeItem(key: string) {
    await SecureStore.deleteItemAsync(key);
  },
};