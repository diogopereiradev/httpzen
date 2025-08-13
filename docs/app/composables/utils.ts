export function useUtils() {
  async function copyToClipboard(text: string): Promise<boolean> {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch {
      return false;
    }
  }

  function debounce(callback: (...args: unknown[]) => void, ttl: number) {
    let timeout: number | undefined;
    return function(this: unknown, ...args: unknown[]) {
      if (timeout !== undefined) {
        clearTimeout(timeout);
      }
      timeout = window.setTimeout(() => callback.apply(this, args), ttl);
    };
  }

  return {
    copyToClipboard,
    debounce,
  }
}